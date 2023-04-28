package sdk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	policiesEndpoint = "policies"
)

// Policy represents an argument struct for making a policy related function calls.
type Policy struct {
	OwnerID   string    `json:"owner_id"`
	Subject   string    `json:"subject"`
	Object    string    `json:"object"`
	Actions   []string  `json:"actions"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AddPolicy creates a policy for the given subject, so that, after
// AddPolicy, `subject` has a `relation` on `object`. Returns a non-nil
// error in case of failures.
func (sdk mfSDK) AddPolicy(p Policy, token string) SDKError {
	data, err := json.Marshal(p)
	if err != nil {
		return NewSDKError(err)
	}

	url := fmt.Sprintf("%s/%s", sdk.usersURL, policiesEndpoint)
	_, _, sdkerr := sdk.processRequest(http.MethodPost, url, token, string(CTJSON), data, http.StatusCreated)
	if sdkerr != nil {
		return sdkerr
	}

	return nil
}

// UpdatePolicy updates policies based on the given policy structure.
func (sdk mfSDK) UpdatePolicy(p Policy, token string) SDKError {
	data, err := json.Marshal(p)
	if err != nil {
		return NewSDKError(err)
	}

	url := fmt.Sprintf("%s/%s", sdk.usersURL, policiesEndpoint)

	_, _, sdkerr := sdk.processRequest(http.MethodPut, url, token, string(CTJSON), data, http.StatusNoContent)
	if sdkerr != nil {
		return sdkerr
	}

	return nil
}

// ListPolicies lists policies based on the given policy structure.
func (sdk mfSDK) ListPolicies(pm PageMetadata, token string) (PolicyPage, SDKError) {
	url, err := sdk.withQueryParams(sdk.usersURL, policiesEndpoint, pm)
	if err != nil {
		return PolicyPage{}, NewSDKError(err)
	}

	_, body, sdkerr := sdk.processRequest(http.MethodGet, url, token, string(CTJSON), nil, http.StatusOK)
	if sdkerr != nil {
		return PolicyPage{}, sdkerr
	}

	var pp PolicyPage
	if err := json.Unmarshal(body, &pp); err != nil {
		return PolicyPage{}, NewSDKError(err)
	}

	return pp, nil
}

// DeletePolicy removes a policy.
func (sdk mfSDK) DeletePolicy(p Policy, token string) SDKError {
	url := fmt.Sprintf("%s/%s/%s/%s", sdk.usersURL, policiesEndpoint, p.Object, p.Subject)

	data, err := json.Marshal(p)
	if err != nil {
		return NewSDKError(err)
	}

	_, _, sdkerr := sdk.processRequest(http.MethodDelete, url, token, string(CTJSON), data, http.StatusNoContent)
	return sdkerr
}

func (sdk mfSDK) Assign(memberType []string, memberID, groupID, token string) SDKError {
	var policy = Policy{
		Subject: memberID,
		Object:  groupID,
		Actions: memberType,
	}
	data, err := json.Marshal(policy)
	if err != nil {
		return NewSDKError(err)
	}
	url := fmt.Sprintf("%s/%s", sdk.usersURL, policiesEndpoint)
	_, _, sdkerr := sdk.processRequest(http.MethodPost, url, token, string(CTJSON), data, http.StatusCreated)
	return sdkerr
}

func (sdk mfSDK) Unassign(memberType []string, groupID string, memberID string, token string) SDKError {
	var policy = Policy{
		Subject: memberID,
		Object:  groupID,
		Actions: memberType,
	}
	data, err := json.Marshal(policy)
	if err != nil {
		return NewSDKError(err)
	}

	url := fmt.Sprintf("%s/%s/%s/%s", sdk.usersURL, policiesEndpoint, groupID, memberID)

	_, _, sdkerr := sdk.processRequest(http.MethodDelete, url, token, string(CTJSON), data, http.StatusNoContent)
	return sdkerr
}

func (sdk mfSDK) Connect(connIDs ConnectionIDs, token string) SDKError {
	data, err := json.Marshal(connIDs)
	if err != nil {
		return NewSDKError(err)
	}

	url := fmt.Sprintf("%s/%s", sdk.thingsURL, connectEndpoint)

	_, _, sdkerr := sdk.processRequest(http.MethodPost, url, token, string(CTJSON), data, http.StatusCreated)
	return sdkerr
}

func (sdk mfSDK) Disconnect(connIDs ConnectionIDs, token string) SDKError {
	data, err := json.Marshal(connIDs)
	if err != nil {
		return NewSDKError(err)
	}

	url := fmt.Sprintf("%s/%s", sdk.thingsURL, disconnectEndpoint)
	_, _, sdkerr := sdk.processRequest(http.MethodPost, url, token, string(CTJSON), data, http.StatusNoContent)
	return sdkerr
}

func (sdk mfSDK) ConnectThing(thingID, chanID, token string) SDKError {
	url := fmt.Sprintf("%s/%s/%s/%s/%s", sdk.thingsURL, channelsEndpoint, chanID, thingsEndpoint, thingID)

	_, _, err := sdk.processRequest(http.MethodPut, url, token, string(CTJSON), nil, http.StatusCreated)
	return err
}

func (sdk mfSDK) DisconnectThing(thingID, chanID, token string) SDKError {
	url := fmt.Sprintf("%s/%s/%s/%s/%s", sdk.thingsURL, channelsEndpoint, chanID, thingsEndpoint, thingID)

	_, _, err := sdk.processRequest(http.MethodDelete, url, token, string(CTJSON), nil, http.StatusNoContent)
	return err
}
