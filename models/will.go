package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
)

type Will struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	Title      string `json:"title" db:"title"`
	Importance int    `json:"importance" db:"importance"`

	Instruction Instructions `json:"instructions" db:"-"`
	Pictures    Pictures     `json:"pictures" db:"-"`
	Videos      Videos       `json:"videos" db:"-"`
	Audios      Audios       `json:"audios" db:"-"`

	UserID   uuid.UUID `json:"user_id" db:"user_id"`
	Trustees Trustees  `json:"trustees" db:"trustees"`
}

// String is not required by pop and may be deleted
func (w Will) String() string {
	jw, _ := json.Marshal(w)
	return string(jw)
}

// Wills is not required by pop and may be deleted
type Wills []Will

// String is not required by pop and may be deleted
func (w Wills) String() string {
	jw, _ := json.Marshal(w)
	return string(jw)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (w *Will) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (w *Will) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (w *Will) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
