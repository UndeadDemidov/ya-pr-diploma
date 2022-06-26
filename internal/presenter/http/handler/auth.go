package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/app"
	errors2 "github.com/UndeadDemidov/ya-pr-diploma/internal/errors"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/presenter/http/utils"
)

var (
	// ToDo - заменить на нормальные ответы в json!
	ErrInvalidContentType   = fmt.Errorf("set header value %v to %v", utils.ContentTypeKey, utils.ContentTypeJSON)
	ErrProperJSONIsExpected = errors.New("proper JSON is expected, read task description carefully")
)

type Auth struct {
	auth app.Authenticator
}

func NewAuth(auth app.Authenticator) *Auth {
	return &Auth{auth: auth}
}

// POST /api/user/register
func (a Auth) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get(utils.ContentTypeKey) != utils.ContentTypeJSON {
		utils.ServerError(w, ErrInvalidContentType, http.StatusBadRequest)
		return
	}
	req := authRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.ServerError(w, ErrProperJSONIsExpected, http.StatusBadRequest)
		return
	}

	err = a.auth.SignIn(r.Context(), req.Login, req.Password)
	if err != nil {
		if errors.Is(err, errors2.ErrLoginIsInUseAlready) {
			utils.ServerError(w, errors2.ErrLoginIsInUseAlready, http.StatusConflict)
			return
		}
		utils.InternalServerError(w, err)
		return
	}
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
