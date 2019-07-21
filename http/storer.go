package http

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	"github.com/volatiletech/authboss"
)

var nextUserID int

// User struct for authboss
type User struct {
	ID int

	// Non-authboss related field
	FirstName string
	LastName string
	Affiliation string

	// Auth
	Email    string
	Username string
	Password string
}

// This pattern is useful in real code to ensure that
// we've got the right interfaces implemented.
var (
	assertUser   = &User{}
	assertStorer = &MemStorer{}

	_ authboss.User            = assertUser
	_ authboss.AuthableUser    = assertUser

	_ authboss.CreatingServerStorer    = assertStorer
	_ authboss.RememberingServerStorer = assertStorer
)

// PutPID into user
func (u *User) PutPID(pid string) { u.Username = pid }

// PutPassword into user
func (u *User) PutPassword(password string) { u.Password = password }

// PutEmail into user
func (u *User) PutEmail(email string) { u.Email = email }

// PutUsername into user
func (u *User) PutUsername(username string) { u.Username = username }

// GetPID from user
func (u User) GetPID() string { return u.Username }

// GetPassword from user
func (u User) GetPassword() string { return u.Password }

// GetEmail from user
func (u User) GetEmail() string { return u.Email }


// MemStorer stores users in memory
// Indexed by Username => must be unique
// TODO: verify that username is unique
type MemStorer struct {
	Users  map[string]User
	Tokens map[string][]string
}

// NewMemStorer constructor
func NewMemStorer() *MemStorer {
	return &MemStorer{
		Users: map[string]User{
			"as29": User{
				ID:                 1,
				FirstName:          "Abhishek",
				LastName:	    "Srivastava",
				Username:           "as29",
				Password:           "$2a$10$XtW/BrS5HeYIuOCXYe8DFuInetDMdaarMUJEOg/VA/JAIDgw3l4aG", // pass = 1234
				Email:              "as29@illinois.edu",
			},
		},
		Tokens: make(map[string][]string),
	}
}

// Save the user
func (m MemStorer) Save(ctx context.Context, user authboss.User) error {
	u := user.(*User)
	m.Users[u.Username] = *u

	//debugln("Saved user:", u.FirstName)
	return nil
}

// Load the user
func (m MemStorer) Load(ctx context.Context, key string) (user authboss.User, err error) {
	u, ok := m.Users[key]
	if !ok {
		return nil, authboss.ErrUserNotFound
	}

	//debugln("Loaded user:", u.FirstName)
	return &u, nil
}

// New user creation
func (m MemStorer) New(ctx context.Context) authboss.User {
	return &User{}
}

// Create the user
func (m MemStorer) Create(ctx context.Context, user authboss.User) error {
	u := user.(*User)

	if _, ok := m.Users[u.Username]; ok {
		return authboss.ErrUserFound
	}

	//debugln("Created new user:", u.FirstName)
	m.Users[u.Username] = *u
	return nil
}

// AddRememberToken to a user
func (m MemStorer) AddRememberToken(ctx context.Context, pid, token string) error {
	m.Tokens[pid] = append(m.Tokens[pid], token)
	//debugf("Adding rm token to %s: %s\n", pid, token)
	spew.Dump(m.Tokens)
	return nil
}

// DelRememberTokens removes all tokens for the given pid
func (m MemStorer) DelRememberTokens(ctx context.Context, pid string) error {
	delete(m.Tokens, pid)
	//debugln("Deleting rm tokens from:", pid)
	spew.Dump(m.Tokens)
	return nil
}

// UseRememberToken finds the pid-token pair and deletes it.
// If the token could not be found return ErrTokenNotFound
func (m MemStorer) UseRememberToken(ctx context.Context, pid, token string) error {
	tokens, ok := m.Tokens[pid]
	if !ok {
		//debugln("Failed to find rm tokens for:", pid)
		return authboss.ErrTokenNotFound
	}

	for i, tok := range tokens {
		if tok == token {
			tokens[len(tokens)-1] = tokens[i]
			m.Tokens[pid] = tokens[:len(tokens)-1]
			//debugf("Used remember for %s: %s\n", pid, token)
			return nil
		}
	}

	return authboss.ErrTokenNotFound
}