package idtoken

import (
	"context"
	"errors"
	"google.golang.org/api/idtoken"
)

type Parser struct {
}

var ErrInvalidToken = errors.New("invalid token")

func (p *Parser) Parse(ctx context.Context, token string) (*TokenInfo, error) {
	payload, err := idtoken.Validate(ctx, token, "")
	if err != nil {
		return nil, ErrInvalidToken
	}

	name, ok := payload.Claims["name"].(string)
	if !ok {
		return nil, errors.New("can't retrieve name")
	}

	email, ok := payload.Claims["email"].(string)
	if !ok {
		return nil, errors.New("can't retrieve email")
	}

	url, ok := payload.Claims["picture"].(string)
	if !ok {
		return nil, errors.New("can't retrieve picture url")
	}

	return &TokenInfo{
		Name:    name,
		Email:   email,
		Picture: url,
	}, nil
}
