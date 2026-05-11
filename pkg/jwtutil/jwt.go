package jwtutil

import (
	"errors"
	"fmt"
	"time"

	"reddit/internal/config"
	"reddit/internal/domain/auth"
	"reddit/internal/domain/user"

	jwt "github.com/golang-jwt/jwt/v5"
)

type userClaims struct {
	User userInfo `json:"user"`
	jwt.RegisteredClaims
}
type userInfo struct {
	UserID   string `json:"id"`
	Username string `json:"username"`
}

type JWT struct {
	secret string
	expire time.Duration
}

func New(cfg config.JWT) (*JWT, error) {
	if cfg.Secret == "" {
		return nil, auth.ErrMissingSecret
	}
	if cfg.Expire <= 0 {
		return nil, auth.ErrInvalidTTL
	}
	return &JWT{secret: cfg.Secret, expire: cfg.Expire}, nil
}

func (j *JWT) GenerateToken(user user.User) (string, error) {
	if err := j.ready(); err != nil {
		return "", err
	}
	if user.ID == "" {
		return "", auth.ErrInvalidUserID
	}
	now := time.Now()
	claims := userClaims{
		User: userInfo{
			UserID:   user.ID,
			Username: user.Username,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.expire)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", fmt.Errorf("%w: %v", auth.ErrSignFailed, err)
	}
	return signed, nil
}

func (j *JWT) ValidateToken(raw string) (string, error) {
	token, claims, err := j.parse(raw)
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", auth.ErrInvalidToken
	}
	if claims.User.UserID == "" {
		return "", auth.ErrInvalidToken
	}
	return claims.User.UserID, nil
}

func (j *JWT) parse(tokenStr string) (*jwt.Token, *userClaims, error) {
	if err := j.ready(); err != nil {
		return nil, nil, err
	}
	claims := &userClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method == nil || token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, auth.ErrInvalidToken
		}
		return []byte(j.secret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, nil, auth.ErrExpiredToken
		}
		return nil, nil, fmt.Errorf("%w: %v", auth.ErrInvalidToken, err)
	}
	return token, claims, nil
}

func (j *JWT) ready() error {
	if j.secret == "" {
		return auth.ErrNotInitialized
	}
	if j.expire <= 0 {
		return auth.ErrInvalidTTL
	}
	return nil
}

var (
	_ auth.TokenIssuer    = (*JWT)(nil)
	_ auth.TokenValidator = (*JWT)(nil)
)
