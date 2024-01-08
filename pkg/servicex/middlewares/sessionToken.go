package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"golang.org/x/oauth2"

	"net/url"

	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
	"github.startlite.cn/itapp/startlite/pkg/lines/errorx"
	"github.startlite.cn/itapp/startlite/pkg/lines/featurex"
	"github.startlite.cn/itapp/startlite/pkg/lines/rds"
	"github.startlite.cn/itapp/startlite/pkg/lines/utilx/claimx"
	"github.startlite.cn/itapp/startlite/pkg/lines/utilx/textx"
	"github.startlite.cn/itapp/startlite/pkg/lines/utilx/tokenx"
	def "github.startlite.cn/itapp/startlite/pkg/servicex/definition"
	"github.startlite.cn/itapp/startlite/pkg/servicex/utils"
)

type SessionToken struct {
	AppCtx      appx.AppContext  `inject:""`
	ReqCtx      appx.ReqContext  `inject:""`
	RedisClient *rds.RedisClient `inject:""`
	OAuthConfig *featurex.OAuth2 `inject:""`

	mu         sync.Mutex
	Token      string
	OAuthToken *oauth2.Token
	UserClaims *jwt.MapClaims
}

type AccessTokener interface {
	GetAccessToken() (string, error)
}

var SessionTokenName = tokenx.WorkflowToken

func CheckSessionToken(reqCtx appx.ReqContext) {
	st := SessionToken{}
	err := reqCtx.Apply(&st)
	if err != nil {
		panic(err)
	}

	rolesRequired := GetRolesRequired(reqCtx)

	st.Token = GetTokenFromRequest(reqCtx.Gin(), SessionTokenName)
	if v, err := url.QueryUnescape(st.Token); err == nil {
		st.Token = v
	}
	if textx.Blank(st.Token) {
		if rolesRequired == nil {
			return
		}
	}

	st.UserClaims, err = st.GetUserClaims()
	if err != nil {
		if rolesRequired == nil {
			return
		}

		st.ReqCtx.Gin().AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{
			"message": err.Error(),
			"xid":     st.ReqCtx.GetXid(),
		})
		return
	}

	reqCtx.Provide(&st)
	reqCtx.Provide(st.UserClaims)
	reqCtx.Provide(st.TokenerFn())
	reqCtx.Provide(st.AccessTokenFn())
	reqCtx.With("user", claimx.GetClaim(st.UserClaims, "email"))

	if len(rolesRequired) == 0 { // only require valid token or needn't token
		return
	}

	rolesHas := claimx.GetRoles(st.UserClaims)
	if !utils.InArrayAnyString(rolesRequired, rolesHas) {
		reqCtx.Gin().AbortWithStatusJSON(http.StatusForbidden, map[string]string{
			"message": fmt.Sprintf("API need role %v", rolesRequired),
			"xid":     reqCtx.GetXid(),
		})
	}
}

func (st *SessionToken) GetOAuthToken() (*oauth2.Token, error) {
	if st.OAuthToken != nil {
		return st.OAuthToken, nil
	}

	var oauthToken *oauth2.Token
	var err error

	// jwt token
	if len(st.Token) > 256 && len(strings.Split(st.Token, ".")) == 3 {
		oauthToken = &oauth2.Token{
			AccessToken: st.Token,
		}

		return oauthToken, nil
	}

	// session token
	oauthToken, err = st.GetOAuthTokenInCache()
	if err != nil {
		return nil, err
	}

	if oauthToken.Valid() {
		return oauthToken, nil
	}

	oauthToken, err = st.DoRefreshToken(false)
	if err != nil {
		return nil, err
	}
	if oauthToken.Valid() {
		return oauthToken, nil
	}

	return nil, errorx.New("refresh token fail")
}

func (st *SessionToken) GetOAuthTokenInCache() (*oauth2.Token, error) {
	if textx.Blank(st.Token) {
		return nil, def.ErrAPINeedSessionToken
	}

	var oauthTokenString string
	err := st.RedisClient.Get(fmt.Sprintf(def.RedisKeySessionToken, SessionTokenName, st.Token), &oauthTokenString)
	if err != nil {
		return nil, def.ErrInvalidSessionToken
	}

	var oauthToken oauth2.Token
	err = json.Unmarshal([]byte(oauthTokenString), &oauthToken)
	if err != nil {
		return nil, def.ErrInvalidAuthToken
	}

	return &oauthToken, nil
}

func (st *SessionToken) DoRefreshToken(mustRefresh bool) (*oauth2.Token, error) {
	st.mu.Lock()
	defer st.mu.Unlock()

	redisMutex := rds.NewRedisMutex(fmt.Sprintf(def.RedisKeyMutexSessionToken, SessionTokenName, st.Token), st.RedisClient.UniversalClient,
		rds.WithExpiry(time.Second*10),
		rds.WithDelay(time.Millisecond*200),
		rds.WithTries(50))

	err := redisMutex.Lock()
	if err != nil {
		return nil, err
	}
	defer redisMutex.Unlock()

	oauthToken, err := st.GetOAuthTokenInCache()
	if err != nil {
		return nil, err
	}
	if oauthToken.Valid() && !mustRefresh {
		return oauthToken, nil
	}

	oauth2Config, err := st.OAuthConfig.GetConfig("")
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	newToken, err := oauth2Config.TokenSource(st.ReqCtx, oauthToken).Token()
	if err != nil {
		return nil, err
	}
	tokenJsonBytes, err := json.Marshal(newToken)
	if err != nil {
		return nil, err
	}
	err = st.RedisClient.Set(
		fmt.Sprintf(def.RedisKeySessionToken, SessionTokenName, st.Token), string(tokenJsonBytes),
		featurex.GetTokenRefreshExpiresIn(newToken)-20)
	if err != nil {
		return nil, err
	}

	return newToken, nil
}

func (st *SessionToken) GetToken() string {
	return st.Token
}

func (st *SessionToken) TokenerFn() func() tokenx.Tokener {
	return func() tokenx.Tokener {
		return st
	}
}

func (st *SessionToken) GetAccessToken() (accessToken string, err error) {
	if st.OAuthToken == nil {
		st.OAuthToken, err = st.GetOAuthToken()
		if err != nil {
			return "", err
		}
	}

	return st.OAuthToken.AccessToken, nil
}

func (st *SessionToken) AccessTokenFn() func() AccessTokener {
	return func() AccessTokener {
		return st
	}
}

func (st *SessionToken) GetUserClaims() (*jwt.MapClaims, error) {
	var err error
	if st.OAuthToken == nil {
		st.OAuthToken, err = st.GetOAuthToken()
		if err != nil {
			return nil, err
		}
	}
	claims, err := st.OAuthConfig.GetClaims(st.OAuthToken.AccessToken)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return &claims, nil
}
