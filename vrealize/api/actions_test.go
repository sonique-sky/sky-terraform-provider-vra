package api

import (
	"testing"
	"gopkg.in/jarcoal/httpmock.v1"
	"errors"
	"encoding/json"
	"io/ioutil"
	"fmt"
)

func TestAPIClient_GetDestroyActionTemplate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	resourceId := "b313acd6-0738-439c-b601-e3ebf9ebb49b"
	actionId := "3da0ca14-e7e2-4d7b-89cb-c6db57440d72"
	path := fmt.Sprintf("http://localhost"+fmtActionTemplate, resourceId, actionId)
	httpmock.RegisterResponder("GET", path,
		httpmock.NewStringResponder(200, testData("destroy_template.json")))

	resourceViewsTemplate := new(ResourceViewsTemplate)
	data, _ := ioutil.ReadFile("test_data/resource_views.json")
	json.Unmarshal(data, resourceViewsTemplate)

	_, err := client.getDestroyAction(resourceViewsTemplate)

	if err != nil {
		t.Errorf("Fail to get destroy action template %v", err)
	}
}

func TestAPIClient_GetDestroyActionTemplate_Failure(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	resourceViewsTemplate := new(ResourceViewsTemplate)

	resourceId := "b313acd6-0738-439c-b601-e3ebf9ebb49b"
	actionId := "3da0ca14-e7e2-4d7b-89cb-c6db57440d72"

	data, _ := ioutil.ReadFile("test_data/resource_views.json")
	json.Unmarshal(data, resourceViewsTemplate)
	httpmock.RegisterResponder("GET", fmt.Sprintf("http://localhost"+fmtActionTemplate, resourceId, actionId),
		httpmock.NewErrorResponder(errors.New(`{"errors":[{"code":50505,"source":null,"message":"System exception.","systemMessage":null,"moreInfoUrl":null}]}`)))

	_, err := client.getDestroyAction(resourceViewsTemplate)

	if err == nil {
		t.Errorf("Fail to get destroy action template exception.")
	}
}
