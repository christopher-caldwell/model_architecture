package auth

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

const jwtAudience = "ops.craftcode.solutions"

type JwtVerifier struct {
	secret []byte
}

func NewJwtVerifier(secret string) *JwtVerifier {
	return &JwtVerifier{secret: []byte(secret)}
}

type jwtClaims struct {
	jwt.RegisteredClaims
}

func (v *JwtVerifier) Verify(_ context.Context, token string) (*Claims, error) {
	parsed, err := jwt.ParseWithClaims(token, &jwtClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return v.secret, nil
	}, jwt.WithAudience(jwtAudience))
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidToken, err)
	}

	c, ok := parsed.Claims.(*jwtClaims)
	if !ok || !parsed.Valid {
		return nil, ErrInvalidToken
	}

	var exp int64
	if c.ExpiresAt != nil {
		exp = c.ExpiresAt.Unix()
	}
	return &Claims{
		Sub: c.Subject,
		Exp: exp,
	}, nil
}
