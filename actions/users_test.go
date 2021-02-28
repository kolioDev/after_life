package actions

import (
	"encoding/json"
	"fmt"

	"github.com/kolioDev/after_life/helpers"
	"github.com/kolioDev/after_life/models"
	"golang.org/x/crypto/bcrypt"
)

func (as *ActionSuite) Test_Users_GetConfirmationKeys() {
	as.LoadFixture("login user")
	u := &models.User{}
	uc := &models.UserConfirmation{}

	as.NoError(as.DB.Where("username=?", "login_usere@test.com").First(u))
	as.NoError(uc.Create(as.DB, u))

	JWT, _, err := helpers.EncodeJWT(u.ID, false)
	as.NoError(err)

	res := as.JSON(fmt.Sprintf("/user/keys?_token=%s", JWT)).Get()
	as.Equal(200, res.Code)

	var resBody = &struct {
		Keys string `json:"keys"`
	}{}

	as.NoError(json.Unmarshal(res.Body.Bytes(), resBody))
	as.NotEqual("", resBody.Keys)

	res = as.JSON(fmt.Sprintf("/user/keys?_token=%s", JWT)).Get()
	as.Equal(403, res.Code)
}

func (as *ActionSuite) Test_Users_UsersConfirmKeysSeen() {
	as.LoadFixture("login user")
	u := &models.User{}
	uc := &models.UserConfirmation{}
	s := &models.Session{}

	as.NoError(as.DB.Where("username=?", "login_usere@test.com").First(u))
	as.NoError(uc.Create(as.DB, u))
	as.NoError(s.Create(as.DB, u, false))

	JWT, _, err := helpers.EncodeJWT(u.ID, false)
	as.NoError(err)

	//get keys
	res := as.JSON(fmt.Sprintf("/user/keys?_token=%s", JWT)).Get()
	as.Equal(200, res.Code)

	var resBody = &struct {
		Keys string `json:"keys"`
	}{}
	as.NoError(json.Unmarshal(res.Body.Bytes(), resBody))
	as.NotEqual("", resBody.Keys)

	//confirm keys received
	res = as.JSON("/user/confirm/keys").Post(map[string]string{
		"_token": JWT,
		"keys":   resBody.Keys + "111",
	})
	as.Equal(406, res.Code)

	res = as.JSON("/user/confirm/keys").Post(map[string]string{
		"_token": JWT,
		"keys":   resBody.Keys,
	})
	as.Equal(200, res.Code)

	//Check the jwt
	var resBodyJWT = &struct {
		Token        string `json:"token"`
		ExpiresAt    int    `json:"expires_at"`
		RefreshToken string `json:"refresh_token"`
	}{}

	as.NoError(json.Unmarshal(res.Body.Bytes(), resBodyJWT))
	claims, err := helpers.DecodeJWT(resBodyJWT.Token)
	as.NoError(err)
	as.Equal(s.UserID, claims.UserID)
	as.True(claims.ProfileConfirmed)
	as.NoError(as.DB.First(s))
	as.NoError(bcrypt.CompareHashAndPassword(s.RefreshTokenHash, []byte(resBodyJWT.RefreshToken)))

	as.NoError(as.DB.Find(uc, uc.ID))
	as.True(uc.KeysSent)
	as.True(uc.KeysSeenConfirmation)
	as.True(uc.Confirmed)
}
