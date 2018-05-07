package machine

import (
	"testing"
	"github.com/sonique-sky/sky-terraform-provider-vra/vrealize/api"
)


type BaseMock struct {

}

type MockClient struct {
	*BaseMock
	mockDelete func() error
}


func (m *BaseMock) GetResourceViews(requestId string) (*api.ResourceViewsTemplate, error) {
	return nil, nil
}

func (m *BaseMock) GetRequestStatus(requestId string) (*api.RequestStatusView, error) {
	return nil, nil
}

func (m *MockClient) DestroyMachine(resourceViewTemplate *api.ResourceViewsTemplate) (error) {
	return nil
}


func TestDelete(t *testing.T) {

}
