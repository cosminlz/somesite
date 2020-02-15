package auth

import (
	"fmt"
	"time"

	"cabhelp.ro/backend/internal/model"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

var jwtKey = []byte("my_secret_key")

// var accessTokenDuration = time.Duration(30) * time.Minute
var accessTokenDuration = time.Duration(1) * time.Second

// var refreshTokenDuration = time.Duration(30*24) * time.Hour
var refreshTokenDuration = time.Duration(1) * time.Second

// Claims ...
type Claims struct {
	UserID model.UserID `json:"userID"`
	jwt.StandardClaims
}

// Tokens ...
type Tokens struct {
	AccessToken           string `json:"accessToken,omitempty"`
	AccessTokenExpiresAt  int64  `json:"accessTokenExpiresAt,omitempty"`
	RefreshToken          string `json:"refreshToken,omitempty"`
	RefreshTokenExpiresAt int64  `json:"-"`
}

// IssueToken returns a new access and refresh tokens
func IssueToken(principal model.Principal) (*Tokens, error) {

	if principal.UserID == model.NilUserID {
		return nil, errors.New("invalid principal")
	}

	accessToken, accessTokenExpiresAt, err := generateToken(principal, accessTokenDuration)
	if err != nil {
		return nil, err
	}
	refreshToken, refreshTokenExpiresAt, err := generateToken(principal, refreshTokenDuration)
	if err != nil {
		return nil, err
	}

	tokens := &Tokens{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
	}

	return tokens, nil
}

func generateToken(principal model.Principal, duration time.Duration) (string, int64, error) {
	now := time.Now()

	tokenClaims := &Claims{
		UserID: principal.UserID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(duration).Unix(),
		},
	}

	tokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS512, tokenClaims)

	token, err := tokenWithClaims.SignedString(jwtKey)
	if err != nil {
		return "", 0, err
	}

	return token, tokenClaims.ExpiresAt, nil
}

func VerifyToken(tokenString string) (*model.Principal, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	fmt.Println(token)
	fmt.Println(claims)

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, err
		}
		return nil, err
	}

	principal := &model.Principal{
		UserID: claims.UserID,
	}

	if !token.Valid {
		return principal, err
	}

	return principal, nil
}
