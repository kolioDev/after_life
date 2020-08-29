package models

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/kolioDev/after_life/helpers"
	"github.com/pkg/errors"

	"github.com/gobuffalo/envy"
	"github.com/tjarratt/babble"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
)

const KEYS_NUMBER = 5

type UserConfirmation struct {
	ID UUID `json:"id" db:"id"`

	UserID UUID `json:"user_id" db:"user_id"`

	Keys                 string `json:"-" db:"-"`
	KeysEncrypted        []byte `json:"-" db:"keys"`
	KeysSent             bool   `json:"keys_sent" db:"keys_sent"`
	KeysSeenConfirmation bool   `json:"keys_seen_confirmation" db:"keys_seen"`
	Confirmed            bool   `json:"confirmed" db:"confirmed"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (u UserConfirmation) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// UserConfirmations is not required by pop and may be deleted
type UserConfirmations []UserConfirmation

// String is not required by pop and may be deleted
func (u UserConfirmations) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *UserConfirmation) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (u *UserConfirmation) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (u *UserConfirmation) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (u *UserConfirmation) Create(tx *pop.Connection, user *User) error {
	k := envy.Get("APP_KEY", "password_123")

	if user.ID == UUIDNil() {
		return errors.New("user cannot be empty")
	}

	u.Keys = generateKeys()
	u.KeysEncrypted = helpers.Encrypt(strings.ToLower(u.Keys), k)
	u.KeysSent = false
	u.Confirmed = false
	u.KeysSeenConfirmation = false
	u.UserID = user.ID
	return tx.Create(u)
}

func (u *UserConfirmation) SendKeys(tx *pop.Connection) error {
	k := envy.Get("APP_KEY", "password_123")
	var err error
	u.Keys, err = helpers.Decrypt(u.KeysEncrypted, k)
	if err != nil {
		return errors.WithStack(err)
	}

	u.KeysSent = true

	if err := tx.Save(u, "keys", "keys_seen", "confirmed", "created_at", "user_id"); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (u *UserConfirmation) CheckKeysMatch(keys string) (bool, error) {
	//TODO::test
	k := envy.Get("APP_KEY", "password_123")
	match := false
	var err error
	u.Keys, err = helpers.Decrypt(u.KeysEncrypted, k)

	if err != nil {
		return false, errors.WithStack(err)
	}

	match = strings.TrimSpace(u.Keys) == strings.TrimSpace(keys)

	u.Keys = ""
	return match, nil
}

func (u *UserConfirmation) SetSeen(tx *pop.Connection) error {
	//TODO::test
	u.KeysSeenConfirmation = true
	return tx.Save(u, "keys", "keys_sent", "confirmed", "created_at", "user_id")
}

func (u *UserConfirmation) Confirm(tx *pop.Connection) error {
	//TODO::test
	u.Confirmed = true
	return tx.Save(u, "keys", "keys_sent", "keys_seen", "created_at", "user_id")
}

func generateKeys() string {
	babbler := babble.NewBabbler()
	babbler.Separator = " "
	babbler.Count = KEYS_NUMBER
	words := babbler.Babble()
	return strings.ReplaceAll(words, "'s", "")
}
