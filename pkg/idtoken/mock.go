package idtoken

import (
	"context"
)

type MockParser struct {
}

const validToken = "VALID_TOKEN"

func (p *MockParser) Parse(_ context.Context, token string) (*TokenInfo, error) {
	if token == validToken {
		return &TokenInfo{
			Name:    "Ivanov Ivan",
			Email:   "ivan@mail.ru",
			Picture: "",
		}, nil
	} else {
		return nil, ErrInvalidToken
	}
}
