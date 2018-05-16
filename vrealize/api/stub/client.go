package stub

import (
	"github.com/sonique-sky/sky-terraform-provider-vra/vrealize/api"
	"io/ioutil"
	"fmt"
	"encoding/json"
)

type Client struct {
	GetRequestStub        func(string) (*api.Request, error)
	ReadCatalogByIdStub   func(string) (*api.RequestTemplate, error)
	ReadCatalogByNameStub func(string) (*api.RequestTemplate, error)
	RequestMachineStub    func(*api.RequestTemplate) (*api.RequestMachineResponse, error)
}

func (c *Client) GetRequest(requestId string) (*api.Request, error) {
	if c.GetRequestStub != nil {
		return c.GetRequestStub(requestId)
	}
	return nil, nil
}

func (c *Client) GetRequestResource(requestId string, resourceType string) (*api.ResourceWrapper, error) {
	return nil, nil
}

func (c *Client) GetMachine(resourceId string) (*api.Resource, error) {
	machine := new(api.Resource)
	readTestData("resource", machine)
	return machine, nil
}

func (c *Client) ReadCatalogByID(catalogId string) (*api.RequestTemplate, error) {
	if c.ReadCatalogByIdStub != nil {
		return c.ReadCatalogByIdStub(catalogId)
	}
	return nil, nil
}

func (c *Client) ReadCatalogByName(catalogName string) (*api.RequestTemplate, error) {
	if c.ReadCatalogByNameStub != nil {
		return c.ReadCatalogByNameStub(catalogName)
	}
	return nil, nil
}

func (c *Client) RequestMachine(template *api.RequestTemplate) (*api.RequestMachineResponse, error) {
	if c.RequestMachineStub != nil {
		return c.RequestMachineStub(template)
	}

	return nil, nil
}

func (c *Client) DestroyMachine(resourceId string) (error) {
	return nil
}

func readTestData(filename string, obj interface{}) {
	bytes, readErr := ioutil.ReadFile(fmt.Sprintf("../test_data/%s.json", filename))
	if readErr != nil {
		panic(readErr)
	}
	jsonErr := json.Unmarshal(bytes, obj)
	if jsonErr != nil {
		panic(jsonErr)
	}
}
