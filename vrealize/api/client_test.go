package api

import (
	"testing"
	"gopkg.in/jarcoal/httpmock.v1"
	"github.com/stretchr/testify/assert"
	"errors"
	"io/ioutil"
)

var client Client

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

	if client.Username != username {
		t.Errorf("Expected username %v, got %v ", username, client.Username)
	}

	if client.Password != password {
		t.Errorf("Expected password %v, got %v ", password, client.Password)
	}

	if client.Tenant != tenant {
		t.Errorf("Expected tenant %v, got %v ", tenant, client.Tenant)
	}

	if client.BaseURL != baseURL {
		t.Errorf("Expected BaseUrl %v, got %v ", baseURL, client.BaseURL)
	}
}

func TestAPIClient_Authenticate_OK(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://localhost/identity/api/tokens",
		httpmock.NewStringResponder(200, testData("token_response.json")))

	err := client.Authenticate()

	assert.Nil(t, err)

	if len(client.BearerToken) == 0 {
		t.Error("Fail to set BearerToken.")
	}

}

func TestAPIClient_Authenticate_Failed(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://localhost/identity/api/tokens",
		httpmock.NewErrorResponder(errors.New(testData("auth_failure.json"))))

	err := client.Authenticate()

	if err == nil {
		t.Errorf("Authentication should fail")
	}

}

func testData(filename string) (string) {
	bytes, err := ioutil.ReadFile("test_data/" + filename)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}