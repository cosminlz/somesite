package auth

import (
	"crypto/sha512"
	"time"

	"cabhelp.ro/backend/internal/model"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

var jwtKey = []byte("my_secret_key")
var tokenDurationHours = 30 * 24

// Tokens
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

// Claims
type Claims struct {
	UserID model.UserID `json:"userID"`
	jwt.StandardClaims
}

func NewTokens() Tokens {
	hasher := sha512.New()

	if _, err := hasher.Write([]byte(jwtKey)); err != nil {
		panic(err)
	}

	tokenDuration := time.Duration(tokenDurationHours) * time.Hour

	return &tokens{
		key:             hasher.Sum(nil),
		duration:        tokenDuration,
		beforeTolerance: -2 * time.Minute,
		signingMethod:   jwt.SigningMethodHS512,
	}
}

func (t *tokens) IssueToken(principal model.Principal) (string, error) {

	if principal.UserID == model.NilUserID {
		return "", errors.New("invalid principal")
	}

	now := time.Now()

	claims := &Claims{
		UserID: principal.UserID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			NotBefore: now.Add(t.beforeTolerance).Unix(),
			ExpiresAt: now.Add(t.duration).Unix(),
		},
	}

	token := jwt.NewWithClaims(t.signingMethod, claims)

	return token.SignedString(t.key)
}

// func (t *tokens) Verify(token string) (*model.Principal, error) {
// }
