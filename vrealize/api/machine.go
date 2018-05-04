package api

import (
	"fmt"
	"encoding/json"
	"log"
)

func (c *Client) DestroyMachine(resourceViewTemplate *ResourceViewsTemplate) (error) {
	destroyMachineAction, errDestroyAction := c.getDestroyAction(resourceViewTemplate)
	if errDestroyAction != nil {
		if errDestroyAction.Error() == "resource is not created or not found" {
			return nil
		}
		return fmt.Errorf("destory Machine action template failed to load: %v", errDestroyAction)
	}

	_, errDestroyMachine := destroyMachineAction.Execute()

	if errDestroyMachine != nil {
		return fmt.Errorf("destory Machine machine operation failed: %v", errDestroyMachine)
	}
	return nil
}

func (c *Client) PowerOffMachine(powerOffTemplate *ActionTemplate, resourceViewTemplate *ResourceViewsTemplate) (*ActionResponseTemplate, error) {
	powerOffMachineActionURL, err := getActionURL(resourceViewTemplate, "POST: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.machine.PowerOff}")

	if err != nil {
		return nil, err
	}

	actionResponse := new(ActionResponseTemplate)

	err = c.post(powerOffMachineActionURL, powerOffTemplate, actionResponse, expectHttpStatus(201))

	if err != nil {
		return nil, err
	}
	return actionResponse, nil
}


func (c *Client) RequestMachine(template *CatalogItemTemplate) (*RequestMachineResponse, error) {
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

