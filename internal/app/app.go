package app

type App struct {
}

func NewApp() *App {
	//TODO: init config

	//TODO: init database

	//TODO: init layers

	//TODO: init server
	return &App{}
}

func (a *App) Run() error {
	//TODO: run server

	//TODO: gracefull shutdown
	return nil
}
