package api

import (
	"testing"
	"gopkg.in/jarcoal/httpmock.v1"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
)

func TestAPIClient_GetDestroyActionTemplate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	resourceId := "b313acd6-0738-439c-b601-e3ebf9ebb49b"
	actionId := "3da0ca14-e7e2-4d7b-89cb-c6db57440d72"

	httpmock.RegisterResponder("GET", fmt.Sprintf("http://localhost"+fmtActionTemplate, resourceId, actionId),
		httpmock.NewStringResponder(200, testData("destroy_template")))

	templateResources := resourceViews()

	_, err := client.getDestroyAction(templateResources)

	assert.Nil(t, err, "Should not error")
}

func TestAPIClient_GetDestroyActionTemplate_Failure(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	resourceId := "b313acd6-0738-439c-b601-e3ebf9ebb49b"
	actionId := "3da0ca14-e7e2-4d7b-89cb-c6db57440d72"

	templateResources := resourceViews()

	httpmock.RegisterResponder("GET", fmt.Sprintf("http://localhost"+fmtActionTemplate, resourceId, actionId),
		httpmock.NewErrorResponder(errors.New(`{"errors":[{"code":50505,"source":null,"message":"System exception.","systemMessage":null,"moreInfoUrl":null}]}`)))

	_, err := client.getDestroyAction(templateResources)

	assert.NotNil(t, err)
}
