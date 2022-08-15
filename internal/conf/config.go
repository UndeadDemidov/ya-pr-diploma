package conf

// App является композицией из составных частей конфигурации.
type App struct {
	Server
	Database
	Externals
}

// NewAppConfig создает новый экземпляр конфигурации.
func NewAppConfig() *App {
	return &App{
		Server:    Server{},
		Database:  Database{},
		Externals: Externals{},
	}
}

// SetPFlag устанавливает все флаги параметров вложенных частей композиции.
func (a *App) SetPFlag() {
	a.Server.SetPFlag()
	a.Database.SetPFlag()
	a.Externals.SetPFlag()
}

// Read считывает конфигурацию по всем вложенным частям композиции.
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
