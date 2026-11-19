package chores

import (
	"context"
	"fmt"
	"time"

	"do-your-dailies/server/internal/domain/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Store interface {
	List(ctx context.Context, userID uint) ([]models.Chore, error)
	Create(ctx context.Context, userID uint, req CreateRequest) (models.Chore, error)
	CreateMany(ctx context.Context, userID uint, reqs []CreateRequest) ([]models.Chore, error)
	Get(ctx context.Context, userID uint, id uint) (models.Chore, error)
	Update(ctx context.Context, userID uint, id uint, req UpdateRequest) (models.Chore, error)
	Delete(ctx context.Context, userID uint, id uint) error
}

type GormStore struct {
	db *gorm.DB
}

func NewGormStore(db *gorm.DB) *GormStore {
	return &GormStore{db: db}
}

func (store *GormStore) List(ctx context.Context, userID uint) ([]models.Chore, error) {
	var chores []models.Chore
	// Explicit order: Postgres returns an updated row at the end of an unordered
	// scan, so without this an edited chore jumps to the bottom of the list.
	result := store.db.WithContext(ctx).Where("user_id = ?", userID).Order("id ASC").Find(&chores)
	return chores, result.Error
}

func (store *GormStore) Create(ctx context.Context, userID uint, req CreateRequest) (models.Chore, error) {
	chore := models.Chore{
		UserID:        userID,
		Name:          req.Name,
		CadenceInDays: req.CadenceInDays,
		DisplayUnit:   req.DisplayUnit,
		Emoji:         req.Emoji,
	}

	err := store.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := EnsureUserExists(tx, userID); err != nil {
			return err
		}
		if err := tx.Create(&chore).Error; err != nil {
			return err
		}
		return createBackdatedCompletion(tx, chore.ID, req.LastCompletedAt)
	})
	return chore, err
}

func (store *GormStore) CreateMany(ctx context.Context, userID uint, reqs []CreateRequest) ([]models.Chore, error) {
	chores := make([]models.Chore, 0, len(reqs))
	for _, req := range reqs {
		chores = append(chores, models.Chore{
			UserID:        userID,
			Name:          req.Name,
			CadenceInDays: req.CadenceInDays,
			DisplayUnit:   req.DisplayUnit,
			Emoji:         req.Emoji,
		})
	}

	err := store.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := EnsureUserExists(tx, userID); err != nil {
			return err
		}
		if err := tx.Create(&chores).Error; err != nil {
			return err
		}
		for i, chore := range chores {
			if err := createBackdatedCompletion(tx, chore.ID, reqs[i].LastCompletedAt); err != nil {
				return err
			}
		}
		return nil
	})
	return chores, err
}

// createBackdatedCompletion records a chore as already completed at lastCompletedAt, so a
// chore imported with known history doesn't read as "just done now" (the fallback used when
// no completion row exists). A nil lastCompletedAt is a no-op, preserving that default.
func createBackdatedCompletion(tx *gorm.DB, choreID uint, lastCompletedAt *time.Time) error {
	if lastCompletedAt == nil {
		return nil
	}
	completion := models.ChoreCompletion{ChoreID: choreID}
	completion.CreatedAt = *lastCompletedAt
	return tx.Create(&completion).Error
}

// EnsureUserExists provisions a User row for userID if one doesn't already exist. Real users
// are always already seeded, so this is a no-op for them; it only actually matters for a
// not-yet-seen synthetic id acting through the admin simulator's impersonation header, letting
// a chore be created under it without violating the chores.user_id foreign key.
func EnsureUserExists(tx *gorm.DB, userID uint) error {
	user := models.User{Model: gorm.Model{ID: userID}, Name: fmt.Sprintf("sim-%d", userID)}
	return tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&user).Error
}

func (store *GormStore) Get(ctx context.Context, userID uint, id uint) (models.Chore, error) {
	var chore models.Chore
	result := store.db.WithContext(ctx).Where("user_id = ?", userID).First(&chore, id)
	return chore, result.Error
}

func (store *GormStore) Update(ctx context.Context, userID uint, id uint, req UpdateRequest) (models.Chore, error) {
	var chore models.Chore
	if err := store.db.WithContext(ctx).Where("user_id = ?", userID).First(&chore, id).Error; err != nil {
		return chore, err
	}
	if req.Name != nil {
		chore.Name = *req.Name
	}
	if req.CadenceInDays != nil {
		chore.CadenceInDays = *req.CadenceInDays
	}
	if req.DisplayUnit != nil {
		chore.DisplayUnit = *req.DisplayUnit
	}
	if req.Emoji != nil {
		chore.Emoji = *req.Emoji
	}
	result := store.db.WithContext(ctx).Save(&chore)
	return chore, result.Error
}

func (store *GormStore) Delete(ctx context.Context, userID uint, id uint) error {
	result := store.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.Chore{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
