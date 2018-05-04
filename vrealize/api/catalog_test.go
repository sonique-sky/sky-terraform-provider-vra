package api

import (
	"testing"
	"gopkg.in/jarcoal/httpmock.v1"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
)

func TestAPIClient_GetCatalogItem(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	catalogId := "e5dd4fba-45ed-4943-b1fc-7f96239286be"

	httpmock.RegisterResponder("GET", fmt.Sprintf("http://localhost"+fmtCatalogItemsTemplate, catalogId),
		httpmock.NewStringResponder(200, testData("catalog_item.json")))

	template, err := client.GetCatalogItem(catalogId)

	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, catalogId, template.CatalogItemID, "Expected catalog ID not returned")
}

func TestAPIClient_GetCatalogItem_Fail(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	catalogId := "e5dd4fba-45ed-4943-b1fc-7f96239286be"

	httpmock.RegisterResponder("GET", fmt.Sprintf("http://localhost"+fmtCatalogItemsTemplate, catalogId),
		httpmock.NewErrorResponder(errors.New(testData("api_error.json"))))

	template, err := client.GetCatalogItem(catalogId)

	assert.NotNil(t,err,"Fail to generate exception")
	assert.Nil(t, template, "No template should be returned")
}

func TestAPIClient_Get(t *testing.T) {

	//id, e := client.ReadCatalogIDByName("The Name")

}
