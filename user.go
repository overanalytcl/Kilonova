package kilonova

import (
	"errors"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserBrief struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Admin    bool   `json:"admin"`
	Proposer bool   `json:"proposer"`
	Bio      string `json:"bio,omitempty"`
}

type UserFull struct {
	UserBrief
	Email             string    `json:"email,omitempty"`
	VerifiedEmail     bool      `json:"verified_email"`
	PreferredLanguage string    `json:"preferred_language"`
	EmailVerifResent  time.Time `json:"-"`
	CreatedAt         time.Time `json:"created_at"`
}

func (uf *UserFull) Brief() *UserBrief {
	if uf == nil {
		return nil
	}
	return &uf.UserBrief
}

// UserFilter is the struct with all filterable fields on the user
// It also provides a Limit and Offset field, for pagination
type UserFilter struct {
	ID  *int  `json:"id"`
	IDs []int `json:"ids"`

	// Name is case insensitive
	Name  *string `json:"name"`
	Email *string `json:"email"`

	Admin    *bool `json:"admin"`
	Proposer *bool `json:"proposer"`

	// For registrations
	ContestID *int `json:"contest_id"`

	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// UserUpdate is the struct with updatable fields that can be easily changed.
// Stuff like admin and proposer should be updated through their dedicated
// SudoAPI methods
type UserUpdate struct {
	//Name    *string `json:"name"`

	Bio *string `json:"bio"`

	PreferredLanguage string `json:"-"`
}

// UserFullUpdate is the struct with all updatable fields on the user
type UserFullUpdate struct {
	UserUpdate

	Email *string `json:"email"`

	Admin    *bool `json:"admin"`
	Proposer *bool `json:"proposer"`

	VerifiedEmail    *bool      `json:"verified_email"`
	EmailVerifSentAt *time.Time `json:"-"`
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), err
}

func CheckPwdHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false
	}
	if err != nil {
		zap.S().Warn(err)
		return false
	}
	return true
}
