package api

import "fmt"

func (c *RestClient) GetRequestStatus(requestId string) (*RequestStatusView, error) {
	path := fmt.Sprintf(fmtRequest, requestId)
	requestStatusViewTemplate := new(RequestStatusView)

	err := c.get(path, requestStatusViewTemplate, noCheck)

	return requestStatusViewTemplate, err
}

func (c *RestClient) GetResourceViews(requestId string) (*ResourceViewsTemplate, error) {
	path := fmt.Sprintf(fmtRequestResourceViews, requestId)
	resourceViewsTemplate := new(ResourceViewsTemplate)

	err := c.get(path, resourceViewsTemplate, noCheck)

	if err != nil {
		return nil, err
	}

	return resourceViewsTemplate, nil
}
