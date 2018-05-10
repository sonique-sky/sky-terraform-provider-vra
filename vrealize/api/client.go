package api

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/dghubble/sling"
	"fmt"
	"log"
)

const (
	fmtRequestBase          = "/catalog-service/api/consumer/requests"
	fmtRequest              = fmtRequestBase + "/%s"
	fmtRequestResourceViews = fmtRequestBase + "/%s/resourceViews"

	fmtCatalogItemsSearch = "/catalog-service/api/consumer/entitledCatalogItems?$filter=name+eq+'%s'"
)


type Client interface {
	GetRequestStatus(requestId string) (*RequestStatusView, error)
	GetResourceViews(requestId string) (*ResourceViewsTemplate, error)
	GetMachine(resourceId string) (*Resource, error)

	ReadCatalogByID(catalogId string) (*CatalogItemTemplate, error)
	ReadCatalogByName(catalogName string) (*CatalogItemTemplate, error)

	RequestMachine(template *CatalogItemTemplate) (*RequestMachineResponse, error)
	DestroyMachine(resourceId string) (error)
}

type RestClient struct {
	Username    string
	Password    string
	BaseURL     string
	Tenant      string
	Insecure    bool
	BearerToken string
	HTTPClient  *sling.Sling
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Tenant   string `json:"tenant"`
}

type AuthResponse struct {
	Expires time.Time `json:"expires"`
	ID      string    `json:"id"`
	Tenant  string    `json:"tenant"`
}

func NewClient(username string, password string, tenant string, baseURL string, insecure bool) *RestClient {
	transport := http.DefaultTransport.(*http.Transport)
	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: insecure,
	}
	return &RestClient{
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

func (c *RestClient) Authenticate() error {
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

func (c *RestClient) post(requestUrl string, requestBody interface{}, response interface{}, validate func(*http.Response) bool) error {
	log.Printf("POST : %s", requestUrl)
	apiError := new(Error)
	resp, err := c.HTTPClient.New().Post(requestUrl).BodyJSON(requestBody).Receive(response, apiError)

	if err != nil {
		log.Printf("Http Error : %v", err)
		return err
	}

	if !validate(resp) {
		return err
	}

	if !apiError.isEmpty() {
		log.Printf("Api Error : %v", apiError.Error())
		return apiError.Error()
	}

	return nil
}

func (c *RestClient) get(requestUrl string, response interface{}, validate func(*http.Response) bool) (error) {
	log.Printf("GET: %s", requestUrl)
	apiError := new(Error)

	resp, err := c.HTTPClient.New().Get(requestUrl).Receive(response, apiError)

	if err != nil {
		log.Printf("Http Error : %v", err)
		return err
	}

	if !validate(resp) {
		return err
	}

	if !apiError.isEmpty() {
		log.Printf("Api Error : %v", apiError.Error())
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
