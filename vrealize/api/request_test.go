package api

import (
	"testing"
	"gopkg.in/jarcoal/httpmock.v1"
	"fmt"
	"errors"
)

func TestAPIClient_GetResourceViews(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	requestId := "937099db-5174-4862-99a3-9c2666bfca28"

	httpmock.RegisterResponder("GET", fmt.Sprintf("http://localhost/catalog-service/api/consumer/requests/%s/resourceViews", requestId),
		httpmock.NewStringResponder(200, testData("resource_views.json")))

	template, err := client.GetResourceViews(requestId)

	if err != nil {
		t.Errorf("Fail to get resource views %v.", err)
	}
	if len(template.Content) == 0 {
		t.Errorf("No resources provisioned.")
	}

	httpmock.RegisterResponder("GET", fmt.Sprintf("http://localhost/catalog-service/api/consumer/requests/%s/resourceViews", requestId),
		httpmock.NewErrorResponder(errors.New(testData("api_error.json"))))

	template, err = client.GetResourceViews(requestId)
	if err == nil {
		t.Errorf("Succeed to get resource views %v.", err)
	}
	if template != nil {
		t.Errorf("Resources provisioned.")
	}
}
