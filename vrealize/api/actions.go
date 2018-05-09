package api

import (
	"fmt"
)

type ActionTemplate struct {
	ActionID string `json:"actionId"`
	Data struct {
		Description  interface{} `json:"description"`
		Reasons      interface{} `json:"reasons"`
		ForceDestroy bool        `json:"ForceDestroy"`
	} `json:"data"`
	Description interface{} `json:"description"`
	ResourceID  string      `json:"resourceId"`
	Type        string      `json:"type"`
}

type Action struct {
	client   *RestClient
	Url      string
	Template ActionTemplate
}

type ActionSpec struct {
	BindingId string `json:"bindingId"`
	ID        string `json:"id"`
}

type ActionList struct {
	Actions []ActionSpec `json:"content"`
}

type ActionResponseTemplate struct {
}

func (a *Action) Execute() (*ActionResponseTemplate, error) {
	actionResponse := new(ActionResponseTemplate)

	err := a.client.post(a.Url, a.Template, actionResponse, expectHttpStatus(201))

	if err != nil {
		return nil, err
	}

	return actionResponse, nil

}

func (c *RestClient) getDestroyAction(resourceId string) (*Action, error) {

	path := fmt.Sprintf("/catalog-service/api/consumer/resources/%s/actions", resourceId)
	actionList := new(ActionList)
	err := c.get(path, actionList, noCheck)
	if err != nil {
		return nil, err
	}

	if spec, found := getActionSpec(actionList, "Infrastructure.Virtual.Action.Destroy"); found {
		template := new(ActionTemplate)
		specErr := c.get(fmt.Sprintf("%s/%s/requests/template", path, spec.ID), template, noCheck)
		if specErr != nil {
			return nil, specErr
		}

		return &Action{
			client:   c,
			Url:      fmt.Sprintf("%s/%s/requests", path, spec.ID),
			Template: *template,
		}, nil

	}
	return nil, fmt.Errorf("action not found")
}

func getActionSpec(actionList *ActionList, bindingId string) (ActionSpec, bool) {
	for _, spec := range actionList.Actions {
		if spec.BindingId == bindingId {
			return spec, true
		}
	}
	return ActionSpec{}, false
}
