package api

import (
	"fmt"
	"encoding/json"
	"log"
)

func (c *RestClient) RequestMachine(template *CatalogItemTemplate) (*RequestMachineResponse, error) {
	path := fmt.Sprintf("/catalog-service/api/consumer/entitledCatalogItems/%s/requests", template.CatalogItemID)

	response := new(RequestMachineResponse)

	jsonBody, jErr := json.Marshal(template)
	if jErr != nil {
		log.Printf("Error marshalling template as JSON")
		return nil, jErr
	} else {
		log.Printf("JSON Request Info: %s", jsonBody)
	}

	err := c.post(path, template, response, noCheck)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *RestClient) GetMachine(resourceId string) (interface{}, error) {
	path := "/catalog-service/api/consumer/resources/" + resourceId

	c.get(path, resourceTemplate, noCheck)
}

func (c *RestClient) DestroyMachine(resourceViewTemplate *ResourceViewsTemplate) (error) {
	action, err := c.getDestroyAction(resourceViewTemplate)
	if err != nil {
		if err.Error() == "resource is not created or not found" {
			return nil
		}
		return fmt.Errorf("destory Machine action template failed to load: %v", err)
	}

	_, errDestroyMachine := action.Execute()

	if errDestroyMachine != nil {
		return fmt.Errorf("destory machine operation failed: %v", errDestroyMachine)
	}
	return nil
}
