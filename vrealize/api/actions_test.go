package api

import (
	"testing"
	"gopkg.in/jarcoal/httpmock.v1"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
)

func TestAPIClient_GetDestroyAction(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	resourceId := "b313acd6-0738-439c-b601-e3ebf9ebb49b"
	actionId := "5883bbea-cd9a-4bf3-a2b3-f6cc20135435"

	actionsPath := fmt.Sprintf("http://localhost/catalog-service/api/consumer/resources/%s/actions", resourceId)
	httpmock.RegisterResponder("GET", actionsPath,
		httpmock.NewStringResponder(200, testData("action_list")))

	path := fmt.Sprintf("http://localhost/catalog-service/api/consumer/resources/%s/actions/%s/requests/template", resourceId, actionId)
	httpmock.RegisterResponder("GET", path,
		httpmock.NewStringResponder(200, testData("destroy_template")))

	action, err := client.getDestroyAction(resourceId)

	assert.Nil(t, err, "Should not error")
	assert.Equal(t, fmt.Sprintf("/catalog-service/api/consumer/resources/%s/actions/%s/requests", resourceId, actionId), action.Url)
}

func TestAPIClient_GetDestroyAction_Failure(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	resourceId := "b313acd6-0738-439c-b601-e3ebf9ebb49b"
	actionId := "5883bbea-cd9a-4bf3-a2b3-f6cc20135435"

	actionsPath := fmt.Sprintf("http://localhost/catalog-service/api/consumer/resources/%s/actions", resourceId)
	httpmock.RegisterResponder("GET", actionsPath,
		httpmock.NewStringResponder(200, testData("action_list")))

	path := fmt.Sprintf("http://localhost/catalog-service/api/consumer/resources/%s/actions/%s/requests/template", resourceId, actionId)
	httpmock.RegisterResponder("GET", path,
		httpmock.NewErrorResponder(errors.New(`{"errors":[{"code":50505,"source":null,"message":"System exception.","systemMessage":null,"moreInfoUrl":null}]}`)))

	_, err := client.getDestroyAction(resourceId)

	assert.NotNil(t, err)
	assert.Equal(t, "goo", err.Error())
}
