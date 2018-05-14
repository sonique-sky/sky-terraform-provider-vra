package api

import (
	"testing"
	"gopkg.in/jarcoal/httpmock.v1"
	"fmt"
	"github.com/stretchr/testify/assert"
	"errors"
	"io/ioutil"
	"encoding/json"
)

func TestAPIClient_RequestMachine(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	catalogId := "e5dd4fba-45ed-4943-b1fc-7f96239286be"
	template := catalogItem()

	path := fmt.Sprintf("http://localhost/catalog-service/api/consumer/entitledCatalogItems/%s/requests", catalogId)
	httpmock.RegisterResponder("POST", path,
		httpmock.NewStringResponder(201, testData("request_list")))

	requestMachine, errorRequestMachine := client.RequestMachine(template)

	assert.Nil(t, errorRequestMachine)
	assert.Equal(t, "b2907df7-6c36-4e30-9c62-a21f293b067a", requestMachine.ID)
}

func TestAPIClient_RequestMachine_Failure(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	catalogId := "e5dd4fba-45ed-4943-b1fc-7f96239286be"
	template := catalogItem()

	path := fmt.Sprintf("http://localhost/catalog-service/api/consumer/entitledCatalogItems/%s/requests", catalogId)
	httpmock.RegisterResponder("POST", path,
		httpmock.NewErrorResponder(errors.New(testData("api_error"))))

	requestMachine, errorRequestMachine := client.RequestMachine(template)

	assert.NotNil(t, errorRequestMachine)
	assert.Nil(t, requestMachine)
}

func TestAPIClient_DestroyMachine(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	resourceId := "b313acd6-0738-439c-b601-e3ebf9ebb49b"
	actionId := "5883bbea-cd9a-4bf3-a2b3-f6cc20135435"

	httpmock.RegisterResponder("GET", fmt.Sprintf("http://localhost/catalog-service/api/consumer/resources/%s/actions", resourceId),
		httpmock.NewStringResponder(200, testData("action_list")))

	httpmock.RegisterResponder("GET", fmt.Sprintf("http://localhost/catalog-service/api/consumer/resources/%s/actions/%s/requests/template", resourceId, actionId),
		httpmock.NewStringResponder(200, testData("destroy_template")))


	httpmock.RegisterResponder("POST", fmt.Sprintf("http://localhost/catalog-service/api/consumer/resources/%s/actions/%s/requests", resourceId, actionId),
		httpmock.NewStringResponder(201, "{}"))

	err := client.DestroyMachine(resourceId)

	assert.Nil(t, err)
}

func TestAPIClient_DestroyMachine_NotFound(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	resourceId := "b313acd6-0738-439c-b601-e3ebf9ebb49b"
	actionId := "5883bbea-cd9a-4bf3-a2b3-f6cc20135435"


	httpmock.RegisterResponder("GET", fmt.Sprintf("http://localhost/catalog-service/api/consumer/resources/%s/actions", resourceId),
		httpmock.NewStringResponder(200, testData("action_list")))

	path := fmt.Sprintf("http://localhost/catalog-service/api/consumer/resources/%s/actions/%s/requests/template", resourceId, actionId)
	httpmock.RegisterResponder("GET", path,
		httpmock.NewStringResponder(200, testData("destroy_template")))

	httpmock.RegisterResponder("POST", fmt.Sprintf("http://localhost"+("/catalog-service/api/consumer/resources")+"/%s/actions/%s/requests", resourceId, actionId),
		httpmock.NewStringResponder(201, "{}"))

	err := client.DestroyMachine(resourceId)

	assert.Nil(t, err)
}

func TestClient_GetMachine(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	resourceId := "dd3ad4bc-f7f2-46dd-bc31-3bb3c1ea460c"

	httpmock.RegisterResponder("GET", fmt.Sprintf("http://localhost/catalog-service/api/consumer/resources/%s", resourceId),
		httpmock.NewStringResponder(200, testData("resource")))

	resource, e := client.GetMachine(resourceId)

	assert.Nil(t, e)

	assert.Equal(t, "vm007203", resource.Name)

	ipAddress, _ := resource.StringValue("ip_address")
	assert.Equal(t, "10.90.64.29", ipAddress)
}

func catalogItem() *RequestTemplate {
	resourceViewsTemplate := new(RequestTemplate)
	data, _ := ioutil.ReadFile("test_data/catalog_item_request_template.json")
	json.Unmarshal(data, resourceViewsTemplate)
	return resourceViewsTemplate
}
