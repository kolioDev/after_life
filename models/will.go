package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/nulls"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
)

type Will struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	Title    string       `json:"title" db:"title"`
	Priority nulls.UInt32 `json:"importance" db:"priority"`

	Instructions *Instructions `json:"instructions" db:"-"`
	Pictures     *File         `json:"pictures" db:"-"`
	Videos       *File         `json:"videos" db:"-"`
	Audios       *File         `json:"audios" db:"-"`

	UserID   uuid.UUID `json:"user_id" db:"user_id"`
	Trustees *Trustees `json:"trustees" db:"-"`
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
	//TODO::implement
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (w *Will) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

//Creates instruction entry in the DB
func (w *Will) Create(tx *pop.Connection) (*validate.Errors, error) {
	//TODO::implement
	//Validate all instructions
	//Save instructions

	//validate&create will
	return validate.NewErrors(), nil
}
