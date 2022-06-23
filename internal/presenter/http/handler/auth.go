package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/presenter/http/utils"
)

//go:generate mockgen -destination=./mocks/mock_auth.go . CredentialManager,CredentialValidator

var (
	// ToDo - заменить на нормальные ответы в json!
	ErrInvalidContentType   = fmt.Errorf("set header value %v to %v", utils.ContentTypeKey, utils.ContentTypeJSON)
	ErrProperJSONIsExpected = errors.New("proper JSON is expected, read task description carefully")
	ErrLoginIsInUseAlready  = errors.New("given login is in use already")
)

type CredentialValidator interface {
	IsValid()
}

type CredentialManager interface {
	New(user, pword string) (cred CredentialValidator, ok bool)
	Add(cred CredentialValidator) error
	Get(user string) (cred CredentialValidator, err error)
	// Remove(credential CredentialValidator) error
}

type Auth struct {
	creds CredentialManager
}

func NewAuth(creds CredentialManager) *Auth {
	return &Auth{creds: creds}
}

// POST /api/user/register
func (a Auth) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get(utils.ContentTypeKey) != utils.ContentTypeJSON {
		http.Error(w, ErrInvalidContentType.Error(), http.StatusBadRequest)
		return
	}
	req := authRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.ServerError(w, ErrProperJSONIsExpected, http.StatusBadRequest)
		return
	}
	if cred, ok := a.creds.New(req.Login, req.Password); !ok {
		utils.ServerError(w, ErrLoginIsInUseAlready, http.StatusConflict)
		return
	} else {
		err = a.creds.Add(cred)
		if err != nil {
			utils.InternalServerError(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// {
//	"login": "<login>",
//	"password": "<password>"
// }
type authRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
