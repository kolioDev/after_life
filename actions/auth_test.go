package actions

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gobuffalo/uuid"
	"github.com/kolioDev/after_life/helpers"
	"github.com/kolioDev/after_life/models"
	"golang.org/x/crypto/bcrypt"
)

func (as *ActionSuite) Test_Auth_AuthCallback() {
	//as.Fail("Not Implemented!")
}

func (as *ActionSuite) Test_Auth_AuthGetInitialTokens() {
	as.LoadFixture("session")

	s := &models.Session{}
	as.NoError(as.DB.First(s))
	as.NotEqual("", s.UniqueToken.String)

	uniqueToken := s.UniqueToken.String

	res := as.JSON(fmt.Sprintf("/auth/access/%s", uniqueToken)).Get()
	as.Equal(200, res.Code)

	var resBody = &struct {
		Token        string    `json:"token"`
		ExpiresAt    int       `json:"expires_at"`
		RefreshToken string    `json:"refresh_token"`
		UserID       uuid.UUID `json:"user_id"`
	}{}

	as.NoError(json.Unmarshal(res.Body.Bytes(), resBody))
	claims, err := helpers.DecodeJWT(resBody.Token)
	as.NoError(err)
	as.Equal(s.UserID, claims.UserID)
	as.False(claims.ProfileConfirmed)
	as.NoError(as.DB.First(s))
	as.NoError(bcrypt.CompareHashAndPassword(s.RefreshTokenHash, []byte(resBody.RefreshToken)))
	as.Equal(s.UserID, resBody.UserID)

	res = as.JSON(fmt.Sprintf("/auth/access/%s", uniqueToken)).Get()
	as.Equal(404, res.Code)
}

func (as *ActionSuite) Test_Auth_AuthRefreshToken() {
	var resBody = &struct {
		Token        string `json:"token"`
		ExpiresAt    int    `json:"expires_at"`
		RefreshToken string `json:"refresh_token"`
	}{}

	/**
	Unprotected session
	*/
	as.LoadFixture("session")

	s := &models.Session{}
	uc := &models.UserConfirmation{}

	as.NoError(as.DB.First(s))
	as.NoError(uc.Create(as.DB, &models.User{ID: s.UserID}))

	res := as.JSON(fmt.Sprintf("/auth/access/%s", s.UniqueToken.String)).Get()
	as.Equal(200, res.Code)
	as.NoError(json.Unmarshal(res.Body.Bytes(), resBody))
	as.NotEqual("", resBody.RefreshToken)

	prevRefreshTkn := ""
	for i := 0; i < 6; i++ {
		prevRefreshTkn = resBody.RefreshToken
		res := as.JSON("/auth/token").Post(map[string]string{
			"refresh_token": resBody.RefreshToken,
			"user_id":       s.UserID.String(),
		})
		as.Equalf(200, res.Code, "Expected 200 but got %d (%s) iteration %d", res.Code, res.Body.String(), i)
		as.NoError(json.Unmarshal(res.Body.Bytes(), resBody))
		claims, err := helpers.DecodeJWT(resBody.Token)
		as.NoError(err)
		as.Equal(s.UserID, claims.UserID)
		as.False(claims.ProfileConfirmed)
	}

	res = as.JSON("/auth/token").Post(map[string]string{
		"refresh_token": resBody.RefreshToken + "123",
		"user_id":       s.UserID.String(),
	})
	as.Equal(406, res.Code)

	res = as.JSON("/auth/token").Post(map[string]string{
		"refresh_token": prevRefreshTkn,
		"user_id":       s.UserID.String(),
	})
	as.Equal(406, res.Code)

	res = as.JSON("/auth/token").Post(map[string]string{
		"refresh_token": resBody.RefreshToken,
		"user_id":       strings.ReplaceAll(s.UserID.String(), "1", "2"),
	})
	as.Equal(404, res.Code)

	/**
	Protected session
	*/
	as.NoError(uc.SetSeen(as.DB))
	as.NoError(uc.Confirm(as.DB))
	res = as.JSON("/auth/token").Post(map[string]string{
		"refresh_token": resBody.RefreshToken,
		"user_id":       s.UserID.String(),
	})
	as.Equal(200, res.Code)
	as.NoError(json.Unmarshal(res.Body.Bytes(), resBody))
	claims, err := helpers.DecodeJWT(resBody.Token)
	as.NoError(err)
	as.Equal(s.UserID, claims.UserID)
	as.True(claims.ProfileConfirmed)
}

func (as *ActionSuite) Test_Auth_AuthResetSignUp_Session() {
	as.LoadFixture("session")

	s := &models.Session{}
	uc := &models.UserConfirmation{}

	as.NoError(as.DB.First(s))
	as.NotEqual("", s.UniqueToken.String)
	as.NoError(uc.Create(models.DB, &models.User{ID: s.UserID}))

	//Only having session identifier
	res := as.JSON("/auth/reset").Post(map[string]string{
		"session_identifier": s.UniqueToken.String + "12",
	})
	as.Equalf(404, res.Code, "Expected 404 but got %d with body %s", res.Code, res.Body.String())

	as.DBDelta(-1, "users", func() {
		as.DBDelta(-1, "sessions", func() {
			as.DBDelta(-1, "user_confirmations", func() {
				res = as.JSON("/auth/reset").Post(map[string]string{
					"session_identifier": s.UniqueToken.String,
				})
				as.Equal(200, res.Code)
			})
		})
	})

}

func (as *ActionSuite) Test_Auth_AuthResetSignUp_UserID() {
	as.LoadFixture("session")

	s := &models.Session{}
	uc := &models.UserConfirmation{}

	as.NoError(as.DB.First(s))
	as.NotEqual("", s.UniqueToken.String)
	as.NoError(uc.Create(models.DB, &models.User{ID: s.UserID}))

	//Trying with confirmed user
	as.NoError(uc.Confirm(models.DB))
	res := as.JSON("/auth/reset").Post(map[string]string{
		"user_id": s.UserID.String(),
	})
	as.Equal(403, res.Code)
	uc.Confirmed = false
	as.NoError(as.DB.Save(uc))

	as.DBDelta(-1, "users", func() {
		as.DBDelta(-1, "sessions", func() {
			as.DBDelta(-1, "user_confirmations", func() {
				res := as.JSON("/auth/reset").Post(map[string]string{
					"user_id": s.UserID.String(),
				})
				as.Equal(200, res.Code)
			})
		})
	})

}
