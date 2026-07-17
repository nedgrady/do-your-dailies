package api

import (
	"do-your-dailies/server/internal/domain/chorecompletions"
	"do-your-dailies/server/internal/domain/choreinqueuecompletion"
	"do-your-dailies/server/internal/domain/chorequeue"
	"do-your-dailies/server/internal/domain/chores"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Application struct {
	Router                      *chi.Mux
	ChoreStore                  chores.Store
	ChoreCompletionStore        chorecompletions.Store
	ChoreQueueStore             chorequeue.Store
	ChoreInQueueCompletionStore choreinqueuecompletion.Store
	Db                          *gorm.DB
}

func New(db *gorm.DB) *Application {
	app := &Application{}
	if db != nil {
		app.ChoreStore = chores.NewGormStore(db)
		app.ChoreCompletionStore = chorecompletions.NewGormStore(db)
		app.ChoreQueueStore = chorequeue.NewGormStore(db)
		app.ChoreInQueueCompletionStore = choreinqueuecompletion.NewGormStore(db)
		app.Db = db
	}
	app.Router = app.setupRoutes()
	return app
}
