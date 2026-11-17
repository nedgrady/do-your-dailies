package idempotency

import (
	"context"
	"errors"

	"do-your-dailies/server/internal/domain/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type store struct {
	db *gorm.DB
}

// begin claims key for this request with one atomic insert. The unique
// (scope, key) index means exactly one of two simultaneous requests inserts
// a row. Returns inserted=true for the winner; everyone else gets the existing row.
func (s store) begin(ctx context.Context, scope, key, requestHash string) (bool, models.IdempotencyKey, error) {
	row := models.IdempotencyKey{
		Scope:       scope,
		Key:         key,
		RequestHash: requestHash,
		State:       StateInProgress,
	}
	result := s.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&row)
	if result.Error != nil {
		return false, models.IdempotencyKey{}, result.Error
	}
	if result.RowsAffected == 1 {
		return true, row, nil
	}

	var existing models.IdempotencyKey
	err := s.db.WithContext(ctx).Where("scope = ? AND key = ?", scope, key).First(&existing).Error
	if err != nil {
		return false, models.IdempotencyKey{}, err
	}
	return false, existing, nil
}

func (s store) complete(ctx context.Context, id uint, status int, contentType string, body []byte) error {
	return s.db.WithContext(ctx).
		Model(&models.IdempotencyKey{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"state":                 StateCompleted,
			"response_status":       status,
			"response_content_type": contentType,
			"response_body":         body,
		}).Error
}

// release deletes an in-progress row so a retry can run the handler again. Used
// after a server error or panic, so failed attempts are not replayed as successes.
func (s store) release(ctx context.Context, id uint) error {
	err := s.db.WithContext(ctx).Delete(&models.IdempotencyKey{}, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}
