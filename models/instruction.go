package models

import (
	"encoding/json"
	"github.com/gobuffalo/nulls"
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
	Text    string `json:"text" db:"content"`
	Picture *File  `json:"picture" db:"-"`
	Audio   *File  `json:"audio" db:"-"`
	Video   *File  `json:"video" db:"-"`

	WillID uuid.UUID `json:"will_id" db:"will_id"`
}

//Used for validation purposes
type instructionIndex struct {
	Index nulls.UInt32 `db:"index"`
}

func (ii instructionIndex) TableName() string {
	return "instructions"
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

		&validators.StringLengthInRange{Name: "text", Field: i.Text, Min: 0, Max: 500},

		&validators.UUIDIsPresent{Name: "will_id", Field: i.WillID},
		&validators.FuncValidator{Name: "will_id", Field: i.WillID.String(), Message: "There is now will with id %s", Fn: func() bool {
			if err := tx.Find(&Will{}, i.WillID); err != nil {
				return false
			}
			return true
		}},

		&validators.FuncValidator{Name: "index", Field: strconv.Itoa(int(i.Index)), Message: "Invalid index %s", Fn: func() bool {

			ii := &instructionIndex{}

			if err := tx.Where("will_id=?", i.WillID).Select("MAX(index) as index").First(ii); err != nil {
				return false
			}

			if !ii.Index.Valid {
				return i.Index == 0
			}

			return i.Index == uint(ii.Index.UInt32)+1
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

//Creates instruction entry in the DB
func (i *Instruction) Create(tx *pop.Connection, w Will) (*validate.Errors, error) {
	i.WillID = w.ID
	return tx.ValidateAndCreate(i)
}

func (is *Instructions) Create(tx *pop.Connection, w Will) (*validate.Errors, error) {
	newIs := Instructions{}
	for _, i := range *is {
		i.WillID = w.ID
		newIs = append(newIs, i)
	}
	return tx.ValidateAndCreate(&newIs)
}
