package helpers

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/uuid"
	"github.com/pkg/errors"
)

const JWT_VALID_TIME = 60 * 2 //seconds

// Create the JWT key used to create the signature
var jwtKey = []byte(envy.Get("JWT_TOKEN", "super_secret_jwt_token"))

//Claims Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	UserID           uuid.UUID
	ProfileConfirmed bool
	jwt.StandardClaims
}

//auth - creating auth JWT
func EncodeJWT(userID uuid.UUID, profileConfirmed bool) (string, int64, error) {

	//Определям, кога ще изтече token-a
	// 2 дни валидност
	expirationTime := time.Now().Add(JWT_VALID_TIME * time.Second)

	//Прави JWT claim, който включва дата на изтичането и UserID
	claims := &Claims{
		UserID:           userID,
		ProfileConfirmed: profileConfirmed,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		return "", 0, errors.New("Cannot generate JWT")
	}
	return tokenString, expirationTime.Unix(), nil
}

func DecodeJWT(tknStr string) (*Claims, error) {
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, errors.WithStack(err)
	}

	if !tkn.Valid {
		return nil, errors.New("token is invalid")
	}

	return claims, nil
}
