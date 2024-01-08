package asset

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	url2 "net/url"

	"github.startlite.cn/itapp/startlite/internal/pkg/features/clients"
	baseclient "github.startlite.cn/itapp/startlite/internal/pkg/features/clients/base"
	"github.startlite.cn/itapp/startlite/internal/pkg/types"
	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
	"github.startlite.cn/itapp/startlite/pkg/lines/errorx"
	"github.startlite.cn/itapp/startlite/pkg/lines/featurex"
	"github.startlite.cn/itapp/startlite/pkg/lines/logx"
	"github.startlite.cn/itapp/startlite/pkg/metrics"
	"github.startlite.cn/itapp/startlite/pkg/utils"
)

var (
	instance *assetClient
)

type AssetClient interface {
	QuerySitesData(ctx context.Context, siteId, din, externalSiteId, assetSiteID, requestId string) (sr interface{}, err error)
	getEnergyAuthToken() (bearerToken string, err error)
}

type assetClient struct {
	reqCtx      appx.ReqContext
	assetAPI    *metrics.HTTPClient
	tokenServer string
	assetServer string
}

func MustNewAssetClient(appCtx appx.AppContext, cl *featurex.ConfigLoader) AssetClient {

	assetConfig := &types.AssetApi{}
	cl.Load(assetConfig)
	tokenConfig := &types.TokenApi{}
	cl.Load(tokenConfig)

	client, _ := baseclient.SetupHttpClient(assetConfig.Cert, assetConfig.Key)

	instance = &assetClient{
		assetAPI:    client,
		tokenServer: tokenConfig.TokenServer,
		assetServer: assetConfig.AssetServer,
	}

	return instance

}

func (a *assetClient) QuerySitesData(ctx context.Context, siteId, din, externalSiteId, assetSiteID, requestId string) (sr interface{}, err error) {
	params := AssetQuery{
		ExternalSiteID: externalSiteId,
		SiteNumber:     siteId,
		Din:            din,
		AssetSiteID:    assetSiteID,
	}

	urlQuery, err := clients.QueryFromInterface(params)
	if err != nil {
		return nil, err
	}
	url, err := url2.ParseRequestURI(a.assetServer + ASSET_SITES_PATH)
	if err != nil {
		return nil, err
	}

	req, err := clients.NewRequestWithPromContextSpecial(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s://%s", url.Scheme, url.Host),
		url.Path,
		nil,
		urlQuery,
		nil,
	)
	if err != nil {
		return sr, &clients.RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	//log().Info("GET ", req.URL.String())

	var resp *http.Response
	var respBody []byte
	token, err := a.getEnergyAuthToken()
	if err != nil {
		return nil, errorx.Errorf("QuerySitesData: Error getting token: %v", err)
	}
	req.Header.Add(`Authorization`, token)
	resp, err = a.assetAPI.Do(req)
	defer utils.CloseAndDiscardRespBody(resp)
	if err != nil {
		statusCode := http.StatusGatewayTimeout
		err = errorx.Errorf("Error making request; Err: %v", err)
		return sr, &clients.RequestError{
			StatusCode: statusCode,
			Err:        err,
		}
	}

	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		err = errorx.Errorf("Error reading request; Err: %v", err)
		return sr, &clients.RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	// This accounts for the case where we only get the { error: ... } field
	if resp.StatusCode != http.StatusOK {
		err = errorx.Errorf("Unexpected status code received from site data query: %d", resp.StatusCode)
		return sr, &clients.RequestError{
			StatusCode: resp.StatusCode,
			Err:        err,
		}
	}

	var siteResponse SiteResponse
	err = json.Unmarshal(respBody, &siteResponse)
	if err == nil {
		// log().Info("the response json string data  is", string(respBody))
		// log().Info("the response data  is", siteResponse)
		log().Info("getting the site info ============== ")
		return siteResponse, nil
	}

	var siteResponseArray SiteResponseArray
	err = json.Unmarshal(respBody, &siteResponseArray)
	if err == nil {
		log().Info("the response data array is", siteResponseArray)
		return siteResponseArray, nil
	}

	return sr, &clients.RequestError{
		StatusCode: resp.StatusCode,
		Err:        errorx.Errorf("Unable to unmarshal asset response: %v", string(respBody)),
	}
}

// Energy Auth token needed to query Asset API Too Endpoints
func (a *assetClient) getEnergyAuthToken() (bearerToken string, err error) {
	url, err := url2.ParseRequestURI(a.tokenServer + DEFAULT_ENERGY_AUTH_TOKEN_PATH)
	if err != nil {
		err = errorx.Errorf("unable to parse energy auth endpoint")
		return
	}
	r, err := clients.NewRequestWithProm(
		http.MethodPost,
		fmt.Sprintf("%s://%s", url.Scheme, url.Host),
		url.Path,
		nil,
		nil,
		nil,
	)
	if err != nil {

		logx.Warn("Error in getEnergyAuthToken. Unable to generate token request --(Error)--> ", err.Error())
		return
	}

	//log().Info("POST", url.Path)

	r.Header.Add(`Content-Type`, "application/json")
	var resp *http.Response

	resp, err = a.assetAPI.Do(r)
	defer utils.CloseAndDiscardRespBody(resp)
	if err != nil {
		err = errorx.Errorf("Unable to retrieve critical resource --(Error)--> [[ %v ]]", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = errorx.Errorf("Received status code (%d) while retrieving token", resp.StatusCode)
		return
	}

	var tokenRes TokenResponse
	if e := json.NewDecoder(resp.Body).Decode(&tokenRes); e != nil {
		err = errorx.Errorf("Unable to read JSON Auth resource --(Error)--> [[ %v ]]", e)
		return
	}

	bearerToken = "Bearer " + tokenRes.Data.Token
	return
}
func log() logx.ReqLogger {
	return instance.reqCtx.GetLogger()
}
