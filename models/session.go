package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/kolioDev/after_life/helpers"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
)

type Session struct {
	ID        UUID      `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	UserID uuid.UUID `json:"user_id" db:"user_id"`

	UniqueToken nulls.String `json:"unique_token" db:"unique_token"` //used only in registration/login

	RefreshToken     string `json:"refresh_token" db:"-"`
	RefreshTokenHash []byte `json:"-" db:"refresh_token"`
}

func (s *Session) TableName() string {
	return "sessions"
}

// String is not required by pop and may be deleted
func (s Session) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Sessions is not required by pop and may be deleted
type Sessions []Session

// String is not required by pop and may be deleted
func (s Sessions) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (s *Session) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (s *Session) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (s *Session) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (s *Session) Create(tx *pop.Connection, u *User, protected bool) error {
	if u.ID == uuid.Nil {
		return errors.New("user cannot be empty")
	}

	s.UserID = u.ID
	if !protected {
		s.UniqueToken = nulls.NewString(helpers.RandString(25))
	}
	s.RefreshToken = helpers.RandString(35)
	var err error
	s.RefreshTokenHash, err = bcrypt.GenerateFromPassword([]byte(s.RefreshToken), bcrypt.DefaultCost)
	if err != nil {
		return errors.WithStack(err)
	}

	return tx.Create(s)
}

func (s *Session) ResetUniqueToken(tx *pop.Connection) error {
	if err := tx.Find(s, s.ID); err != nil {
		return errors.WithStack(err)
	}

	sUserID := s.UserID

	if err := tx.Destroy(s); err != nil {
		return errors.WithStack(err)
	}
	s.UniqueToken.Valid = false
	s.UniqueToken.String = ""

	return s.Create(tx, &User{ID: sUserID}, true)
}
