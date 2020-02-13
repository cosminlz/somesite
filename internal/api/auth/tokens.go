package auth

import (
	"time"

	"cabhelp.ro/backend/internal/model"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

var jwtKey = []byte("my_secret_key")
var accessTokenDuration = time.Duration(30) * time.Minute
var refreshTokenDuration = time.Duration(30*24) * time.Hour

/*
type Tokens interface {
	IssueToken(principal model.Principal) (string, error)
	//Verify(token string) (*model.Principal, error)
}

type tokens struct {
	key             []byte
	duration        time.Duration
	beforeTolerance time.Duration
	signingMethod   jwt.SigningMethod
}
*/

// Claims ...
type Claims struct {
	UserID model.UserID `json:"userID"`
	jwt.StandardClaims
}

/*
// NewTokens creates a new Tokens object
func NewTokens() Tokens {

	tokenDuration := time.Duration(tokenDurationHours) * time.Hour
	return &tokens{
		key:           jwtKey,
		duration:      tokenDuration,
		signingMethod: jwt.SigningMethodHS512,
	}
}
*/

// Tokens ...
type Tokens struct {
	AccessToken  string `json:"accessToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

// IssueToken returns a new token
func IssueToken(principal model.Principal) (*Tokens, error) {

	if principal.UserID == model.NilUserID {
		return nil, errors.New("invalid principal")
	}

	now := time.Now()

	claims := &Claims{
		UserID: principal.UserID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(accessTokenDuration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	accessToken, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, err
	}

	tokens := &Tokens{
		AccessToken: accessToken,
	}

	return tokens, nil
}

func VerifyToken(accessToken string) (*model.Principal, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, err
		}
		return nil, err
	}
	if !token.Valid {
		return nil, err
	}

	principal := &model.Principal{
		UserID: claims.UserID,
	}

	return principal, nil
}
