package conf

type Configurer interface {
	SetPFlag()
	Read() error
}

var _ Configurer = (*App)(nil)

type App struct {
	Server
	Database
	Externals
}

func NewAppConfig() *App {
	return &App{
		Server:    Server{},
		Database:  Database{},
		Externals: Externals{},
	}
}

func (a *App) SetPFlag() {
	a.Server.SetPFlag()
	a.Database.SetPFlag()
	a.Externals.SetPFlag()
}

func (a *App) Read() error {
	err := a.Server.Read()
	if err != nil {
		return err
	}
	err = a.Database.Read()
	if err != nil {
		return err
	}
	err = a.Externals.Read()
	if err != nil {
		return err
	}
	return nil
}
