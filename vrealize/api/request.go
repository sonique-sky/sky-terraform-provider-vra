package api

import "fmt"

func (c *RestClient) GetRequest(requestId string) (*Request, error) {
	path := fmt.Sprintf(fmtRequest, requestId)
	requestStatusViewTemplate := new(Request)

	err := c.get(path, requestStatusViewTemplate, noCheck)

	return requestStatusViewTemplate, err
}

func (c *RestClient) GetRequestResource(requestId string, resourceType string) (*ResourceWrapper, error) {
	path := fmt.Sprintf("/catalog-service/api/consumer/requests/%s/resources?$filter=resourceType/id+eq+'%s'", requestId, resourceType)
	resourceViewsTemplate := new(ResourceWrapper)

	err := c.get(path, resourceViewsTemplate, noCheck)

	if err != nil {
		return nil, err
	}

	return resourceViewsTemplate, nil
}
