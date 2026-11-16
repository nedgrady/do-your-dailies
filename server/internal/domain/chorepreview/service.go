package chorepreview

import (
	"context"
	"errors"
	"time"

	"do-your-dailies/server/internal/domain/chorequeue"
	"do-your-dailies/server/internal/domain/models"

	"gorm.io/gorm"
)

type DaySchedule struct {
	DayOffset int
	Chores    []chorequeue.ChoreInQueue
}

// errAlwaysRollback is returned from the Transaction closure below purely
// to force GORM to roll back rather than commit - this is a pure dry run,
// so nothing it writes must ever persist, regardless of success or failure.
var errAlwaysRollback = errors.New("chorepreview: intentional rollback, not a real failure")

// Simulate projects the next `days` days assuming every day's
// capacity-first recommendation is completed, so a chore due more than
// once in the window (e.g. daily) correctly reappears on every day it
// would actually be done - unlike a single higher-LIMIT query, which can
// only ever return each chore once.
//
// It works by running the exact same ranking query
// (chorequeue.RankChoresAsOf) once per simulated day inside a transaction
// that is never committed, writing each day's picks as real
// ChoreCompletion rows (backdated to that simulated day) so the next
// iteration's identical query naturally sees them as the chore's most
// recent completion. db.Transaction is used rather than db.Begin so that
// calling Simulate against a db that is *already* a transaction (as in
// tests using testhelpers.NewTransactionDB) correctly nests via a Postgres
// SAVEPOINT instead of breaking the outer transaction.
func Simulate(db *gorm.DB, ctx context.Context, userID uint, days int) ([]DaySchedule, error) {
	schedule := make([]DaySchedule, 0, days)
	now := time.Now()

	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for day := 1; day <= days; day++ {
			asOf := now.AddDate(0, 0, day)

			picks, err := chorequeue.RankChoresAsOf(tx, ctx, userID, asOf, chorequeue.DefaultMaxChores)
			if err != nil {
				return err
			}

			schedule = append(schedule, DaySchedule{DayOffset: day, Chores: picks})

			for _, pick := range picks {
				completion := models.ChoreCompletion{ChoreID: pick.ChoreID}
				completion.CreatedAt = asOf
				if err := tx.Create(&completion).Error; err != nil {
					return err
				}
			}
		}
		return errAlwaysRollback
	})

	if err != nil && !errors.Is(err, errAlwaysRollback) {
		return nil, err
	}
	return schedule, nil
}
