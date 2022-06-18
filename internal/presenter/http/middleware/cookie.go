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

	utils "github.com/UndeadDemidov/ya-pr-diploma/internal/presenter/http"
)

const (
	userIDCookie = "Gopher-martUserID"
	secretKey    = "secret key" // Такие вещи выносятся в Vault.
	saltStartIdx = 4
	saltEndIdx   = 9
	maxAge       = 60 * 60 * 24 * 180
)

var (
	ErrSignedCookieInvalidValueOrUnsigned = errors.New("invalid cookie value or it is unsigned")
	ErrSignedCookieInvalidSign            = errors.New("invalid sign")
	ErrSignedCookieSaltNotSetProperly     = errors.New("SaltStartIdx and SaltEndIdx must be set properly")
	ContextUserIDKey                      = LocalContext(userIDCookie)
)

type LocalContext string

func UserCookie(next http.Handler) http.Handler {
	middleware := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user, err := getUserID(r)
		if err != nil {
			utils.InternalServerError(w, err)
		}
		ctx = context.WithValue(ctx, ContextUserIDKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(middleware)
}

func getUserID(r *http.Request) (userID string, err error) {
	// получить куку пользователя
	c, err := r.Cookie(userIDCookie)
	// куки нет
	if errors.Is(err, http.ErrNoCookie) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	// кука есть
	cookie := NewUserIDSignedCookie("")
	cookie.Cookie = c
	err = cookie.DetachSign()
	switch {
	case err == nil: // кука подписана верно
		return cookie.BaseValue, nil
	case errors.Is(err, ErrSignedCookieInvalidSign): // кука подписана неверно
		return "", nil
	}
	return "", err
}

// GetUserID возвращает сохраненный в контексте куку userIDCookie
func GetUserID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if userID, ok := ctx.Value(ContextUserIDKey).(string); ok {
		return userID
	}
	return ""
}

type SignedCookie struct {
	*http.Cookie
	SaltStartIdx uint
	SaltEndIdx   uint
	key          []byte
	sign         []byte
	BaseValue    string
}

func NewUserIDSignedCookie(usedID string) (sc SignedCookie) {
	sc = SignedCookie{
		Cookie: &http.Cookie{
			Path:   "/",
			Name:   userIDCookie,
			Value:  usedID,
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

	var key = []byte(secretKey)
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
