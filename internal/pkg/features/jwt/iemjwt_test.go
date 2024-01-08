package jwt

import (
	"testing"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/stretchr/testify/assert"
)

func TestIsValidJWT(t *testing.T) {
	secret := "secret"
	tchjwt := IEMJWT{}

	claims := jwt.StandardClaims{}

	validToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	invalidToken := "invalidToken"

	// valid return no err
	assert.Nil(t, tchjwt.IsValidJWT(validToken, secret))
	// notValid return err
	err := tchjwt.IsValidJWT(invalidToken, secret)
	assert.NotNil(t, err, "Expected error but got nil")
}
