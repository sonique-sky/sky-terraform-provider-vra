package api

import (
	"testing"
	"gopkg.in/jarcoal/httpmock.v1"
	"github.com/stretchr/testify/assert"
	"errors"
	"io/ioutil"
	"encoding/json"
)

var client RestClient

func init() {
	client = NewClient(
		"admin@myvra.local",
		"pass!@#",
		"vsphere.local",
		"http://localhost/",
		true,
	)
}

func TestNewAPIClient(t *testing.T) {
	username := "admin@myvra.local"
	password := "pass!@#"
	tenant := "vshpere.local"
	baseURL := "http://localhost/"

	client := NewClient(
		username,
		password,
		tenant,
		baseURL,
		true,
	)

	assert.Equal(t, username, client.Username)
	assert.Equal(t, password, client.Password)
	assert.Equal(t, tenant, client.Tenant)
	assert.Equal(t, baseURL, client.BaseURL)
}

func TestAPIClient_Authenticate_OK(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://localhost/identity/api/tokens",
		httpmock.NewStringResponder(200, testData("token_response")))

	data, _ := ioutil.ReadFile("test_data/token_response.json")
	response := new(AuthResponse)
	json.Unmarshal(data, response)

	err := client.Authenticate()

	assert.Nil(t, err)
	assert.Equal(t, response.ID, client.BearerToken)
}

func TestAPIClient_Authenticate_Failed(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://localhost/identity/api/tokens",
		httpmock.NewErrorResponder(errors.New(testData("auth_failure"))))

	err := client.Authenticate()

	assert.NotNil(t, err)
}

func testData(filename string) (string) {
	bytes, err := ioutil.ReadFile("test_data/" + filename +".json")
	if err != nil {
		panic(err)
	}
	return string(bytes)
}