package api

import "fmt"

func (c *Client) GetRequestStatus(ResourceID string) (*RequestStatusView, error) {
	path := fmt.Sprintf(fmtRequest, ResourceID)
	requestStatusViewTemplate := new(RequestStatusView)

	err := c.get(path, requestStatusViewTemplate, noCheck)

	return requestStatusViewTemplate, err
}

func (c *Client) GetResourceViews(ResourceID string) (*ResourceViewsTemplate, error) {
	path := fmt.Sprintf(fmtRequestResourceViews, ResourceID)
	resourceViewsTemplate := new(ResourceViewsTemplate)

	err := c.get(path, resourceViewsTemplate, noCheck)

	if err != nil {
		return nil, err
	}

	return resourceViewsTemplate, nil
}
