package app

import (
	"fmt"
)

type App struct {
	*Engine
}

func NewApp() *App {
	app := &App{
		Engine: NewEngine(),
	}
	return app
}

func (a App) Start(port int) error {
	return a.Run(fmt.Sprintf(":%d", port))
}
