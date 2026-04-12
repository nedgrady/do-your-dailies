package api

import (
	"do-your-dailies/server/internal/domain/chorecompletions"
	"do-your-dailies/server/internal/domain/chores"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Application struct {
	Router               *chi.Mux
	ChoreStore           chores.Store
	ChoreCompletionStore chorecompletions.Store
}

func New(db *gorm.DB) *Application {
	app := &Application{}
	if db != nil {
		app.ChoreStore = chores.NewGormStore(db)
		app.ChoreCompletionStore = chorecompletions.NewGormStore(db)
	}
	app.Router = app.setupRoutes()
	return app
}
