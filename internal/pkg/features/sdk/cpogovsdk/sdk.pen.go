package cpogovsdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.startlite.cn/itapp/startlite/pkg/lines/appx"
	"github.startlite.cn/itapp/startlite/pkg/lines/constantx"
	"github.startlite.cn/itapp/startlite/pkg/lines/errorx"
	"github.startlite.cn/itapp/startlite/pkg/lines/utilx/textx"
	"github.startlite.cn/itapp/startlite/pkg/lines/utilx/tokenx"
	"github.startlite.cn/itapp/startlite/pkg/lines/xidx"
)

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = http.DefaultClient
	}

	client.RequestEditors = append(client.RequestEditors, withSetXidFromRequest, withSetTokenFromRequest)

	// ensure will must occur deep copy on slice in WithRequestEditorFn action later
	client.RequestEditors = client.RequestEditors[0:len(client.RequestEditors):len(client.RequestEditors)]
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)

		// ensure will must occur deep copy on slice in WithRequestEditorFn action later
		c.RequestEditors = c.RequestEditors[0:len(c.RequestEditors):len(c.RequestEditors)]
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// QueryStationStatsUsingPOST request  with any body
	QueryStationStatsUsingPOSTWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	QueryStationStatsUsingPOST(ctx context.Context, body QueryStationStatsUsingPOSTJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// QueryStationStatusUsingPOST request  with any body
	QueryStationStatusUsingPOSTWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	QueryStationStatusUsingPOST(ctx context.Context, body QueryStationStatusUsingPOSTJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// QueryStationInfoUsingPOST request  with any body
	QueryStationInfoUsingPOSTWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	QueryStationInfoUsingPOST(ctx context.Context, body QueryStationInfoUsingPOSTJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// QueryTokenUsingPOST request  with any body
	QueryTokenUsingPOSTWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	QueryTokenUsingPOST(ctx context.Context, body QueryTokenUsingPOSTJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) QueryStationStatsUsingPOSTWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewQueryStationStatsUsingPOSTRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}

	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req.Clone(ctx))
}

func (c *Client) QueryStationStatsUsingPOST(ctx context.Context, body QueryStationStatsUsingPOSTJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewQueryStationStatsUsingPOSTRequest(c.Server, body)
	if err != nil {
		return nil, err
	}

	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req.Clone(ctx))
}

func (c *Client) QueryStationStatusUsingPOSTWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewQueryStationStatusUsingPOSTRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}

	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req.Clone(ctx))
}

func (c *Client) QueryStationStatusUsingPOST(ctx context.Context, body QueryStationStatusUsingPOSTJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewQueryStationStatusUsingPOSTRequest(c.Server, body)
	if err != nil {
		return nil, err
	}

	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req.Clone(ctx))
}

func (c *Client) QueryStationInfoUsingPOSTWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewQueryStationInfoUsingPOSTRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}

	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req.Clone(ctx))
}

func (c *Client) QueryStationInfoUsingPOST(ctx context.Context, body QueryStationInfoUsingPOSTJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewQueryStationInfoUsingPOSTRequest(c.Server, body)
	if err != nil {
		return nil, err
	}

	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req.Clone(ctx))
}

func (c *Client) QueryTokenUsingPOSTWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewQueryTokenUsingPOSTRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}

	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req.Clone(ctx))
}

func (c *Client) QueryTokenUsingPOST(ctx context.Context, body QueryTokenUsingPOSTJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewQueryTokenUsingPOSTRequest(c.Server, body)
	if err != nil {
		return nil, err
	}

	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req.Clone(ctx))
}

// NewQueryStationStatsUsingPOSTRequest calls the generic QueryStationStatsUsingPOST builder with application/json body
func NewQueryStationStatsUsingPOSTRequest(server string, body QueryStationStatsUsingPOSTJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewQueryStationStatsUsingPOSTRequestWithBody(server, "application/json", bodyReader)
}

// NewQueryStationStatsUsingPOSTRequestWithBody generates requests for QueryStationStatsUsingPOST with any type of body
func NewQueryStationStatsUsingPOSTRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	queryUrl, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	basePath := "/query_station_stats"

	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	queryUrl, err = queryUrl.Parse(basePath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryUrl.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewQueryStationStatusUsingPOSTRequest calls the generic QueryStationStatusUsingPOST builder with application/json body
func NewQueryStationStatusUsingPOSTRequest(server string, body QueryStationStatusUsingPOSTJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewQueryStationStatusUsingPOSTRequestWithBody(server, "application/json", bodyReader)
}

// NewQueryStationStatusUsingPOSTRequestWithBody generates requests for QueryStationStatusUsingPOST with any type of body
func NewQueryStationStatusUsingPOSTRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	queryUrl, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	basePath := "/query_station_status"

	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	queryUrl, err = queryUrl.Parse(basePath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryUrl.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewQueryStationInfoUsingPOSTRequest calls the generic QueryStationInfoUsingPOST builder with application/json body
func NewQueryStationInfoUsingPOSTRequest(server string, body QueryStationInfoUsingPOSTJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewQueryStationInfoUsingPOSTRequestWithBody(server, "application/json", bodyReader)
}

// NewQueryStationInfoUsingPOSTRequestWithBody generates requests for QueryStationInfoUsingPOST with any type of body
func NewQueryStationInfoUsingPOSTRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	queryUrl, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	basePath := "/query_stations_info"

	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	queryUrl, err = queryUrl.Parse(basePath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryUrl.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewQueryTokenUsingPOSTRequest calls the generic QueryTokenUsingPOST builder with application/json body
func NewQueryTokenUsingPOSTRequest(server string, body QueryTokenUsingPOSTJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewQueryTokenUsingPOSTRequestWithBody(server, "application/json", bodyReader)
}

// NewQueryTokenUsingPOSTRequestWithBody generates requests for QueryTokenUsingPOST with any type of body
func NewQueryTokenUsingPOSTRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	queryUrl, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	basePath := "/query_token"

	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	queryUrl, err = queryUrl.Parse(basePath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryUrl.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// WithChangeToInternalApi is a RequestEditorFn
func WithChangeToInternalApi(ctx context.Context, req *http.Request) error {
	if strings.HasPrefix(req.URL.Path, constantx.APIPrefix) {
		req.URL.Path = strings.Replace(req.URL.Path, constantx.APIPrefix, constantx.InternalPrefix, 1)
	}

	return nil
}

// withSetXidFromRequest is a RequestEditorFn to set xid in header
func withSetXidFromRequest(ctx context.Context, req *http.Request) error {
	var xid string
	if v, ok := ctx.(appx.ReqContext); ok {
		xid = v.GetXid()
	}
	if !textx.Blank(xid) {
		req.Header.Set(xidx.HeaderXid, xid)
	}

	return nil
}

// withSetTokenFromRequest is a RequestEditorFn to set Authorization in header
func withSetTokenFromRequest(ctx context.Context, req *http.Request) error {
	var token string
	if v, ok := ctx.(appx.ReqContext); ok {
		token = tokenx.GetToken(v)
	}
	if !textx.Blank(token) {
		req.Header.Set(tokenx.HEADERAuthorization, fmt.Sprintf("%s %s", tokenx.HEADERBearer, token))
	}
	return nil
}

func (cwr ClientWithResponses) ExtendWithOptions(opts ...ClientOption) ClientWithResponsesInterface {
	client, ok := cwr.ClientInterface.(*Client)
	if !ok {
		panic("can't extend client")
	}

	newClient := &Client{
		Server:         client.Server,
		Client:         client.Client,
		RequestEditors: client.RequestEditors,
	}
	for _, o := range opts {
		if err := o(newClient); err != nil {
			panic(errorx.Wrap(err, "can't apply option"))
		}
	}
	return &ClientWithResponses{newClient}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	ExtendWithOptions(options ...ClientOption) ClientWithResponsesInterface
	// QueryStationStatsUsingPOST request  with any body
	QueryStationStatsUsingPOSTWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*QueryStationStatsUsingPOSTRespByPen, error)

	QueryStationStatsUsingPOSTWithResponse(ctx context.Context, body QueryStationStatsUsingPOSTJSONRequestBody) (*QueryStationStatsUsingPOSTRespByPen, error)

	// QueryStationStatusUsingPOST request  with any body
	QueryStationStatusUsingPOSTWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*QueryStationStatusUsingPOSTRespByPen, error)

	QueryStationStatusUsingPOSTWithResponse(ctx context.Context, body QueryStationStatusUsingPOSTJSONRequestBody) (*QueryStationStatusUsingPOSTRespByPen, error)

	// QueryStationInfoUsingPOST request  with any body
	QueryStationInfoUsingPOSTWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*QueryStationInfoUsingPOSTRespByPen, error)

	QueryStationInfoUsingPOSTWithResponse(ctx context.Context, body QueryStationInfoUsingPOSTJSONRequestBody) (*QueryStationInfoUsingPOSTRespByPen, error)

	// QueryTokenUsingPOST request  with any body
	QueryTokenUsingPOSTWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*QueryTokenUsingPOSTRespByPen, error)

	QueryTokenUsingPOSTWithResponse(ctx context.Context, body QueryTokenUsingPOSTJSONRequestBody) (*QueryTokenUsingPOSTRespByPen, error)
}

type QueryStationStatsUsingPOSTRespByPen struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *GovChargingResponseEncryption
}

// Status returns HTTPResponse.Status
func (r QueryStationStatsUsingPOSTRespByPen) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r QueryStationStatsUsingPOSTRespByPen) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetBody returns HTTPResponse payload
func (r QueryStationStatsUsingPOSTRespByPen) GetBody() []byte {
	if r.HTTPResponse != nil {
		return r.Body
	}
	return nil
}

type QueryStationStatusUsingPOSTRespByPen struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *GovChargingResponseEncryption
}

// Status returns HTTPResponse.Status
func (r QueryStationStatusUsingPOSTRespByPen) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r QueryStationStatusUsingPOSTRespByPen) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetBody returns HTTPResponse payload
func (r QueryStationStatusUsingPOSTRespByPen) GetBody() []byte {
	if r.HTTPResponse != nil {
		return r.Body
	}
	return nil
}

type QueryStationInfoUsingPOSTRespByPen struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *GovChargingResponseEncryption
}

// Status returns HTTPResponse.Status
func (r QueryStationInfoUsingPOSTRespByPen) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r QueryStationInfoUsingPOSTRespByPen) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetBody returns HTTPResponse payload
func (r QueryStationInfoUsingPOSTRespByPen) GetBody() []byte {
	if r.HTTPResponse != nil {
		return r.Body
	}
	return nil
}

type QueryTokenUsingPOSTRespByPen struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *GovChargingResponseEncryption
}

// Status returns HTTPResponse.Status
func (r QueryTokenUsingPOSTRespByPen) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r QueryTokenUsingPOSTRespByPen) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetBody returns HTTPResponse payload
func (r QueryTokenUsingPOSTRespByPen) GetBody() []byte {
	if r.HTTPResponse != nil {
		return r.Body
	}
	return nil
}

// QueryStationStatsUsingPOSTWithBodyWithResponse request with arbitrary body returning *QueryStationStatsUsingPOSTRespByPen
func (c *ClientWithResponses) QueryStationStatsUsingPOSTWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*QueryStationStatsUsingPOSTRespByPen, error) {
	rsp, err := c.QueryStationStatsUsingPOSTWithBody(ctx, contentType, body)
	if err != nil {
		return nil, err
	}
	return ParseQueryStationStatsUsingPOSTRespByPen(rsp)
}

func (c *ClientWithResponses) QueryStationStatsUsingPOSTWithResponse(ctx context.Context, body QueryStationStatsUsingPOSTJSONRequestBody) (*QueryStationStatsUsingPOSTRespByPen, error) {
	rsp, err := c.QueryStationStatsUsingPOST(ctx, body)
	if err != nil {
		return nil, err
	}
	return ParseQueryStationStatsUsingPOSTRespByPen(rsp)
}

// QueryStationStatusUsingPOSTWithBodyWithResponse request with arbitrary body returning *QueryStationStatusUsingPOSTRespByPen
func (c *ClientWithResponses) QueryStationStatusUsingPOSTWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*QueryStationStatusUsingPOSTRespByPen, error) {
	rsp, err := c.QueryStationStatusUsingPOSTWithBody(ctx, contentType, body)
	if err != nil {
		return nil, err
	}
	return ParseQueryStationStatusUsingPOSTRespByPen(rsp)
}

func (c *ClientWithResponses) QueryStationStatusUsingPOSTWithResponse(ctx context.Context, body QueryStationStatusUsingPOSTJSONRequestBody) (*QueryStationStatusUsingPOSTRespByPen, error) {
	rsp, err := c.QueryStationStatusUsingPOST(ctx, body)
	if err != nil {
		return nil, err
	}
	return ParseQueryStationStatusUsingPOSTRespByPen(rsp)
}

// QueryStationInfoUsingPOSTWithBodyWithResponse request with arbitrary body returning *QueryStationInfoUsingPOSTRespByPen
func (c *ClientWithResponses) QueryStationInfoUsingPOSTWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*QueryStationInfoUsingPOSTRespByPen, error) {
	rsp, err := c.QueryStationInfoUsingPOSTWithBody(ctx, contentType, body)
	if err != nil {
		return nil, err
	}
	return ParseQueryStationInfoUsingPOSTRespByPen(rsp)
}

func (c *ClientWithResponses) QueryStationInfoUsingPOSTWithResponse(ctx context.Context, body QueryStationInfoUsingPOSTJSONRequestBody) (*QueryStationInfoUsingPOSTRespByPen, error) {
	rsp, err := c.QueryStationInfoUsingPOST(ctx, body)
	if err != nil {
		return nil, err
	}
	return ParseQueryStationInfoUsingPOSTRespByPen(rsp)
}

// QueryTokenUsingPOSTWithBodyWithResponse request with arbitrary body returning *QueryTokenUsingPOSTRespByPen
func (c *ClientWithResponses) QueryTokenUsingPOSTWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*QueryTokenUsingPOSTRespByPen, error) {
	rsp, err := c.QueryTokenUsingPOSTWithBody(ctx, contentType, body)
	if err != nil {
		return nil, err
	}
	return ParseQueryTokenUsingPOSTRespByPen(rsp)
}

func (c *ClientWithResponses) QueryTokenUsingPOSTWithResponse(ctx context.Context, body QueryTokenUsingPOSTJSONRequestBody) (*QueryTokenUsingPOSTRespByPen, error) {
	rsp, err := c.QueryTokenUsingPOST(ctx, body)
	if err != nil {
		return nil, err
	}
	return ParseQueryTokenUsingPOSTRespByPen(rsp)
}

// ParseQueryStationStatsUsingPOSTRespByPen parses an HTTP response from a QueryStationStatsUsingPOSTWithResponse call
func ParseQueryStationStatsUsingPOSTRespByPen(rsp *http.Response) (*QueryStationStatsUsingPOSTRespByPen, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &QueryStationStatsUsingPOSTRespByPen{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest GovChargingResponseEncryption
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseQueryStationStatusUsingPOSTRespByPen parses an HTTP response from a QueryStationStatusUsingPOSTWithResponse call
func ParseQueryStationStatusUsingPOSTRespByPen(rsp *http.Response) (*QueryStationStatusUsingPOSTRespByPen, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &QueryStationStatusUsingPOSTRespByPen{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest GovChargingResponseEncryption
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseQueryStationInfoUsingPOSTRespByPen parses an HTTP response from a QueryStationInfoUsingPOSTWithResponse call
func ParseQueryStationInfoUsingPOSTRespByPen(rsp *http.Response) (*QueryStationInfoUsingPOSTRespByPen, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &QueryStationInfoUsingPOSTRespByPen{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest GovChargingResponseEncryption
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseQueryTokenUsingPOSTRespByPen parses an HTTP response from a QueryTokenUsingPOSTWithResponse call
func ParseQueryTokenUsingPOSTRespByPen(rsp *http.Response) (*QueryTokenUsingPOSTRespByPen, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &QueryTokenUsingPOSTRespByPen{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest GovChargingResponseEncryption
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}
