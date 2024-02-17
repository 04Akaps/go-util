package paseto

import (
	auth "github.com/04Akaps/go-util/auth/proto"
	"github.com/o1egl/paseto"
	"math/rand"
)

type PasetoMaker struct {
	pt  *paseto.V2
	key []byte
}

type PasetoInterface interface {
	CreateToken(*auth.AuthData) (string, error)
	VerifyToken(token string) error
}

func NewPasetoMaker(key string) PasetoInterface {
	return &PasetoMaker{
		pt:  paseto.NewV2(),
		key: []byte(key),
	}
}

func (m *PasetoMaker) CreateToken(auth *auth.AuthData) (string, error) {
	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)
	return m.pt.Encrypt(m.key, auth, nil)
}

func (m *PasetoMaker) VerifyToken(token string) error {
	payload := &auth.AuthData{}
	return m.pt.Decrypt(token, m.key, payload, nil)
}
