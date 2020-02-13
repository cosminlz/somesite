package model

// Credentials represents the set of creds for login
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
