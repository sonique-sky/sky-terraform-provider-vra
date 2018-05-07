package api

import (
	"testing"
	"gopkg.in/jarcoal/httpmock.v1"
	"fmt"
	"errors"
	"github.com/stretchr/testify/assert"
)

func TestAPIClient_GetResourceViews(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	requestId := "937099db-5174-4862-99a3-9c2666bfca28"

	httpmock.RegisterResponder("GET", fmt.Sprintf("http://localhost"+fmtRequestResourceViews, requestId),
		httpmock.NewStringResponder(200, testData("resource_views")))

	template, err := client.GetResourceViews(requestId)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(template.Content))
}

func TestAPIClient_GetResourceViews_Failure(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	requestId := "937099db-5174-4862-99a3-9c2666bfca28"

	httpmock.RegisterResponder("GET", fmt.Sprintf("http://localhost"+fmtRequestResourceViews, requestId),
		httpmock.NewErrorResponder(errors.New(testData("api_error"))))

	template, err := client.GetResourceViews(requestId)

	assert.NotNil(t, err)
	assert.Nil(t, template)
}

func TestAPIClient_GetRequestStatus(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	requestId := "937099db-5174-4862-99a3-9c2666bfca28"

	httpmock.RegisterResponder("GET", fmt.Sprintf("http://localhost"+fmtRequest, requestId),
		httpmock.NewStringResponder(200, testData("request_list")))

	template, err := client.GetRequestStatus(requestId)

	assert.Nil(t, err)
	assert.Equal(t, "PENDING_PRE_APPROVAL", template.Phase)
}
