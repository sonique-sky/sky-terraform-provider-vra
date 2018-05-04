package api

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/dghubble/sling"
	"fmt"
)

const (
	fmtRequestBase = "/catalog-service/api/consumer/requests"
	fmtRequest = fmtRequestBase + "/%s"
	fmtRequestResourceViews = fmtRequestBase + "/%s/resourceViews"

	fmtResourcesBase = "/catalog-service/api/consumer/resources"
	fmtActionTemplate = fmtResourcesBase + "/%s/actions/%s/requests/template"
	fmtActionRequest = fmtResourcesBase + "/%s/actions/%s/requests"

	fmtCatalogItemsTemplate = "/catalog-service/api/consumer/entitledCatalogItems/%s/requests/template"
)
//Client - This struct is used to store information provided in .tf file under provider block
//Later on, this stores bearToken after successful authentication and uses that token for next
//REST get or post calls.
type Client struct {
	Username    string
	Password    string
	BaseURL     string
	Tenant      string
	Insecure    bool
	BearerToken string
	HTTPClient  *sling.Sling
}

//AuthRequest - This struct contains the user information provided by user
//and for authentication details of this struct are used.
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Tenant   string `json:"tenant"`
}

//AuthResponse - This struct contains response of user authentication call.
type AuthResponse struct {
	Expires time.Time `json:"expires"`
	ID      string    `json:"id"`
	Tenant  string    `json:"tenant"`
}

func NewClient(username string, password string, tenant string, baseURL string, insecure bool) Client {
	// This overrides the DefaultTransport which is probably ok
	// since we're generally only using a single client.
	transport := http.DefaultTransport.(*http.Transport)
	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: insecure,
	}
	return Client{
		Username: username,
		Password: password,
		Tenant:   tenant,
		BaseURL:  baseURL,
		Insecure: insecure,
		HTTPClient: sling.New().Base(baseURL).
			Set("Accept", "application/json").
			Set("Content-Type", "application/json"),
	}
}

func (c *Client) Authenticate() error {
	params := &AuthRequest{
		Username: c.Username,
		Password: c.Password,
		Tenant:   c.Tenant,
	}

	authRes := new(AuthResponse)

	err := c.post("/identity/api/tokens", params, authRes, noCheck)

	if err != nil {
		return err
	}

	c.BearerToken = authRes.ID
	c.HTTPClient = c.HTTPClient.New().Set("Authorization", fmt.Sprintf("Bearer %s", authRes.ID))

	return nil
}

func (c *Client) post(requestUrl string, requestBody interface{}, response interface{}, validate func(*http.Response) bool) error {
	apiError := new(Error)
	resp, err := c.HTTPClient.New().Post(requestUrl).BodyJSON(requestBody).Receive(response, apiError)

	if err != nil {
		return err
	}

	if !validate(resp) {
		return err
	}

	if !apiError.isEmpty() {
		return apiError.Error()
	}

	return nil
}

func (c *Client) get(requestUrl string, response interface{}, validate func(*http.Response) bool) (error) {
	apiError := new(Error)

	resp, err := c.HTTPClient.New().Get(requestUrl).Receive(response, apiError)

	if err != nil {
		return err
	}

	if !validate(resp) {
		return err
	}

	if !apiError.isEmpty() {
		return apiError.Error()
	}

	return nil
}

func expectHttpStatus(code int) (func(response *http.Response) bool) {
	return func(resp *http.Response) bool {
		return code == resp.StatusCode
	}
}

var noCheck = func(response *http.Response) bool {
	return true
}
