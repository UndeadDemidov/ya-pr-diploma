package auth

import (
	"context"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/app"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/user"
	errors2 "github.com/UndeadDemidov/ya-pr-diploma/internal/errors"
	_ "github.com/golang/mock/mockgen/model"
)

//go:generate mockgen -destination=./mocks/mock_service.go . CredentialManager

type CredentialManager interface {
	AddNewUser(ctx context.Context, usr user.User, login, pword string) (err error)
	GetUser(ctx context.Context, login string) (usr user.User, err error)
	AuthenticateUser(ctx context.Context, login, pword string) (usr user.User, err error)
	// DisableUser(user user.User) error
}

var _ app.Authenticator = (*Service)(nil)

type Service struct {
	userSvc user.Registerer
	credMan CredentialManager
}

func NewService(userSvc user.Registerer, credMan CredentialManager) *Service {
	if userSvc == nil {
		panic("missing user.Registerer, parameter must not be nil")
	}
	if credMan == nil {
		panic("missing CredentialManager, parameter must not be nil")
	}
	return &Service{
		userSvc: userSvc,
		credMan: credMan,
	}
}

func NewServiceWithDefaultCredMan(repo Repository, userSvc user.Registerer) *Service {
	if repo == nil {
		panic("missing Repository, parameter must not be nil")
	}
	if userSvc == nil {
		panic("missing user.Registerer, parameter must not be nil")
	}
	return NewService(userSvc, NewManager(repo))
}

// SignIn регистрирует нового пользователя с новым id и добавляет ему логин/пароль
// В случае если такой логин уже есть, то возвращает ошибку.
// id пользователя можно получить только после Login с этой же парой логин/пароль
// ToDo удалить пользователя (компенсация), если ошибка при добавлении кред, так как пользователь и его креды должны быть в БД
// ToDo альтернативно можно проврять, что пользовтель есть, а кред нет, тогда просто добавить креды
// В реальном проекте я бы наплевал на архитектурную красоту в сервисе и сделал бы транзакцию: добавление пользователя+креды
func (s Service) SignIn(ctx context.Context, login, pword string) (usr user.User, err error) {
	// найти пользователя по логину - если есть, то занят
	_, err = s.credMan.GetUser(ctx, login)
	if err == nil {
		return user.User{}, errors2.ErrLoginIsInUseAlready
	}
	// если не занят, то создаем пустого пользователя и регистрируем его
	usr = user.NewUser()
	err = s.userSvc.RegisterNewUser(ctx, usr)
	if err != nil {
		return user.User{}, err
	}
	// создаем креды на пользователя
	err = s.credMan.AddNewUser(ctx, usr, login, pword)
	if err != nil {
		return user.User{}, err
	}
	return usr, nil
}

func (s Service) Login(ctx context.Context, login, pword string) (usr user.User, err error) {
	return s.credMan.AuthenticateUser(ctx, login, pword)
}
