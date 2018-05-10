package api

import (
	"time"
)

//ResourceViewsTemplate - is used to store information
//related to resource template information.
type ResourceViewsTemplate struct {
	Content []struct {
		ResourceID   string `json:"resourceId"`
		RequestState string `json:"requestState"`
		ResourceType string `json:"resourceType"`
		Links []struct {
			Href string `json:"href"`
			Rel  string `json:"rel"`
		} `json:"links"`
	} `json:"content"`
	Links []interface{} `json:"links"`
}

//RequestStatusView - used to store REST response of
//request triggered against any resource.
type RequestStatusView struct {
	RequestCompletion struct {
		RequestCompletionState string `json:"requestCompletionState"`
		CompletionDetails      string `json:"CompletionDetails"`
	} `json:"requestCompletion"`
	Phase string `json:"phase"`
}

//RequestMachineResponse - used to store response of request
//created against machine provision.
type RequestMachineResponse struct {
	ID           string      `json:"id"`
	IconID       string      `json:"iconId"`
	Version      int         `json:"version"`
	State        string      `json:"state"`
	Description  string      `json:"description"`
	Reasons      interface{} `json:"reasons"`
	RequestedFor string      `json:"requestedFor"`
	RequestedBy  string      `json:"requestedBy"`
	Organization struct {
		TenantRef      string `json:"tenantRef"`
		TenantLabel    string `json:"tenantLabel"`
		SubtenantRef   string `json:"subtenantRef"`
		SubtenantLabel string `json:"subtenantLabel"`
	} `json:"organization"`

	RequestorEntitlementID   string                 `json:"requestorEntitlementId"`
	PreApprovalID            string                 `json:"preApprovalId"`
	PostApprovalID           string                 `json:"postApprovalId"`
	DateCreated              time.Time              `json:"dateCreated"`
	LastUpdated              time.Time              `json:"lastUpdated"`
	DateSubmitted            time.Time              `json:"dateSubmitted"`
	DateApproved             time.Time              `json:"dateApproved"`
	DateCompleted            time.Time              `json:"dateCompleted"`
	Quote                    interface{}            `json:"quote"`
	RequestData              map[string]interface{} `json:"requestData"`
	RequestCompletion        string                 `json:"requestCompletion"`
	RetriesRemaining         int                    `json:"retriesRemaining"`
	RequestedItemName        string                 `json:"requestedItemName"`
	RequestedItemDescription string                 `json:"requestedItemDescription"`
	Components               string                 `json:"components"`
	StateName                string                 `json:"stateName"`

	CatalogItemProviderBinding struct {
		BindingID string `json:"bindingId"`
		ProviderRef struct {
			ID    string `json:"id"`
			Label string `json:"label"`
		} `json:"providerRef"`
	} `json:"catalogItemProviderBinding"`

	Phase           string `json:"phase"`
	ApprovalStatus  string `json:"approvalStatus"`
	ExecutionStatus string `json:"executionStatus"`
	WaitingStatus   string `json:"waitingStatus"`
	CatalogItemRef struct {
		ID    string `json:"id"`
		Label string `json:"label"`
	} `json:"catalogItemRef"`
}

type ResourceDataEntry struct {
	Key string `json:"key"`
	Value struct {
		Type  string      `json:"type"`
		Value interface{} `json:"value"`
	}
}

type Resource struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	RequestID   string `json:"requestId"`
	ResourceData struct {
		Entries []ResourceDataEntry `json:"entries"`
	} `json:"resourceData"`
}

func (r *Resource) StringValue(key string) (string, bool) {
	for _, val := range r.ResourceData.Entries {
		if val.Key == key && val.Value.Type == "string" {
			return val.Value.Value.(string), true
		}
	}
	return "", false
}

func (r *Resource) IntValue(key string) (int, bool) {
	for _, val := range r.ResourceData.Entries {
		if val.Key == key && val.Value.Type == "integer" {
			return val.Value.Value.(int), true
		}
	}
	return 0, false
}

func (r *Resource) BoolValue(key string) (bool, bool) {
	for _, val := range r.ResourceData.Entries {
		if val.Key == key && val.Value.Type == "boolean" {
			return val.Value.Value.(bool), true
		}
	}
	return false, false
}
