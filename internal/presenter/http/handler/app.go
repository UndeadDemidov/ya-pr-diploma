package handler

type App struct {
	*Auth
}

func NewApp(man CredentialManager) *App {
	return &App{Auth: NewAuth(man)}
}
