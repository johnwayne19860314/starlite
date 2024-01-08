package jwt

import (
	"github.com/dgrijalva/jwt-go/v4"

	"github.startlite.cn/itapp/startlite/internal/pkg/types"
	"github.startlite.cn/itapp/startlite/pkg/lines/errorx"
	"github.startlite.cn/itapp/startlite/pkg/lines/featurex"
)

type IEMJWT struct {
	types.JWTConfig
}

func NewIEMJWT(cl *featurex.ConfigLoader) *IEMJWT {
	jwt := &IEMJWT{}
	cl.Load(jwt)
	return jwt
}

// IsValidJWT check if JWT token is valid
func (j *IEMJWT) IsValidJWT(tokenString string, secret string) error {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}, jwt.WithLeeway(10))

	if err != nil {
		return errorx.Errorf("parse error: %w", err)
	}

	if _, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		return nil
	} else {
		return errorx.Errorf("invalid token")
	}
}
