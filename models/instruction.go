package models

import (
	"encoding/json"
	"github.com/gobuffalo/validate/validators"
	"strconv"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
)

type Instruction struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	Index   uint   `json:"index" db:"index"`
	Text    string `json:"text"`
	Picture `json:"picture" db:"-"`
	Audio   `json:"audio" db:"audio"`

	WillID uuid.UUID `json:"will_id" db:"will_id"`
}

// String is not required by pop and may be deleted
func (i Instruction) String() string {
	ji, _ := json.Marshal(i)
	return string(ji)
}

// Instructions is not required by pop and may be deleted
type Instructions []Instruction

// String is not required by pop and may be deleted
func (i Instructions) String() string {
	ji, _ := json.Marshal(i)
	return string(ji)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (i *Instruction) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(

		&validators.UUIDIsPresent{Name: "will_id", Field: i.WillID},
		&validators.FuncValidator{Name: "will_id", Field: i.WillID.String(), Message: "There is now will with id %s", Fn: func() bool {
			if err := tx.Find(&Will{}, i.WillID); err != nil {
				return false
			}
			return true
		}},

		&validators.FuncValidator{Name: "index", Field: strconv.Itoa(int(i.Index)), Fn: func() bool {
			w := &Instructions{}
			if err := tx.Where("will_id", i.WillID).All(w); err != nil {
				return false
			}

			return true
		}},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (i *Instruction) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (i *Instruction) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
