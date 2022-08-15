package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/app"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/user"
	errors2 "github.com/UndeadDemidov/ya-pr-diploma/internal/errors"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/presenter/http/middleware"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/presenter/http/utils"
)

var (
	// ToDo - заменить на нормальные ответы в json!
	ErrInvalidContentType   = fmt.Errorf("set header key %v to proper value", utils.ContentTypeKey)
	ErrProperJSONIsExpected = errors.New("proper JSON is expected, read task description carefully")
)

type Auth struct {
	auth     app.Authenticator
	sessions *middleware.Sessions
}

func NewAuth(auth app.Authenticator, sessions *middleware.Sessions) *Auth {
	if auth == nil {
		panic("missing app.Authenticator, parameter must not be nil")
	}
	if sessions == nil {
		panic("missing *middleware.Sessions, parameter must not be nil")
	}
	return &Auth{auth: auth, sessions: sessions}
}

// RegisterUser
// POST /api/user/register
func (a *Auth) RegisterUser(w http.ResponseWriter, r *http.Request) {
	req := authRequest{}
	err := req.Read(r)
	if err != nil {
		utils.ServerError(w, ErrInvalidContentType, http.StatusBadRequest)
		return
	}

	usr, err := a.auth.SignIn(r.Context(), req.Login, req.Password)
	if err != nil {
		if errors.Is(err, errors2.ErrLoginIsInUseAlready) {
			utils.ServerError(w, errors2.ErrLoginIsInUseAlready, http.StatusConflict)
			return
		}
		utils.InternalServerError(w, err)
		return
	}
	// Можно было и JWT поюзать, но решил для практики поизобретать велосипеды в отпуске,
	// чтобы не обнулиться в дно за месяц академа
	cookie := middleware.NewSessionSignedCookie(a.sessions.AddNewSession(usr))
	cookie.Set(w)
	w.WriteHeader(http.StatusOK)
}

// LoginUser
// POST /api/user/login
func (a *Auth) LoginUser(w http.ResponseWriter, r *http.Request) {
	req := authRequest{}
	err := req.Read(r)
	if err != nil {
		utils.ServerError(w, ErrInvalidContentType, http.StatusBadRequest)
		return
	}

	usr, err := a.auth.Login(r.Context(), req.Login, req.Password)
	if err != nil {
		if errors.Is(err, errors2.ErrPairLoginPwordIsNotExist) {
			utils.ServerError(w, errors2.ErrPairLoginPwordIsNotExist, http.StatusUnauthorized)
			return
		}
		utils.InternalServerError(w, err)
		return
	}
	// Можно было и JWT поюзать, но решил для практики поизобретать велосипеды в отпуске,
	// чтобы не обнулиться в дно за месяц академа
	cookie := middleware.NewSessionSignedCookie(a.sessions.AddNewSession(usr))
	cookie.Set(w)
	w.WriteHeader(http.StatusOK)
}

// {
//	"login": "<login>",
//	"password": "<password>"
// }
type authRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (ar *authRequest) Read(r *http.Request) error {
	if r.Header.Get(utils.ContentTypeKey) != utils.ContentTypeJSON {
		return ErrInvalidContentType
	}
	err := json.NewDecoder(r.Body).Decode(ar)
	if err != nil {
		// ToDo wrapper?
		return ErrProperJSONIsExpected
	}
	return nil
}

// GetUserFromContext возвращает сохраненного в контексте пользователя
func GetUserFromContext(ctx context.Context) user.User {
	if ctx == nil {
		return user.User{}
	}
	if userID, ok := ctx.Value(middleware.ContextUserIDKey).(string); ok {
		return user.User{ID: userID}
	}
	return user.User{}
}
