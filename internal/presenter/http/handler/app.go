package handler

import (
	"github.com/UndeadDemidov/ya-pr-diploma/internal/app"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/presenter/http/middleware"
)

type App struct {
	*Auth
}

func NewApp(mart *app.GopherMart, sessions *middleware.Sessions) *App {
	if mart == nil {
		panic("missing *app.GopherMart, parameter must not be nil")
	}
	if sessions == nil {
		panic("missing *middleware.Sessions, parameter must not be nil")
	}
	return &App{Auth: NewAuth(mart.Authenticator, sessions)}
}
