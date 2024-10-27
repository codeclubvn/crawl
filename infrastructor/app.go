package infrastructor

import (
	"crawl/initialization"
	mongodriven "go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	Env     *initialization.Database
	MongoDB *mongodriven.Client
}

func App() (*Application, *mongodriven.Client) {
	app := &Application{}
	app.Env = initialization.NewEnv()
	app.MongoDB = NewMongoDatabase(app.Env)
	return app, app.MongoDB
}

func (app *Application) CloseDBConnection() {
	CloseMongoDBConnection(app.MongoDB)
}
