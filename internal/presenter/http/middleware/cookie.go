package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"

	errors2 "github.com/UndeadDemidov/ya-pr-diploma/internal/errors"
	"github.com/UndeadDemidov/ya-pr-diploma/internal/presenter/http/utils"
)

const (
	sessionIDCookie = "GopherMartSessionID"
	salt            = "secret key" // Можно использовать IP или еще что-то присущее конкретному пользователю/машине
	saltStartIdx    = 4
	saltEndIdx      = 9
	maxAge          = 60 * 60 * 24 * 180
)

var (
	ErrSignedCookieInvalidValueOrUnsigned = errors.New("invalid cookie value or it is unsigned")
	ErrSignedCookieInvalidSign            = errors.New("invalid sign")
	ErrSignedCookieSaltNotSetProperly     = errors.New("SaltStartIdx and SaltEndIdx must be set properly")
	ContextUserIDKey                      = LocalContext(sessionIDCookie)
)

type LocalContext string

func SessionsCookie(sessions *Sessions) func(next http.Handler) http.Handler {
	if sessions == nil {
		panic("missing *Sessions, parameter must not be nil")
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// получить из куки id сессии
			token, err := getSessionTokenFromCookie(sessionIDCookie, r)
			// если сессии нет - прерываем работу
			if err != nil {
				utils.ServerError(w, err, http.StatusUnauthorized)
				return
			}
			// если сессия есть - проверяем валидность
			// если не валидна - прерываем работу
			if sessions.IsExpired(token) {
				utils.ServerError(w, errors2.ErrSessionIsExpired, http.StatusUnauthorized)
				return
			}
			// если сессия валидна - записываем id пользователя в контекст
			ctx := context.WithValue(r.Context(), ContextUserIDKey, sessions.GetReference(token))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func getSessionTokenFromCookie(name string, r *http.Request) (token SessionToken, err error) {
	c, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	cookie := NewSessionSignedCookie("")
	cookie.Cookie = c
	err = cookie.DetachSign()
	if err != nil {
		return "", err
	}
	return SessionToken(cookie.BaseValue), nil
}

type SignedCookie struct {
	*http.Cookie
	SaltStartIdx uint
	SaltEndIdx   uint
	key          []byte
	sign         []byte
	BaseValue    string
}

func NewSessionSignedCookie(val SessionToken) (sc SignedCookie) {
	return NewSignedCookie("/", sessionIDCookie, string(val), maxAge, saltStartIdx, saltEndIdx)
}

func NewSignedCookie(path, name, val string, maxAge int, saltStartIdx, saltEndIdx uint) (sc SignedCookie) {
	sc = SignedCookie{
		Cookie: &http.Cookie{
			Path:   path,
			Name:   name,
			Value:  val,
			MaxAge: maxAge,
		},
		SaltStartIdx: saltStartIdx,
		SaltEndIdx:   saltEndIdx,
	}

	sc.AttachSign()
	return sc
}

func (sc *SignedCookie) AttachSign() {
	sc.BaseValue = sc.Value
	if len(sc.key) == 0 {
		sc.RecalcKey()
	}
	sc.sign = sc.calcSign()
	sc.Value = fmt.Sprintf("%s|%s", sc.Value, hex.EncodeToString(sc.sign))
}

func (sc *SignedCookie) calcSign() []byte {
	h := hmac.New(sha256.New, sc.key)
	h.Write([]byte(sc.BaseValue))
	return h.Sum(nil)
}

func (sc *SignedCookie) RecalcKey() {
	if sc.SaltStartIdx == 0 || sc.SaltEndIdx == 0 ||
		sc.SaltEndIdx < sc.SaltStartIdx || sc.SaltEndIdx > uint(len(sc.BaseValue)) {
		panic(ErrSignedCookieSaltNotSetProperly)
	}

	var key = []byte(salt)
	key = append(key, []byte(sc.BaseValue)[sc.SaltStartIdx:sc.SaltEndIdx]...)
	sc.key = key
}

func (sc *SignedCookie) DetachSign() (err error) {
	ss := strings.Split(sc.Value, "|")
	if len(ss) < 2 {
		return ErrSignedCookieInvalidValueOrUnsigned
	}
	sc.BaseValue = ss[0]
	sc.RecalcKey()

	sign := ss[1]
	if hex.EncodeToString(sc.calcSign()) != sign {
		return ErrSignedCookieInvalidSign
	}

	return nil
}

func (sc *SignedCookie) Set(w http.ResponseWriter) {
	http.SetCookie(w, sc.Cookie)
}
