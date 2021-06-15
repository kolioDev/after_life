package resolvers

import (
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	"github.com/kolioDev/after_life/models"
)

var User *models.User
var TX *pop.Connection

func SetUser(u *models.User) {
	User = u
}

func SetTX(tx *pop.Connection) {
	TX = tx
}

func GetNullableString(s *string) nulls.String {
	if s == nil {
		return nulls.String{}
	}
	return nulls.NewString(*s)
}

func GetNullableInt(i *int) nulls.Int {
	if i == nil {
		return nulls.Int{}
	}
	return nulls.NewInt(*i)
}
