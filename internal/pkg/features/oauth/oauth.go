package oauth

import (
	"context"
	"fmt"

	"github.startlite.cn/itapp/startlite/internal/pkg/types"
	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
	"github.startlite.cn/itapp/startlite/pkg/lines/featurex"
	"golang.org/x/oauth2"
)

var instance *azureADConfig

type OauthService interface {
	GetAuthCodeUrl(state string) string
	Exchange(ctx context.Context, code string) (*oauth2.Token, error)
}

type azureADConfig struct {
	config *oauth2.Config
}

func MustNewAzureADClient(appCtx appx.AppContext, cl *featurex.ConfigLoader) OauthService {

	azureAdConfig := &types.AzureADConfig{}
	cl.Load(azureAdConfig)
	tokenUrl := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", azureAdConfig.TenantID)
	authUrl := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/authorize", azureAdConfig.TenantID)
	instance = &azureADConfig{
		config: &oauth2.Config{
			ClientID:     azureAdConfig.ClientID,
			ClientSecret: azureAdConfig.ClientSecret,
			RedirectURL:  azureAdConfig.RedirectURL,
			Endpoint: oauth2.Endpoint{
				AuthURL:   authUrl,
				TokenURL:  tokenUrl,
				AuthStyle: oauth2.AuthStyleInParams,
			},
			Scopes: []string{"openid", "profile"},
		},
	}
	return instance
}

func (a *azureADConfig) GetAuthCodeUrl(state string) string {
	url := a.config.AuthCodeURL(state, oauth2.AccessTypeOnline)
	return url
}

func (a *azureADConfig) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {

	token, err := a.config.Exchange(ctx, code, oauth2.AccessTypeOffline)
	return token, err
}
