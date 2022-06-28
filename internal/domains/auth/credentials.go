package auth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"hash"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/domains/user"
	_ "github.com/golang/mock/mockgen/model"
)

const salt = "gPhRmRt"

//go:generate mockgen -destination=./mocks/mock_credentials.go . Repository

type Repository interface {
	Create(ctx context.Context, user user.User, login, pword string) error
	Read(ctx context.Context, login string) (usr user.User, err error)
	ReadWithPassword(ctx context.Context, login, pword string) (usr user.User, err error)
}

var _ CredentialManager = (*Manager)(nil)

// Никакого хранения в памяти кредов не делаем. Чем меньше мест, где храняться логины/пароли,
// тем меньше мест нужно защищать от хакеров - даже если это дополнительное место - память.
// Имеет смысл кешировать столько с очень большой нагрузкой по аутентификации.
type Manager struct {
	repo   Repository
	hasher *SaltedHash
}

func NewManager(repo Repository) *Manager {
	return &Manager{repo: repo, hasher: NewSaltedHashWithDefaultMixer(sha256.New(), []byte(salt))}
}

func (man *Manager) AddNewUser(ctx context.Context, usr user.User, login, pword string) error {
	return man.repo.Create(ctx, usr, login, man.getHashedPassword(pword))
}

func (man *Manager) GetUser(ctx context.Context, login string) (usr user.User, err error) {
	usr, err = man.repo.Read(ctx, login)
	if err != nil {
		return user.User{}, err
	}
	return usr, nil
}

func (man *Manager) AuthenticateUser(ctx context.Context, login, pword string) (usr user.User, err error) {
	usr, err = man.repo.ReadWithPassword(ctx, login, man.getHashedPassword(pword))
	if err != nil {
		return user.User{}, err
	}
	return usr, nil
}

func (man *Manager) getHashedPassword(pword string) string {
	return base64.URLEncoding.EncodeToString(man.hasher.Sum([]byte(pword)))
}

type SaltedHash struct {
	hash.Hash
	salt  []byte
	mixer func([]byte, []byte) []byte
}

func NewSaltedHash(h hash.Hash, salt []byte, mixer func([]byte, []byte) []byte) *SaltedHash {
	return &SaltedHash{Hash: h, salt: salt, mixer: mixer}
}

func NewSaltedHashWithDefaultMixer(h hash.Hash, salt []byte) *SaltedHash {
	m := func(body []byte, salt []byte) []byte {
		return append(body, salt...)
	}
	return NewSaltedHash(h, salt, m)
}

func (h SaltedHash) Sum(b []byte) []byte {
	return h.Hash.Sum(h.mixer(b, h.salt))
}
