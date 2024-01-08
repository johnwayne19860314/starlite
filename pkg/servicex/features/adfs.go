package features

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/s12v/go-jwks"
	"gopkg.in/square/go-jose.v2"

	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
	"github.startlite.cn/itapp/startlite/pkg/lines/errorx"
	"github.startlite.cn/itapp/startlite/pkg/lines/featurex"
	"github.startlite.cn/itapp/startlite/pkg/servicex/types"
	"github.startlite.cn/itapp/startlite/pkg/servicex/utils"
)

type ADFS struct {
	types.ADFSConfig
}

func MustNewADFSClient(appCtx appx.AppContext, cl *featurex.ConfigLoader) *ADFS {
	adfs := &ADFS{}
	cl.Load(adfs)

	return adfs
}

func (client *ADFS) GetJwkByKid(kid string) (*jose.JSONWebKey, error) {
	jwksSource := jwks.NewWebSource(client.Endpoint, http.DefaultClient)
	jwksClient := jwks.NewDefaultClient(
		jwksSource,
		time.Hour,    // Refresh keys every 1 hour
		12*time.Hour, // Expire keys after 12 hours
	)

	var jwk *jose.JSONWebKey

	jwk, err := jwksClient.GetEncryptionKey(context.Background(), kid)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return jwk, nil
}

func (client *ADFS) GetJwkByToken(token string) (*jose.JSONWebKey, error) {
	kid, err := utils.JwtGetHeaderKeyId(token)
	if err != nil {
		return nil, err
	}

	return client.GetJwkByKid(kid)
}

func (client *ADFS) ValidateToken(token string) error {
	jwk, err := client.GetJwkByToken(token)
	if err != nil {
		return err
	}

	parts := strings.Split(token, ".")
	err = jwt.SigningMethodRS256.Verify(strings.Join(parts[0:2], "."), parts[2], jwk.Key)
	if err != nil {
		return errorx.New("failed to verify jwt token")
	}

	return nil
}

func (client *ADFS) GetClaimFromToken(token string) (jwt.MapClaims, error) {
	jwk, err := client.GetJwkByToken(token)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(token, ".")
	err = jwt.SigningMethodRS256.Verify(strings.Join(parts[0:2], "."), parts[2], jwk.Key)
	if err != nil {
		return nil, errorx.New("failed to verify jwt token")
	}

	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return jwk.Key, nil
	}, jwt.WithoutAudienceValidation())
	if err != nil {
		switch err.(type) {
		case *jwt.TokenExpiredError:
			return nil, errorx.New("token expired")
		case *jwt.TokenNotValidYetError:
			return nil, errorx.New("token not valid yet")
		case *jwt.InvalidIssuerError:
			return nil, errorx.New("token issuer is invalid")
		default:
			return nil, errorx.WithStack(err)
		}
	}

	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
