package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// NewToken generates a new signed JWT token.
func NewToken(user User, secret []byte, expiration time.Duration) (
	token *Token, err error) {
	jwt, err := genToken(tokenParams{
		User:      user,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(expiration),
		Secret:    secret,
	})
	if err != nil {
		err = fmt.Errorf("%v; %w", err, ErrInternalError)
		return
	}
	token = &Token{
		JWT: jwt,
	}
	return
}

// Token defines the token model.
type Token struct {
	JWT string
}

type tokenParams struct {
	User      User
	IssuedAt  time.Time
	ExpiredAt time.Time
	Secret    []byte
}

// TokenClaims defines the token claims struct.
type TokenClaims struct {
	Email     string `json:"email"`
	Org       string `json:"org"`
	IssuedAt  int64  `json:"issued_at"`
	ExpiredAt int64  `json:"expired_at"`
}

// Valid implements the JWT interface.
func (c *TokenClaims) Valid() error {
	if c.ExpiredAt < time.Now().Unix() {
		return ErrTokenExpired
	}
	return nil
}

func genToken(params tokenParams) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		Email:     params.User.Email,
		Org:       params.User.Organization,
		IssuedAt:  params.IssuedAt.Unix(),
		ExpiredAt: params.ExpiredAt.Unix(),
	})
	return jwtToken.SignedString(params.Secret)
}

func verifyToken(tokenStr string, secret []byte) (
	claims *TokenClaims, err error) {
	claims = &TokenClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims,
		func(token *jwt.Token) (interface{}, error) {
			// validate the alg
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v",
					token.Header["alg"])
			}
			return secret, nil
		})

	if err != nil {
		err = fmt.Errorf("%v; %w", err, ErrInvalidToken)
		return
	}
	if !token.Valid {
		err = ErrInvalidToken
		return
	}
	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		err = ErrInvalidToken
		return
	}
	return
}
