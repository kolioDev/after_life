package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/kolioDev/after_life/helpers"
	"github.com/kolioDev/after_life/models"
	"github.com/pkg/errors"
)

/**
* @api {get} /user/keys Get unique user keys
* @apiName getKeys
* @apiPermission protected unverified
* @apiDescription  Gets unique keys that the user needs to save securely himself
* @apiGroup User
*
* @apiSuccess {String} keys List of keys separated by " "
*
* @apiError (403) KeysAlreadySent  Keys have been marked as <code>sent</code>
*
 */
func UsersGetConfirmationKeys(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	u := c.Value("current_user").(*models.User)
	uc := &models.UserConfirmation{}

	if err := tx.Where("user_id=?", u.ID).First(uc); err != nil {
		return errors.WithStack(err)
	}

	if uc.KeysSent {
		return c.Render(403, r.JSON("Keys already sent"))
	}

	if err := uc.SendKeys(tx); err != nil {
		return errors.WithStack(err)
	}

	return c.Render(200, r.JSON(map[string]string{
		"keys": uc.Keys,
	}))
}

/**
* @api {post} /user/confirm/keys Confirm keys seen
* @apiName confirmKeys
* @apiPermission protected unverified
* @apiDescription Checks the seen column for the keys in the DB and returns new, verified JWT and refresh token
* @apiGroup User
*
* @apiParam (Body) {String} keys the keys received
*
* @apiSuccess {String} token JWT
* @apiSuccess {Number} expires_at UNIX expires at
* @apiSuccess {String} refresh_token token used to get new JWT
*
* @apiError (406) KeysDoNotMatch   <code>keys</code> are not the ones sent
*
 */
func UsersConfirmKeysSeen(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	s := &models.Session{}
	uc := &models.UserConfirmation{}
	u := c.Value("current_user").(*models.User)
	var keys string

	parsedC, _ := postParse(c)
	keys = parsedC.Request().PostForm.Get("keys")

	//get UserConfirmation
	if err := tx.Where("user_id=?", u.ID).First(uc); err != nil {
		return errors.WithStack(err)
	}

	keysMatch, err := uc.CheckKeysMatch(keys)
	if err != nil {
		return errors.WithStack(err)
	}

	if !keysMatch {
		return c.Render(406, r.JSON("keys do not match"))
	}

	//get session
	if err := tx.Where("user_id=?", u.ID).First(s); err != nil {
		return errors.WithStack(err)
	}

	if err := tx.Destroy(s); err != nil {
		return errors.WithStack(err)
	}

	if err := s.Create(tx, u, true); err != nil {
		return errors.WithStack(err)
	}

	if err := uc.SetSeen(tx); err != nil {
		return errors.WithStack(err)
	}

	if err := uc.Confirm(tx); err != nil {
		return errors.WithStack(err)
	}

	JWT, expiresAt, err := helpers.EncodeJWT(u.ID, true)
	if err != nil {
		return errors.WithStack(err)
	}

	return c.Render(200, r.JSON(map[string]interface{}{
		"token":         JWT,
		"expires_at":    expiresAt,
		"refresh_token": s.RefreshToken,
	}))

}
