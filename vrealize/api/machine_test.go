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
	actionId := "3da0ca14-e7e2-4d7b-89cb-c6db57440d72"

	httpmock.RegisterResponder("GET", fmt.Sprintf("http://localhost"+fmtActionTemplate, resourceId, actionId),
		httpmock.NewStringResponder(200, testData("destroy_template")))

	httpmock.RegisterResponder("POST", fmt.Sprintf("http://localhost"+fmtActionRequest, resourceId, actionId),
		httpmock.NewStringResponder(201, "{}"))

	templateResources := resourceViews()

	err := client.DestroyMachine(templateResources)

	assert.Nil(t, err)
}

func TestAPIClient_DestroyMachine_NotFound(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	resourceId := "b313acd6-0738-439c-b601-e3ebf9ebb49b"
	actionId := "3da0ca14-e7e2-4d7b-89cb-c6db57440d72"

	httpmock.RegisterResponder("GET", fmt.Sprintf("http://localhost"+fmtActionTemplate, resourceId, actionId),
		httpmock.NewStringResponder(200, testData("destroy_template")))

	httpmock.RegisterResponder("POST", fmt.Sprintf("http://localhost"+fmtActionRequest, resourceId, actionId),
		httpmock.NewStringResponder(201, "{}"))

	templateResources := resourceViews()

	err := client.DestroyMachine(templateResources)

	assert.Nil(t, err)
}

func resourceViews() *ResourceViewsTemplate {
	resourceViewsTemplate := new(ResourceViewsTemplate)
	data, _ := ioutil.ReadFile("test_data/resource_views.json")
	json.Unmarshal(data, resourceViewsTemplate)
	return resourceViewsTemplate
}

func catalogItem() *CatalogItemTemplate {
	resourceViewsTemplate := new(CatalogItemTemplate)
	data, _ := ioutil.ReadFile("test_data/catalog_item_request_template.json")
	json.Unmarshal(data, resourceViewsTemplate)
	return resourceViewsTemplate
}
