package api

import (
	"fmt"
	"log"
)

type RequestTemplate struct {
	Type            string                 `json:"type"`
	CatalogItemID   string                 `json:"catalogItemId"`
	RequestedFor    string                 `json:"requestedFor"`
	BusinessGroupID string                 `json:"businessGroupId"`
	Description     string                 `json:"description"`
	Reasons         string                 `json:"reasons"`
	Data            map[string]interface{} `json:"data"`
}

type CatalogItem struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type CatalogItemWrapper struct {
	Item CatalogItem `json:"catalogItem"`
}

type entitledCatalogItemViews struct {
	Links   interface{}          `json:"links"`
	Content []CatalogItemWrapper `json:"content"`
	Metadata struct {
		TotalElements int `json:"totalElements"`
	} `json:"metadata"`
}

func (c *RestClient) ReadCatalogByID(catalogId string) (*RequestTemplate, error) {
	catalog := new(RequestTemplate)
	path := fmt.Sprintf("/catalog-service/api/consumer/entitledCatalogItems/%s/requests/template", catalogId)

	log.Printf("Path : %s", path)
	err := c.get(path, catalog, noCheck)
	if err != nil {
		return nil, err
	}

	return catalog, nil
}

func (c *RestClient) ReadCatalogByName(catalogName string) (*RequestTemplate, error) {
	template := new(entitledCatalogItemViews)
	path := fmt.Sprintf("/catalog-service/api/consumer/entitledCatalogItems?$filter=name+eq+'%s'", catalogName)

	err := c.get(path, template, noCheck)
	if err != nil {
		return nil, err
	}

	if template.Metadata.TotalElements != 1 {
		return nil, fmt.Errorf("could not identify catalog item named '%s'", catalogName)
	}

	catalog := template.Content[0]

	return c.ReadCatalogByID(catalog.Item.ID)
}
