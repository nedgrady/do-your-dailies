package migrations

import (
	"do-your-dailies/server/internal/domain/models"

	"gorm.io/gorm"
)

// migrationLockKey is an arbitrary, fixed advisory lock key. Cloud Run can
// boot multiple instances concurrently (new deploys, scale-out), so
// migrations are serialized with a transaction-scoped Postgres advisory
// lock: only one instance runs AutoMigrate's DDL at a time, others block
// until it commits, then their own (now idempotent) AutoMigrate is a no-op.
const migrationLockKey = 727384

func Migrate(database *gorm.DB) error {
	return database.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT pg_advisory_xact_lock(?)", migrationLockKey).Error; err != nil {
			return err
		}

		println("Running migrations...")
		return tx.AutoMigrate(&models.User{}, &models.Chore{}, &models.ChoreCompletion{}, &models.ChoreInQueueCompletion{})
	})
}
