package cpopartnersdk

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
	// NotificationChargeOrderInfo request  with any body
	NotificationChargeOrderInfoWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	NotificationChargeOrderInfo(ctx context.Context, body NotificationChargeOrderInfoJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// NotificationEquipChargeStatus request  with any body
	NotificationEquipChargeStatusWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	NotificationEquipChargeStatus(ctx context.Context, body NotificationEquipChargeStatusJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// NotificationStartChargeResult request  with any body
	NotificationStartChargeResultWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	NotificationStartChargeResult(ctx context.Context, body NotificationStartChargeResultJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// NotificationStationFee request  with any body
	NotificationStationFeeWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	NotificationStationFee(ctx context.Context, body NotificationStationFeeJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// NotificationStopChargeResult request  with any body
	NotificationStopChargeResultWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	NotificationStopChargeResult(ctx context.Context, body NotificationStopChargeResultJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// QueryToken request  with any body
	QueryTokenWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	QueryToken(ctx context.Context, body QueryTokenJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) NotificationChargeOrderInfoWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewNotificationChargeOrderInfoRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}

	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req.Clone(ctx))
}

func (c *Client) NotificationChargeOrderInfo(ctx context.Context, body NotificationChargeOrderInfoJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewNotificationChargeOrderInfoRequest(c.Server, body)
	if err != nil {
		return nil, err
	}

	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req.Clone(ctx))
}

func (c *Client) NotificationEquipChargeStatusWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewNotificationEquipChargeStatusRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}

	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req.Clone(ctx))
}

func (c *Client) NotificationEquipChargeStatus(ctx context.Context, body NotificationEquipChargeStatusJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewNotificationEquipChargeStatusRequest(c.Server, body)
	if err != nil {
		return nil, err
	}

	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req.Clone(ctx))
}

func (c *Client) NotificationStartChargeResultWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewNotificationStartChargeResultRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}

	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req.Clone(ctx))
}

func (c *Client) NotificationStartChargeResult(ctx context.Context, body NotificationStartChargeResultJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewNotificationStartChargeResultRequest(c.Server, body)
	if err != nil {
		return nil, err
	}

	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req.Clone(ctx))
}

func (c *Client) NotificationStationFeeWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewNotificationStationFeeRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}

	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req.Clone(ctx))
}

func (c *Client) NotificationStationFee(ctx context.Context, body NotificationStationFeeJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewNotificationStationFeeRequest(c.Server, body)
	if err != nil {
		return nil, err
	}

	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req.Clone(ctx))
}

func (c *Client) NotificationStopChargeResultWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewNotificationStopChargeResultRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}

	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req.Clone(ctx))
}

func (c *Client) NotificationStopChargeResult(ctx context.Context, body NotificationStopChargeResultJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewNotificationStopChargeResultRequest(c.Server, body)
	if err != nil {
		return nil, err
	}

	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req.Clone(ctx))
}

func (c *Client) QueryTokenWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewQueryTokenRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}

	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req.Clone(ctx))
}

func (c *Client) QueryToken(ctx context.Context, body QueryTokenJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewQueryTokenRequest(c.Server, body)
	if err != nil {
		return nil, err
	}

	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req.Clone(ctx))
}

// NewNotificationChargeOrderInfoRequest calls the generic NotificationChargeOrderInfo builder with application/json body
func NewNotificationChargeOrderInfoRequest(server string, body NotificationChargeOrderInfoJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewNotificationChargeOrderInfoRequestWithBody(server, "application/json", bodyReader)
}

// NewNotificationChargeOrderInfoRequestWithBody generates requests for NotificationChargeOrderInfo with any type of body
func NewNotificationChargeOrderInfoRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	queryUrl, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	basePath := "/notification_charge_order_info"

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

// NewNotificationEquipChargeStatusRequest calls the generic NotificationEquipChargeStatus builder with application/json body
func NewNotificationEquipChargeStatusRequest(server string, body NotificationEquipChargeStatusJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewNotificationEquipChargeStatusRequestWithBody(server, "application/json", bodyReader)
}

// NewNotificationEquipChargeStatusRequestWithBody generates requests for NotificationEquipChargeStatus with any type of body
func NewNotificationEquipChargeStatusRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	queryUrl, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	basePath := "/notification_equip_charge_status"

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

// NewNotificationStartChargeResultRequest calls the generic NotificationStartChargeResult builder with application/json body
func NewNotificationStartChargeResultRequest(server string, body NotificationStartChargeResultJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewNotificationStartChargeResultRequestWithBody(server, "application/json", bodyReader)
}

// NewNotificationStartChargeResultRequestWithBody generates requests for NotificationStartChargeResult with any type of body
func NewNotificationStartChargeResultRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	queryUrl, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	basePath := "/notification_start_charge_result"

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

// NewNotificationStationFeeRequest calls the generic NotificationStationFee builder with application/json body
func NewNotificationStationFeeRequest(server string, body NotificationStationFeeJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewNotificationStationFeeRequestWithBody(server, "application/json", bodyReader)
}

// NewNotificationStationFeeRequestWithBody generates requests for NotificationStationFee with any type of body
func NewNotificationStationFeeRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	queryUrl, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	basePath := "/notification_stationFee"

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

// NewNotificationStopChargeResultRequest calls the generic NotificationStopChargeResult builder with application/json body
func NewNotificationStopChargeResultRequest(server string, body NotificationStopChargeResultJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewNotificationStopChargeResultRequestWithBody(server, "application/json", bodyReader)
}

// NewNotificationStopChargeResultRequestWithBody generates requests for NotificationStopChargeResult with any type of body
func NewNotificationStopChargeResultRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	queryUrl, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	basePath := "/notification_stop_charge_result"

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

// NewQueryTokenRequest calls the generic queryToken builder with application/json body
func NewQueryTokenRequest(server string, body QueryTokenJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewQueryTokenRequestWithBody(server, "application/json", bodyReader)
}

// NewQueryTokenRequestWithBody generates requests for queryToken with any type of body
func NewQueryTokenRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
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
	// NotificationChargeOrderInfo request  with any body
	NotificationChargeOrderInfoWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*NotificationChargeOrderInfoRespByPen, error)

	NotificationChargeOrderInfoWithResponse(ctx context.Context, body NotificationChargeOrderInfoJSONRequestBody) (*NotificationChargeOrderInfoRespByPen, error)

	// NotificationEquipChargeStatus request  with any body
	NotificationEquipChargeStatusWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*NotificationEquipChargeStatusRespByPen, error)

	NotificationEquipChargeStatusWithResponse(ctx context.Context, body NotificationEquipChargeStatusJSONRequestBody) (*NotificationEquipChargeStatusRespByPen, error)

	// NotificationStartChargeResult request  with any body
	NotificationStartChargeResultWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*NotificationStartChargeResultRespByPen, error)

	NotificationStartChargeResultWithResponse(ctx context.Context, body NotificationStartChargeResultJSONRequestBody) (*NotificationStartChargeResultRespByPen, error)

	// NotificationStationFee request  with any body
	NotificationStationFeeWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*NotificationStationFeeRespByPen, error)

	NotificationStationFeeWithResponse(ctx context.Context, body NotificationStationFeeJSONRequestBody) (*NotificationStationFeeRespByPen, error)

	// NotificationStopChargeResult request  with any body
	NotificationStopChargeResultWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*NotificationStopChargeResultRespByPen, error)

	NotificationStopChargeResultWithResponse(ctx context.Context, body NotificationStopChargeResultJSONRequestBody) (*NotificationStopChargeResultRespByPen, error)

	// queryToken request  with any body
	QueryTokenWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*QueryTokenRespByPen, error)

	QueryTokenWithResponse(ctx context.Context, body QueryTokenJSONRequestBody) (*QueryTokenRespByPen, error)
}

type NotificationChargeOrderInfoRespByPen struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *PartnerTCECResponse
}

// Status returns HTTPResponse.Status
func (r NotificationChargeOrderInfoRespByPen) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r NotificationChargeOrderInfoRespByPen) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetBody returns HTTPResponse payload
func (r NotificationChargeOrderInfoRespByPen) GetBody() []byte {
	if r.HTTPResponse != nil {
		return r.Body
	}
	return nil
}

type NotificationEquipChargeStatusRespByPen struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *PartnerTCECResponse
}

// Status returns HTTPResponse.Status
func (r NotificationEquipChargeStatusRespByPen) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r NotificationEquipChargeStatusRespByPen) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetBody returns HTTPResponse payload
func (r NotificationEquipChargeStatusRespByPen) GetBody() []byte {
	if r.HTTPResponse != nil {
		return r.Body
	}
	return nil
}

type NotificationStartChargeResultRespByPen struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *PartnerTCECResponse
}

// Status returns HTTPResponse.Status
func (r NotificationStartChargeResultRespByPen) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r NotificationStartChargeResultRespByPen) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetBody returns HTTPResponse payload
func (r NotificationStartChargeResultRespByPen) GetBody() []byte {
	if r.HTTPResponse != nil {
		return r.Body
	}
	return nil
}

type NotificationStationFeeRespByPen struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *PartnerTCECResponse
}

// Status returns HTTPResponse.Status
func (r NotificationStationFeeRespByPen) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r NotificationStationFeeRespByPen) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetBody returns HTTPResponse payload
func (r NotificationStationFeeRespByPen) GetBody() []byte {
	if r.HTTPResponse != nil {
		return r.Body
	}
	return nil
}

type NotificationStopChargeResultRespByPen struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *PartnerTCECResponse
}

// Status returns HTTPResponse.Status
func (r NotificationStopChargeResultRespByPen) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r NotificationStopChargeResultRespByPen) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetBody returns HTTPResponse payload
func (r NotificationStopChargeResultRespByPen) GetBody() []byte {
	if r.HTTPResponse != nil {
		return r.Body
	}
	return nil
}

type QueryTokenRespByPen struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *PartnerTCECResponse
}

// Status returns HTTPResponse.Status
func (r QueryTokenRespByPen) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r QueryTokenRespByPen) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetBody returns HTTPResponse payload
func (r QueryTokenRespByPen) GetBody() []byte {
	if r.HTTPResponse != nil {
		return r.Body
	}
	return nil
}

// NotificationChargeOrderInfoWithBodyWithResponse request with arbitrary body returning *NotificationChargeOrderInfoRespByPen
func (c *ClientWithResponses) NotificationChargeOrderInfoWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*NotificationChargeOrderInfoRespByPen, error) {
	rsp, err := c.NotificationChargeOrderInfoWithBody(ctx, contentType, body)
	if err != nil {
		return nil, err
	}
	return ParseNotificationChargeOrderInfoRespByPen(rsp)
}

func (c *ClientWithResponses) NotificationChargeOrderInfoWithResponse(ctx context.Context, body NotificationChargeOrderInfoJSONRequestBody) (*NotificationChargeOrderInfoRespByPen, error) {
	rsp, err := c.NotificationChargeOrderInfo(ctx, body)
	if err != nil {
		return nil, err
	}
	return ParseNotificationChargeOrderInfoRespByPen(rsp)
}

// NotificationEquipChargeStatusWithBodyWithResponse request with arbitrary body returning *NotificationEquipChargeStatusRespByPen
func (c *ClientWithResponses) NotificationEquipChargeStatusWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*NotificationEquipChargeStatusRespByPen, error) {
	rsp, err := c.NotificationEquipChargeStatusWithBody(ctx, contentType, body)
	if err != nil {
		return nil, err
	}
	return ParseNotificationEquipChargeStatusRespByPen(rsp)
}

func (c *ClientWithResponses) NotificationEquipChargeStatusWithResponse(ctx context.Context, body NotificationEquipChargeStatusJSONRequestBody) (*NotificationEquipChargeStatusRespByPen, error) {
	rsp, err := c.NotificationEquipChargeStatus(ctx, body)
	if err != nil {
		return nil, err
	}
	return ParseNotificationEquipChargeStatusRespByPen(rsp)
}

// NotificationStartChargeResultWithBodyWithResponse request with arbitrary body returning *NotificationStartChargeResultRespByPen
func (c *ClientWithResponses) NotificationStartChargeResultWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*NotificationStartChargeResultRespByPen, error) {
	rsp, err := c.NotificationStartChargeResultWithBody(ctx, contentType, body)
	if err != nil {
		return nil, err
	}
	return ParseNotificationStartChargeResultRespByPen(rsp)
}

func (c *ClientWithResponses) NotificationStartChargeResultWithResponse(ctx context.Context, body NotificationStartChargeResultJSONRequestBody) (*NotificationStartChargeResultRespByPen, error) {
	rsp, err := c.NotificationStartChargeResult(ctx, body)
	if err != nil {
		return nil, err
	}
	return ParseNotificationStartChargeResultRespByPen(rsp)
}

// NotificationStationFeeWithBodyWithResponse request with arbitrary body returning *NotificationStationFeeRespByPen
func (c *ClientWithResponses) NotificationStationFeeWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*NotificationStationFeeRespByPen, error) {
	rsp, err := c.NotificationStationFeeWithBody(ctx, contentType, body)
	if err != nil {
		return nil, err
	}
	return ParseNotificationStationFeeRespByPen(rsp)
}

func (c *ClientWithResponses) NotificationStationFeeWithResponse(ctx context.Context, body NotificationStationFeeJSONRequestBody) (*NotificationStationFeeRespByPen, error) {
	rsp, err := c.NotificationStationFee(ctx, body)
	if err != nil {
		return nil, err
	}
	return ParseNotificationStationFeeRespByPen(rsp)
}

// NotificationStopChargeResultWithBodyWithResponse request with arbitrary body returning *NotificationStopChargeResultRespByPen
func (c *ClientWithResponses) NotificationStopChargeResultWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*NotificationStopChargeResultRespByPen, error) {
	rsp, err := c.NotificationStopChargeResultWithBody(ctx, contentType, body)
	if err != nil {
		return nil, err
	}
	return ParseNotificationStopChargeResultRespByPen(rsp)
}

func (c *ClientWithResponses) NotificationStopChargeResultWithResponse(ctx context.Context, body NotificationStopChargeResultJSONRequestBody) (*NotificationStopChargeResultRespByPen, error) {
	rsp, err := c.NotificationStopChargeResult(ctx, body)
	if err != nil {
		return nil, err
	}
	return ParseNotificationStopChargeResultRespByPen(rsp)
}

// QueryTokenWithBodyWithResponse request with arbitrary body returning *QueryTokenRespByPen
func (c *ClientWithResponses) QueryTokenWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*QueryTokenRespByPen, error) {
	rsp, err := c.QueryTokenWithBody(ctx, contentType, body)
	if err != nil {
		return nil, err
	}
	return ParseQueryTokenRespByPen(rsp)
}

func (c *ClientWithResponses) QueryTokenWithResponse(ctx context.Context, body QueryTokenJSONRequestBody) (*QueryTokenRespByPen, error) {
	rsp, err := c.QueryToken(ctx, body)
	if err != nil {
		return nil, err
	}
	return ParseQueryTokenRespByPen(rsp)
}

// ParseNotificationChargeOrderInfoRespByPen parses an HTTP response from a NotificationChargeOrderInfoWithResponse call
func ParseNotificationChargeOrderInfoRespByPen(rsp *http.Response) (*NotificationChargeOrderInfoRespByPen, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &NotificationChargeOrderInfoRespByPen{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest PartnerTCECResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseNotificationEquipChargeStatusRespByPen parses an HTTP response from a NotificationEquipChargeStatusWithResponse call
func ParseNotificationEquipChargeStatusRespByPen(rsp *http.Response) (*NotificationEquipChargeStatusRespByPen, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &NotificationEquipChargeStatusRespByPen{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest PartnerTCECResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseNotificationStartChargeResultRespByPen parses an HTTP response from a NotificationStartChargeResultWithResponse call
func ParseNotificationStartChargeResultRespByPen(rsp *http.Response) (*NotificationStartChargeResultRespByPen, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &NotificationStartChargeResultRespByPen{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest PartnerTCECResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseNotificationStationFeeRespByPen parses an HTTP response from a NotificationStationFeeWithResponse call
func ParseNotificationStationFeeRespByPen(rsp *http.Response) (*NotificationStationFeeRespByPen, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &NotificationStationFeeRespByPen{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest PartnerTCECResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseNotificationStopChargeResultRespByPen parses an HTTP response from a NotificationStopChargeResultWithResponse call
func ParseNotificationStopChargeResultRespByPen(rsp *http.Response) (*NotificationStopChargeResultRespByPen, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &NotificationStopChargeResultRespByPen{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest PartnerTCECResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseQueryTokenRespByPen parses an HTTP response from a QueryTokenWithResponse call
func ParseQueryTokenRespByPen(rsp *http.Response) (*QueryTokenRespByPen, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &QueryTokenRespByPen{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest PartnerTCECResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}
