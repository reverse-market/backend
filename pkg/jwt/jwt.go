package jwt

import (
	"errors"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	ErrInvalidToken   = errors.New("invalid token")
	ErrExpiredToken   = errors.New("expired token")
	ErrInvalidSubject = errors.New("invalid subject")
)

type Manager struct {
	Secret     string
	Expiration time.Duration
}

func (m *Manager) CreateToken(id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(m.Expiration).Unix(),
		Subject:   strconv.Itoa(id),
	})

	return token.SignedString([]byte(m.Secret))
}

func (m *Manager) GetIdFromToken(token string) (int, error) {
	parsed, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.Secret), nil
	})
	if err != nil {
		return 0, ErrInvalidToken
	}

	claims, ok := parsed.Claims.(*jwt.StandardClaims)
	if !ok {
		return 0, ErrInvalidToken
	}

	if err := claims.Valid(); err != nil {
		var jwtErr *jwt.ValidationError
		if errors.As(err, jwtErr); jwtErr.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return 0, ErrExpiredToken
		}
		return 0, ErrInvalidToken
	}

	id, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return 0, ErrInvalidSubject
	}

	return id, nil
}
