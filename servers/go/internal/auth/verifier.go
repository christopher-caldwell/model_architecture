package auth

import (
	"context"
	"errors"
)

var ErrInvalidToken = errors.New("invalid token")

type Verifier interface {
	Verify(ctx context.Context, token string) (*Claims, error)
}
