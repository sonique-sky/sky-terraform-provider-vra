package api

import (
	"fmt"
)

func (c *RestClient) RequestMachine(template *RequestTemplate) (*RequestMachineResponse, error) {
	path := fmt.Sprintf("/catalog-service/api/consumer/entitledCatalogItems/%s/requests", template.CatalogItemID)

	response := new(RequestMachineResponse)

	err := c.post(path, template, response, noCheck)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *RestClient) GetMachine(resourceId string) (*Resource, error) {
	var resource = new(Resource)
	path := "/catalog-service/api/consumer/resources/" + resourceId

	err := c.get(path, resource, noCheck)
	if err != nil {
		return nil, err
	}

	return resource, nil

}

func (c *RestClient) DestroyMachine(resourceId string) (error){
	action, err := c.getDestroyAction(resourceId)
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
