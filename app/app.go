package app

import (
	"log/slog"
	"net/http"
)

type App struct {
	Server *http.Server // Pointer to an http.Server instance
	DB     DBManager
}

func (app *App) Init() {
	slog.Info("Init app")
}

func (app *App) MountRoutes() {
	slog.Info("Mounting routes")
	http.HandleFunc("/api/set", app.ApiSetUrl)
	http.HandleFunc("/api/get", app.ApiGetUrl)
	http.HandleFunc("/", app.RootHandler)
}

// Creates and inits a new app.
func Default(c Configuration) App {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}
	app := App{Server: &server, DB: SetupDB(c.BaseUrl)}
	app.Init()
	app.MountRoutes()
	return app
}

// Starts all services.
func (app *App) Run(port string) {
	slog.Info("Running app.")
	app.Server.ListenAndServe()
}
