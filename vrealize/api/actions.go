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
	RelLink  string
	Url      string
	Template ActionTemplate
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

func (c *RestClient) getAction(resourceViewsTemplate *ResourceViewsTemplate, templateLink string, actionLink string) (*Action, error) {
	templateUrl, err := getActionURL(resourceViewsTemplate, templateLink)

	if err != nil {
		return nil, err
	}

	actionTemplate := new(ActionTemplate)
	err = c.get(templateUrl, actionTemplate, noCheck)
	if err != nil {
		return nil, err
	}

	actionUrl, err := getActionURL(resourceViewsTemplate, actionLink)

	return &Action{
		client:   c,
		Template: *actionTemplate,
		Url:      actionUrl,
	}, nil
}

func getActionURL(template *ResourceViewsTemplate, relLink string) (templateActionURL string, err error) {
	var actionURL string
	l := len(template.Content)
	for i := 0; i < l; i++ {
		lengthLinks := len(template.Content[i].Links)
		for j := 0; j < lengthLinks; j++ {
			if template.Content[i].Links[j].Rel == relLink {
				actionURL = template.Content[i].Links[j].Href
			}
		}
	}

	if len(actionURL) == 0 {
		return "", fmt.Errorf("resource is not created or not found")
	}

	return actionURL, nil
}

func (c *RestClient) getPowerOffAction(resourceViewsTemplate *ResourceViewsTemplate) (*Action, error) {
	relLink := "GET Template: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.machine.PowerOff}"
	actionLink := "POST: {com.vmware.csp.component.iaas.proxy.provider@resource.action.name.machine.PowerOff}"
	return c.getAction(resourceViewsTemplate, relLink, actionLink)
}

func (c *RestClient) getDestroyAction(resourceViewsTemplate *ResourceViewsTemplate) (*Action, error) {
	relLink := "GET Template: {com.vmware.csp.component.cafe.composition@resource.action.deployment.destroy.name}"
	actionLink := "POST: {com.vmware.csp.component.cafe.composition@resource.action.deployment.destroy.name}"
	return c.getAction(resourceViewsTemplate, relLink, actionLink)
}
