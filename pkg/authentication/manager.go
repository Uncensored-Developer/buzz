package authentication

import (
	"fmt"
	"github.com/form3tech-oss/jwt-go"
	"github.com/pkg/errors"
	"time"
)

type ITokenManager interface {
	NewToken(userId int64, ttl time.Duration) (string, error)
	Parse(token string) (int64, error)
}

type Manager struct {
	signingKey string
}

func NewManager(signingKey string) (*Manager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}
	return &Manager{signingKey: signingKey}, nil
}

func (m *Manager) NewToken(userId int64, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"userId": userId})
	tokenStr, err := token.SignedString([]byte(m.signingKey))
	if err != nil {
		return "", errors.Wrap(err, "token sign failed")
	}
	return tokenStr, nil
}

func (m *Manager) Parse(token string) (int64, error) {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			msg := fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"])
			return nil, errors.New(msg)
		}
		return []byte(m.signingKey), nil
	})
	if err != nil {
		return 0, errors.Wrap(err, "token parse failed")
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.Wrap(err, "token parse failed")
	}
	return int64(claims["userId"].(float64)), nil
}
