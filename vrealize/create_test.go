package vrealize

import (
	"fmt"
	"testing"
	"io/ioutil"
	"encoding/json"

	"github.com/stretchr/testify/assert"

	"github.com/sonique-sky/sky-terraform-provider-vra/vrealize/api/stub"
	"github.com/sonique-sky/sky-terraform-provider-vra/vrealize/api"
)

func TestCreate_FailsWithNoCatalogItems(t *testing.T) {
	client := &stub.Client{
		ReadCatalogByNameStub: func(reqId string) (*api.RequestTemplate, error) {
			return nil, fmt.Errorf("fail message")
		},
	}

	req := &TerraformRequest{
	}

	_, err := Create(client, req)

	assert.NotNil(t, err)
	assert.Equal(t, "catalog lookup failed: fail message", err.Error())
}

func TestCreate_WithCatalogID(t *testing.T) {
	actualTemplate := new(api.RequestTemplate)

	stubTemplate := new(api.RequestTemplate)
	bytes, _ := ioutil.ReadFile("api/test_data/catalog_item_request_template.json")
	json.Unmarshal(bytes, stubTemplate)

	client := &stub.Client{
		ReadCatalogByIdStub: func(reqId string) (*api.RequestTemplate, error) {
			return stubTemplate, nil
		},
		RequestMachineStub: func(requestTemplate *api.RequestTemplate) (*api.RequestMachineResponse, error) {
			actualTemplate = requestTemplate
			return nil, nil
		},
	}

	req := &TerraformRequest{
		CatalogId: "cat_id",
		DeploymentConfiguration: map[string]interface{}{
			"description": "A Deployment Description",
			"reasons":     "A Deployment Reason",
		},
		ResourceConfiguration: map[string]interface{}{
			"CentOS_6.3.description": "A Machine Description",
			"CentOS_6.3.cpu":         "4",
			"CentOS_6.3.memory":      "16384",
			"CentOS_6.3.untemplated": "okay",
			"CentOS_6.3.foo":         "bar",
		},
	}

	_, err := Create(client, req)

	assert.Nil(t, err)
	assert.Equal(t, "A Deployment Description", actualTemplate.Description)
	assert.Equal(t, "A Deployment Reason", actualTemplate.Reasons)

	assert.Equal(t, "A Machine Description", deNest(actualTemplate.Data,"CentOS_6.3","data", "description"))
	assert.Equal(t, "4", deNest(actualTemplate.Data,"CentOS_6.3","data", "cpu"))
	assert.Equal(t, "16384", deNest(actualTemplate.Data,"CentOS_6.3","data", "memory"))
	assert.Equal(t, "okay", deNest(actualTemplate.Data,"CentOS_6.3", "data", "untemplated"))

	assert.Equal(t, "bar", deNest(actualTemplate.Data, "CentOS_6.3", "data", "mo_data", "foo"))

	assert.Equal(t, "okay", deNest(actualTemplate.Data, "CentOS_6.3", "data", "data", "untemplated"))
}

func deNest(source map[string]interface{}, keys ...string) string {
	var resp = source[keys[0]]
	for i := 1; i < len(keys)-1; i++ {
		resp = resp.(map[string]interface{})[keys[i]]
	}
	return resp.(map[string]interface{})[keys[len(keys)-1]].(string)

}
