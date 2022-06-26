package handler

import "github.com/UndeadDemidov/ya-pr-diploma/internal/app"

type App struct {
	*Auth
}

func NewApp(mart *app.GopherMart) *App {
	return &App{Auth: NewAuth(mart.Authenticator)}
}
