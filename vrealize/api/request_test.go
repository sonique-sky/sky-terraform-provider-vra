package api

import (
	"testing"
	"gopkg.in/jarcoal/httpmock.v1"
	"fmt"
	"errors"
	"github.com/stretchr/testify/assert"
	"net/url"
)

func TestRestClient_GetRequestResource(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	requestId := "937099db-5174-4862-99a3-9c2666bfca28"
	resourceType := "Infrastructure.Virtual"

	searchPath := fmt.Sprintf("http://localhost/catalog-service/api/consumer/requests/%s/resources?%%24filter=resourceType%%2Fid+eq+%s", requestId,
		url.QueryEscape(fmt.Sprintf("'%s'", resourceType)))

	httpmock.RegisterResponder("GET", searchPath,
		httpmock.NewStringResponder(200, testData("request_resource_search_results")))


	resource, err := client.GetRequestResource(requestId, resourceType)

	assert.Nil(t, err)
	assert.NotNil(t, resource)
	assert.Equal(t, 1, len(resource.Resources))

	res := resource.Resources[0]

	if val, found := res.StringValue("foo"); found {
		assert.Equal(t, "foo", val)
	}

}

func TestRestClient_GetRequestResource_Failure(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	requestId := "937099db-5174-4862-99a3-9c2666bfca28"
	resourceType := "Infrastructure.Virtual"

	searchPath := fmt.Sprintf("http://localhost/catalog-service/api/consumer/requests/%s/resources?%%24filter=resourceType%%2Fid+eq+%s", requestId,
		url.QueryEscape(fmt.Sprintf("'%s'", resourceType)))

	httpmock.RegisterResponder("GET", searchPath,
		httpmock.NewErrorResponder(errors.New(testData("api_error"))))

	template, err := client.GetRequestResource(requestId, "Infrastructure.Virtual")

	assert.NotNil(t, err)
	assert.Nil(t, template)
}

func TestAPIClient_GetRequestStatus(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	requestId := "937099db-5174-4862-99a3-9c2666bfca28"

	httpmock.RegisterResponder("GET", fmt.Sprintf("http://localhost"+fmtRequest, requestId),
		httpmock.NewStringResponder(200, testData("request_list")))

	template, err := client.GetRequest(requestId)

	assert.Nil(t, err)
	assert.Equal(t, "PENDING_PRE_APPROVAL", template.Phase)
}
