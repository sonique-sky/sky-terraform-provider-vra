package api

import (
	"fmt"
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
	Links    interface{}          `json:"links"`
	Content  []CatalogItemWrapper `json:"content"`
	Metadata Metadata             `json:"metadata"`
}

type Metadata struct {
	TotalElements int `json:"totalElements"`
}

func (c *Client) GetCatalogItem(uuid string) (*CatalogItemTemplate, error) {
	template := new(CatalogItemTemplate)
	path := fmt.Sprintf("/catalog-service/api/consumer/entitledCatalogItems/%s/requests/template", uuid)

	err := c.get(path, template, noCheck)
	if err != nil {
		return nil, err
	}

	return template, nil
}

func (c *Client) ReadCatalogNameByID(catalogID string) (interface{}, error) {
	catalog := new(CatalogItemWrapper)
	path := fmt.Sprintf("/catalog-service/api/consumer/entitledCatalogItems/%s", catalogID)

	err := c.get(path, catalog, noCheck)
	if err != nil {
		return nil, err
	}

	return catalog.Item.Name, nil
}

func (c *Client) ReadCatalogIDByName(catalogName string) (interface{}, error) {
	template := new(entitledCatalogItemViews)
	path := fmt.Sprintf("catalog-service/api/consumer/entitledCatalogItemViews?$filter=name+eq+'%s'", catalogName)

	err := c.get(path, template, noCheck)
	if err != nil {
		return nil, err
	}

	if template.Metadata.TotalElements != 1 {
		return nil, fmt.Errorf("could not identify catalog item named '%s'", catalogName)
	}

	catalog := template.Content[0]

	return catalog.Item.ID, nil
}
