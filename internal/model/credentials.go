package model

import "fmt"

// Credentials represents the set of creds for login
type Credentials struct {
	SessionData

	Email    string `json:"email"`
	Password string `json:"password"`
}

// Principal is an authenticated entity
type Principal struct {
	UserID UserID `json:"userID,omitempty"`
}

var NilPrincipal Principal

func (p Principal) String() string {
	if p.UserID != "" {
		return fmt.Sprintf("UserID[%s]", p.UserID)
	}
	return "(nil)"
}
