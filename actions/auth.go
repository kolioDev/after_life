// actions/auth.go
package actions

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/kolioDev/after_life/helpers"
	"github.com/kolioDev/after_life/models"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/instagram"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/markbates/goth/gothic"
)

/**
Auth sign up flow:
FE -> request to BE for oauth
BE - generate and save User, UserConfirmation, Session (with unique identifier)
BE -> redirect to FE width /{unique_id}
FE -> request JWT and refresh_token
BE - reset the unique_id of the session
BE -> send JWT and refresh_token
FE -> request user_confirmation keys [protected]
BE -> send keys and check "sent" column
FE -> click a button to confirm keys read
BE -> create new JWT with flag confirmed

*/
func init() {
	gothic.Store = App().SessionStore
	goth.UseProviders(
		facebook.New(os.Getenv("FACEBOOK_KEY"), os.Getenv("FACEBOOK_SECRET"), fmt.Sprintf("%s%s", App().Host, "/oauth/facebook/callback")),
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), fmt.Sprintf("%s%s", App().Host, "/oauth/google/callback")),
		instagram.New(os.Getenv("INSTAGRAM_KEY"), os.Getenv("INSTAGRAM_SECRET"), fmt.Sprintf("%s%s", App().Host, "/oauth/instagram/callback")),
	)
}

/**
Receives OAuth and sends keys
*/
func AuthOAuthCallback(c buffalo.Context) error {
	tx := models.DB //c.Value("tx").(*pop.Connection)
	gothU, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return c.Error(401, err)
	}
	u := &models.User{}
	q := tx.Where("provider = ? and provider_id = ?", gothU.Provider, gothU.UserID)

	exists, err := q.Exists(u)
	if err != nil {
		return errors.WithStack(err)
	}

	createSessAndRedirect := func(u *models.User) error {
		//create session and redirect to a something/{session_unique_token} path
		s := models.Session{}
		if err := s.Create(tx, u, false); err != nil {
			return errors.WithStack(err)
		}

		return c.Redirect(303,
			envy.Get("FRONTEND_URL", "127.0.0.1:8080")+
				strings.Replace(
					envy.Get("FRONTEND_AUTH_SUCCESS_PATH", "/oath/confirmed/{session_unique_token}"),
					"{session_unique_token}", s.UniqueToken.String, 1),
		)
	}

	//If the user already exists
	if exists {

		q.First(u)

		uc := &models.UserConfirmation{}
		err = tx.Where("user_id=?", u.ID).First(uc)
		if err != nil {
			if errors.Cause(err) != sql.ErrNoRows {
				return errors.WithStack(err)
			}
		}

		//Log in existing user
		if exists && uc.Confirmed {
			return createSessAndRedirect(u)
		} else {
			return c.Redirect(303,
				envy.Get("FRONTEND_URL", "127.0.0.1:8080")+
					envy.Get("FRONTEND_AUTH_ERROR_PATH", "/oath/error/"),
			)
		}

	}

	_, err = u.Create(tx)
	if err != nil {
		return errors.WithStack(err)
	}

	u = &models.User{}
	u.Provider = gothU.Provider
	u.ProviderID = gothU.UserID
	u.Name = gothU.Name
	if u.Name == "" {
		u.Name = gothU.Email
	}
	if u.Name == "" {
		u.Name = gothU.NickName
	}
	if u.Name == "" {
		u.Name = gothU.FirstName + " " + gothU.LastName
	}
	u.Name = strings.TrimSpace(u.Name)

	verrs, err := u.Create(tx)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs!=nil && verrs.HasAny() {
		return c.Render(406, r.JSON(verrs.Errors))
	}

	return createSessAndRedirect(u)
}

/**
* @api {get} /auth/access/:session_identifier Sign up JWT
* @apiName GetUnconfirmedJWT
* @apiDescription  Gets JWT and refresh token after OAuth sign up
@apiGroup Auth
*
* @apiParam {String} session_identifier session unique indentifier
*
* @apiSuccess {String} token JWT
* @apiSuccess {Number} expires_at UNIX expires at
* @apiSuccess {String} refresh_token token used to get new JWT
* @apiSuccess {String} user_id the uuid of the user connected to the session
* @apiSuccess {Boolean} profile_confirmed if true the user has confirmed keys seen and received
*
* @apiError (404) SessionNotFund  Session with identifier <code>session_identifier</code> was not found.
**
*/
func AuthGetInitialTokens(c buffalo.Context) error {
	tx := models.DB //c.Value("tx").(*pop.Connection)
	s := &models.Session{}
	u := &models.User{}
	uc := &models.UserConfirmation{}

	if err := tx.Where("unique_token=?", c.Param("session_identifier")).Where("unique_token<>?", "").First(s); err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return c.Render(404, r.JSON("Session do not exists"))
		}
		return errors.WithStack(err)
	}

	if err := tx.Find(u, s.UserID); err != nil {
		return errors.WithStack(err)
	}

	if err := tx.Where("user_id=?", u.ID).First(uc); err != nil {
		log.Println("---169---")
		return errors.WithStack(err)
	}

	JWT, expiresAt, err := helpers.EncodeJWT(u.ID, uc.Confirmed)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := s.ResetUniqueToken(tx); err != nil {
		return errors.WithStack(err)
	}

	return c.Render(200, r.JSON(map[string]interface{}{
		"token":             JWT,
		"expires_at":        expiresAt,
		"refresh_token":     s.RefreshToken,
		"user_id":           s.UserID,
		"profile_confirmed": uc.Confirmed,
	}))
}

/**
* @api {post} /auth/token Refresh the JWT
* @apiName RefreshJWT
* @apiDescription  Gets new refresh token and JWT
* @apiGroup Auth
*
* @apiParam (Body) {String} refresh_token the refresh token of the previous session
* @apiParam (Body) {String} user_id the id of the logged in user
*
* @apiSuccess {String} token JWT
* @apiSuccess {Number} expires_at UNIX expires at
* @apiSuccess {String} refresh_token token used to get new JWT
* @apiSuccess {Boolean} profile_confirmed if true the user has confirmed keys seen and received`
*
* @apiError (404) UserNotFund User with id <code>user_id</code> was not found.
* @apiError (404) SessionNotFund  Session for user with <code>user_id</code> was not found.
* @apiError (406) InvalidRefreshToken  <code>refresh_token</code> did not match the last session.
**
 */
func AuthRefreshToken(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	s := &models.Session{}
	uc := &models.UserConfirmation{}
	u := &models.User{}
	//Post params
	var refreshToken string
	var uID string

	parsedC, _ := postParse(c)
	refreshToken = parsedC.Request().PostForm.Get("refresh_token")
	uID = parsedC.Request().PostForm.Get("user_id")

	if err := tx.Find(u, uID); err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return c.Render(404, r.JSON("User not found"))
		}
		return errors.WithStack(err)
	}

	if err := tx.Where("user_id=?", uID).Last(s); err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return c.Render(404, r.JSON("Session do not exists"))
		}
		return errors.WithStack(err)
	}

	if bcrypt.CompareHashAndPassword(s.RefreshTokenHash, []byte(refreshToken)) != nil {
		return c.Render(406, r.JSON("invalid refresh token"))
	}

	//determine if the jwt should be with flag protected or not and generate JWT
	if err := tx.Where("user_id=?", uID).First(uc); err != nil {
		return errors.WithStack(err)
	}

	//generate new session
	if err := tx.Destroy(s); err != nil {
		return errors.WithStack(err)
	}

	if err := s.Create(tx, u, uc.Confirmed); err != nil {
		return errors.WithStack(err)
	}

	//generate new JWT
	JWT, expiresAt, err := helpers.EncodeJWT(uuid.FromStringOrNil(uID), uc.Confirmed)
	if err != nil {
		return errors.WithStack(err)
	}

	return c.Render(200, r.JSON(map[string]interface{}{
		"token":             JWT,
		"expires_at":        expiresAt,
		"refresh_token":     s.RefreshToken,
		"profile_confirmed": uc.Confirmed,
	}))
}

/**
* @api {post} /auth/reset Resets signUP
* @apiName ResetAuth
* @apiDescription Resets the registration process in case of error
* @apiGroup Auth
*
* @apiParam (Body) {String} [session_identifier] the refresh token of the previous session
* @apiParam (Body) {String} [user_id] the id of the user
*
*
* @apiError (403) AccountConfirmed The account associated with the user was confirmed
* @apiError (404) SessionNotFund  Session for user with <code>user_id</code> was not found.
* @apiError (406) CannotDeleteUser  Unknown error
**
 */
func AuthResetSignUp(c buffalo.Context) error {
	tx := models.DB //c.Value("tx").(*pop.Connection)

	parsedC, _ := postParse(c)
	uID := parsedC.Request().PostForm.Get("user_id")
	sIdentifier := parsedC.Request().PostForm.Get("session_identifier")

	u := &models.User{}
	uc := &models.UserConfirmation{}
	s := &models.Session{}

	if uID != "" {
		if err := tx.Where("user_id=?", uID).First(uc); err != nil {
			if errors.Cause(err) == sql.ErrNoRows {
				return c.Render(404, r.JSON("User confirmation do not exists"))
			}
			return errors.WithStack(err)
		}
		if uc.Confirmed {
			return c.Render(403, r.JSON("Cannot delete confirmed account"))
		}
		if err := tx.Find(u, uID); err != nil {
			return errors.WithStack(err)
		}
	} else if sIdentifier != "" {
		if err := tx.Where("unique_token=?", sIdentifier).First(s); err != nil {
			if errors.Cause(err) == sql.ErrNoRows {
				return c.Render(404, r.JSON("Session do not exists"))
			}
			return errors.WithStack(err)
		}
		if err := tx.Find(u, s.UserID); err != nil {
			return errors.WithStack(err)
		}
	} else {
		return c.Render(406, r.JSON("cannot delete user - unknown error"))
	}

	if err := tx.Destroy(u); err != nil {
		return errors.WithStack(err)
	}

	return c.Render(200, r.JSON("user deleted"))
}

//AuthMiddleware - requires JWT
//params - true = do not require the user to be confirmed
func AuthMiddleware(next buffalo.Handler, params ...bool) buffalo.Handler {
	return func(c buffalo.Context) error {

		parsedC, _ := postParse(c)
		tknStr := parsedC.Request().PostForm.Get("_token")

		if tknStr == "" {
			tknStr = c.Param("_token")
		}

		claims, err := helpers.DecodeJWT(tknStr)
		if err != nil {
			return c.Render(401, r.JSON("Unauthorized"))
		}

		if len(params) == 0 && !claims.ProfileConfirmed {
			return c.Render(401, r.JSON("Unauthorized"))
		}

		if len(params) > 0 {
			if params[0] && claims.ProfileConfirmed {
				return c.Render(401, r.JSON("Unauthorized"))
			}
			if !params[0] && !claims.ProfileConfirmed {
				return c.Render(401, r.JSON("Unauthorized"))
			}
		}

		u := &models.User{}
		if err := models.DB.Find(u, claims.UserID); err != nil {
			return c.Render(401, r.JSON("Unauthorized -  cannot find user"))
		}

		//Add the user to the context
		c.Set("current_user", u)

		return next(c)
	}
}
