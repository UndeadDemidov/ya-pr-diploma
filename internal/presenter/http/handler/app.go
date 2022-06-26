package handler

type App struct {
	*Auth
}

func NewApp(auth Authenticator) *App {
	return &App{Auth: NewAuth(auth)}
}
