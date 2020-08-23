package graph

import (
	"github.com/gobuffalo/nulls"
	"github.com/kolioDev/after_life/models"
)

var User *models.User

func SetUser(u *models.User) {
	User = u
}

func GetNullable(s *string) nulls.String {
	if s == nil {
		return nulls.String{}
	}
	return nulls.NewString(*s)
}
