package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate/validators"
	"github.com/kolioDev/after_life/scalars"
	"github.com/pkg/errors"

	"github.com/gobuffalo/pop"

	"github.com/gobuffalo/validate"

	"github.com/kolioDev/after_life/graphql/model"
)

type Trustee struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	UserID uuid.UUID `json:"user_id" db:"user_id"`

	Name           string       `json:"name" db:"name"`
	Relationship   string       `json:"relationship" db:"relationship"`
	Phone          string       `json:"phone" db:"phone"`
	Email          string       `json:"email" db:"email"`
	FacebookLink   nulls.String `json:"facebook_link" db:"facebook_link"`
	TwitterLink    nulls.String `json:"twitter_link" db:"twitter_link"`
	AdditionalInfo nulls.String `json:"additional_info" db:"additional_information"`
}

var TRUSTEE_RELATIONSHIP_TYPES = []string{
	"father",
	"mother",
	"other_relative",
	"best_friend",
	"school_friend",
	"college_friend",
	"other_friend",
	"husband",
	"wife",
	"girlfriend",
	"boyfriend",
	"fiance",
	"acquaintance",
}

// String is not required by pop and may be deleted
func (t Trustee) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Trustees is not required by pop and may be deleted
type Trustees []Trustee

// String is not required by pop and may be deleted
func (t Trustees) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *Trustee) Validate(tx *pop.Connection) (*validate.Errors, error) {
	vals := []validate.Validator{
		&validators.UUIDIsPresent{Name: "user_id", Field: t.UserID},
		&validators.StringIsPresent{Name: "name", Field: t.Name},
		&validators.StringIsPresent{Name: "relationship", Field: t.Relationship},
		&validators.StringIsPresent{Name: "email", Field: t.Email},
		&validators.StringIsPresent{Name: "phone", Field: t.Phone},

		&validators.StringLengthInRange{Name: "name", Field: t.Name, Min: 2, Max: 250},
		&validators.RegexMatch{Name: "relationship", Field: t.Relationship, Expr: fmt.Sprintf("(%s)", strings.Join(TRUSTEE_RELATIONSHIP_TYPES, "|")), Message: "invalid relationship type"},
		&validators.RegexMatch{Name: "phone", Field: t.Phone, Expr: `^(\+|00)\d{6,13}$`, Message: "invalid phone"},
		&validators.EmailIsPresent{Name: "email", Field: t.Email},

		&validators.StringLengthInRange{Name: "additional_info", Min: 0, Max: 500, Field: t.AdditionalInfo.String},

		&validators.FuncValidator{Name: "user_id", Field: t.UserID.String(), Message: "%s is not a valid user id", Fn: func() bool {
			if err := tx.Find(&User{}, t.UserID); err != nil {
				return false
			}
			return true
		}},
	}

	if t.FacebookLink.Valid {
		vals = append(vals,
			&validators.URLIsPresent{Name: "facebook", Field: t.FacebookLink.String},
			&validators.RegexMatch{Name: "facebook", Field: t.FacebookLink.String, Expr: `https?:\/\/(m.|mobile.|www.)?facebook\.com\/.*`},
		)

	}

	if t.TwitterLink.Valid {
		vals = append(vals,
			&validators.URLIsPresent{Name: "twitter", Field: t.TwitterLink.String},
			&validators.RegexMatch{Name: "twitter", Field: t.TwitterLink.String, Expr: `^https?://(m.|mobile.|www.)?twitter\.com/.*$`},
		)
	}

	return validate.Validate(vals...), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *Trustee) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (t *Trustee) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	oldT := &Trustee{}
	if err := tx.Find(oldT, t.ID); err != nil {
		return validate.NewErrors(), errors.WithStack(err)
	}

	return validate.Validate(
		&validators.StringsMatch{Name: "user_id", Field: t.UserID.String(), Field2: oldT.UserID.String()},
	), nil
}

//Creates trustee entry in the DB
func (t *Trustee) Create(tx *pop.Connection, userID uuid.UUID) (*validate.Errors, error) {
	t.UserID = userID
	return tx.ValidateAndCreate(t)
}

//Updates trustee entry in the DB
func (t *Trustee) Update(tx *pop.Connection) (*validate.Errors, error) {
	return tx.ValidateAndUpdate(t)
}

func (t *Trustees) GetForUser(tx *pop.Connection, userID uuid.UUID, orderBy string, order string) error {

	if orderBy == "" {
		orderBy = "created_at"
	}

	if order == "" {
		order = "asc"
	}

	return tx.Where("user_id=?", userID).Order(fmt.Sprintf("%s %s", orderBy, order)).All(t)
}

func (t Trustee) ToGraphQL() *model.Trustee {
	return &model.Trustee{
		ID:           scalars.ModelsUUID2GhqlUUID(t.ID),
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
		Relationship: model.TrusteeType(t.Relationship),
		Name:         t.Name,
		Email:        t.Email,
		Phone:        t.Phone,
		FacebookLink: NullableToString(t.FacebookLink),
		TwitterLink:  NullableToString(t.TwitterLink),
	}
}

func (ts Trustees) ToGraphQL() []*model.Trustee {
	var QLTrustees []*model.Trustee
	for _, t := range ts {
		QLTrustees = append(QLTrustees, t.ToGraphQL())
	}
	return QLTrustees
}
