package api

import (
	"fmt"
	"log"
)

type CatalogItemTemplate struct {
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

func (c *RestClient) getCatalogItem(uuid string) (*CatalogItemTemplate, error) {
	template := new(CatalogItemTemplate)
	path := fmt.Sprintf("/catalog-service/api/consumer/entitledCatalogItems/%s/requests/template", uuid)

	err := c.get(path, template, noCheck)
	if err != nil {
		return nil, err
	}

	return template, nil
}

func (c *RestClient) ReadCatalogByID(catalogId string) (*CatalogItemTemplate, error) {
	catalog := new(CatalogItemTemplate)
	path := fmt.Sprintf("/catalog-service/api/consumer/entitledCatalogItems/%s/requests/template", catalogId)

	log.Printf("Path : %s", path)
	err := c.get(path, catalog, noCheck)
	if err != nil {
		return nil, err
	}

	return catalog, nil
}

func (c *RestClient) ReadCatalogByName(catalogName string) (*CatalogItemTemplate, error) {
	template := new(entitledCatalogItemViews)
	path := fmt.Sprintf(fmtCatalogItemsSearch, catalogName)

	err := c.get(path, template, noCheck)
	if err != nil {
		return nil, err
	}

	if template.Metadata.TotalElements != 1 {
		return nil, fmt.Errorf("could not identify catalog item named '%s'", catalogName)
	}

	catalog := template.Content[0]

	return c.getCatalogItem(catalog.Item.ID)
}
