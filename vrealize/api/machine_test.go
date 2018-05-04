package api

import (
	"testing"
	"gopkg.in/jarcoal/httpmock.v1"
	"fmt"
	"github.com/stretchr/testify/assert"
	"errors"
)

func TestAPIClient_RequestMachine(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://localhost/catalog-service/"+
		"api/consumer/entitledCatalogItems/e5dd4fba-45ed-4943-b1fc-7f96239286be/requests/template",
		httpmock.NewStringResponder(200, testData("request_template.json")))

	httpmock.RegisterResponder("POST", "http://localhost/catalog-service/"+
		"api/consumer/entitledCatalogItems/e5dd4fba-45ed-4943-b1fc-7f96239286be/requests",
		httpmock.NewStringResponder(201, testData("request_list.json")))

	template, err := client.GetCatalogItem("e5dd4fba-45ed-4943-b1fc-7f96239286be")
	if err != nil {
		t.Errorf("Failed to get catalog item template %v.", err)
	}
	if len(template.CatalogItemID) == 0 {
		t.Errorf("Catalog Id is empty.")
	}

	requestMachine, errorRequestMachine := client.RequestMachine(template)

	if errorRequestMachine != nil {
		t.Errorf("Failed to request the machine %v.", errorRequestMachine)
	}

	if len(requestMachine.ID) == 0 {
		t.Errorf("Failed to request machine.")
	}

	httpmock.RegisterResponder("POST", "http://localhost/catalog-service/"+
		"api/consumer/entitledCatalogItems/e5dd4fba-45ed-4943-b1fc-7f96239286be/requests",
		httpmock.NewErrorResponder(errors.New(testData("api_error.json"))))

	requestMachine, errorRequestMachine = client.RequestMachine(template)

	if errorRequestMachine == nil {
		t.Errorf("Failed to generate exception.")
	}

	if requestMachine != nil {
		t.Errorf("Deploy machine request succeeded.")
	}
}

func TestAPIClient_DestroyMachine(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	resourceId := "b313acd6-0738-439c-b601-e3ebf9ebb49b"
	actionId := "3da0ca14-e7e2-4d7b-89cb-c6db57440d72"

	httpmock.RegisterResponder("GET", fmt.Sprintf("http://localhost"+fmtRequestResourceViews, resourceId),
		httpmock.NewStringResponder(200, testData("resource_views.json")))

	httpmock.RegisterResponder("GET", fmt.Sprintf("http://localhost"+fmtActionTemplate, resourceId, actionId),
		httpmock.NewStringResponder(200, testData("destroy_template.json")))

	httpmock.RegisterResponder("POST", fmt.Sprintf("http://localhost"+fmtActionRequest, resourceId, actionId),
		httpmock.NewStringResponder(201, "{}"))

	templateResources, errTemplate := client.GetResourceViews(resourceId)
	if errTemplate != nil {
		t.Errorf("Failed to get the template resources %v", errTemplate)
	}
	err := client.DestroyMachine(templateResources)

	assert.Nil(t, err)
}

func TestAPIClient_DestroyMachine_NotFound(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	resourceId := "b313acd6-0738-439c-b601-e3ebf9ebb49b"
	actionId := "3da0ca14-e7e2-4d7b-89cb-c6db57440d72"

	httpmock.RegisterResponder("GET", fmt.Sprintf("http://localhost"+fmtRequestResourceViews, resourceId),
		httpmock.NewStringResponder(200, testData("resource_views.json")))

	httpmock.RegisterResponder("GET", fmt.Sprintf("http://localhost"+fmtActionTemplate, resourceId, actionId),
		httpmock.NewStringResponder(200, testData("destroy_template.json")))

	httpmock.RegisterResponder("POST", fmt.Sprintf("http://localhost"+fmtActionRequest, resourceId, actionId),
		httpmock.NewStringResponder(201, "{}"))

	templateResources, errTemplate := client.GetResourceViews(resourceId)
	if errTemplate != nil {
		t.Errorf("Failed to get the template resources %v", errTemplate)
	}
	err := client.DestroyMachine(templateResources)

	assert.Nil(t, err)
}
