package api

import (
	"testing"
	"gopkg.in/jarcoal/httpmock.v1"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/url"
)

func TestAPIClient_GetCatalogItem(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	catalogId := "e5dd4fba-45ed-4943-b1fc-7f96239286be"

	path := fmt.Sprintf("http://localhost/catalog-service/api/consumer/entitledCatalogItems/%s/requests/template", catalogId)
	httpmock.RegisterResponder("GET", path,
		httpmock.NewStringResponder(200, testData("catalog_item_request_template")))

	template, err := client.getCatalogItem(catalogId)

	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, catalogId, template.CatalogItemID, "Expected catalog ID not returned")
}

func TestAPIClient_GetCatalogItem_Fail(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	catalogId := "e5dd4fba-45ed-4943-b1fc-7f96239286be"

	path := fmt.Sprintf("http://localhost/catalog-service/api/consumer/entitledCatalogItems/%s/requests/template", catalogId)
	httpmock.RegisterResponder("GET", path,
		httpmock.NewErrorResponder(errors.New(testData("api_error"))))

	template, err := client.getCatalogItem(catalogId)

	assert.NotNil(t, err, "Fail to generate exception")
	assert.Nil(t, template, "No template should be returned")
}

func TestAPIClient_ReadCatalogByName(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	catalogId := "e5dd4fba-45ed-4943-b1fc-7f96239286be"
	catalogName := "CentOS 7"

	searchPath := "http://localhost/catalog-service/api/consumer/entitledCatalogItems?%24filter=name+eq+" +
		url.QueryEscape(fmt.Sprintf("'%s'", catalogName))
	httpmock.RegisterResponder("GET", searchPath,
		httpmock.NewStringResponder(200, testData("catalog_item_search_results")))

	clientPath := fmt.Sprintf("http://localhost/catalog-service/api/consumer/entitledCatalogItems/%s/requests/template", catalogId)
	httpmock.RegisterResponder("GET", clientPath,
		httpmock.NewStringResponder(200, testData("catalog_item_request_template")))

	template, err := client.ReadCatalogByName(catalogName)

	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, template)

}

func TestAPIClient_ReadCatalogByID(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	catalogId := "e5dd4fba-45ed-4943-b1fc-7f96239286be"

	httpmock.RegisterResponder("GET", "http://localhost"+("/catalog-service/api/consumer/entitledCatalogItems")+"/"+catalogId,
		httpmock.NewStringResponder(200, testData("catalog_item")))

	httpmock.RegisterResponder("GET", fmt.Sprintf("http://localhost"+("/catalog-service/api/consumer/entitledCatalogItems")+"/%s/requests/template", catalogId),
		httpmock.NewStringResponder(200, testData("catalog_item_request_template")))

	template, err := client.ReadCatalogByID(catalogId)

	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, catalogId, template.CatalogItemID, "Expected catalog ID not returned")
}
