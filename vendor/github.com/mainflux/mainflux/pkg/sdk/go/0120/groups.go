// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type GroupMetadata map[string]interface{}

type Member struct {
	ID   string
	Type string
}

type PageMetadata struct {
	Total    uint64
	Offset   uint64
	Limit    uint64
	Size     uint64
	Level    uint64
	Name     string
	Type     string
	Metadata GroupMetadata
}

type GroupPage struct {
	PageMetadata
	Groups []Group
}

type MemberPage struct {
	PageMetadata
	Members []Member
}

const (
	groupsEndpoint = "groups"
	MaxLevel       = uint64(5)
	MinLevel       = uint64(1)
)

type assignRequest struct {
	Type    string   `json:"type,omitempty"`
	Members []string `json:"members"`
}

func (sdk mfSDK) CreateGroup(g Group, token string) (string, error) {
	data, err := json.Marshal(g)
	if err != nil {
		return "", err
	}

	url := createURL(sdk.baseURL, sdk.groupsPrefix, groupsEndpoint)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", Wrap(ErrFailedCreation, New(resp.Status))
	}

	id := strings.TrimPrefix(resp.Header.Get("Location"), fmt.Sprintf("/%s/", groupsEndpoint))
	return id, nil
}

func (sdk mfSDK) DeleteGroup(id, token string) error {
	endpoint := fmt.Sprintf("%s/%s", groupsEndpoint, id)

	url := createURL(sdk.baseURL, sdk.groupsPrefix, endpoint)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return Wrap(ErrFailedRemoval, New(resp.Status))
	}

	return nil
}

func (sdk mfSDK) Assign(memberIDs []string, memberType, groupID string, token string) error {
	var ids []string
	endpoint := fmt.Sprintf("%s/%s/members", groupsEndpoint, groupID)
	url := createURL(sdk.baseURL, sdk.groupsPrefix, endpoint)

	ids = append(ids, memberIDs...)
	assignReq := assignRequest{
		Type:    memberType,
		Members: ids,
	}

	data, err := json.Marshal(assignReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return Wrap(ErrMemberAdd, New(resp.Status))
	}

	return nil
}

func (sdk mfSDK) Unassign(token, groupID string, memberIDs ...string) error {
	var ids []string
	endpoint := fmt.Sprintf("%s/%s/members", groupsEndpoint, groupID)
	url := createURL(sdk.baseURL, sdk.groupsPrefix, endpoint)

	ids = append(ids, memberIDs...)
	assignReq := assignRequest{
		Members: ids,
	}

	data, err := json.Marshal(assignReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewReader(data))
	if err != nil {
		return err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return Wrap(ErrFailedRemoval, New(resp.Status))
	}

	return nil
}

func (sdk mfSDK) Members(groupID, token string, offset, limit uint64) (MemberPage, error) {
	endpoint := fmt.Sprintf("%s/%s/members?offset=%d&limit=%d&", groupsEndpoint, groupID, offset, limit)
	url := createURL(sdk.baseURL, sdk.groupsPrefix, endpoint)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return MemberPage{}, err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return MemberPage{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return MemberPage{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return MemberPage{}, Wrap(ErrFailedFetch, New(resp.Status))
	}

	var tp MemberPage
	if err := json.Unmarshal(body, &tp); err != nil {
		return MemberPage{}, err
	}

	return tp, nil
}

func (sdk mfSDK) Groups(offset, limit uint64, token string) (GroupPage, error) {
	endpoint := fmt.Sprintf("%s?offset=%d&limit=%d&tree=false", groupsEndpoint, offset, limit)
	url := createURL(sdk.baseURL, sdk.groupsPrefix, endpoint)
	return sdk.getGroups(token, url)
}

func (sdk mfSDK) Parents(id string, offset, limit uint64, token string) (GroupPage, error) {
	endpoint := fmt.Sprintf("%s/%s/parents?offset=%d&limit=%d&tree=false&level=%d", groupsEndpoint, id, offset, limit, MaxLevel)
	url := createURL(sdk.baseURL, sdk.groupsPrefix, endpoint)
	return sdk.getGroups(token, url)
}

func (sdk mfSDK) Children(id string, offset, limit uint64, token string) (GroupPage, error) {
	endpoint := fmt.Sprintf("%s/%s/children?offset=%d&limit=%d&tree=false&level=%d", groupsEndpoint, id, offset, limit, MaxLevel)
	url := createURL(sdk.baseURL, sdk.groupsPrefix, endpoint)
	return sdk.getGroups(token, url)
}

func (sdk mfSDK) getGroups(token, url string) (GroupPage, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return GroupPage{}, err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return GroupPage{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return GroupPage{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return GroupPage{}, Wrap(ErrFailedFetch, New(resp.Status))
	}

	var tp GroupPage
	if err := json.Unmarshal(body, &tp); err != nil {
		return GroupPage{}, err
	}
	return tp, nil
}

func (sdk mfSDK) Group(id, token string) (Group, error) {
	endpoint := fmt.Sprintf("%s/%s", groupsEndpoint, id)
	url := createURL(sdk.baseURL, sdk.groupsPrefix, endpoint)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Group{}, err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return Group{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Group{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return Group{}, Wrap(ErrFailedFetch, New(resp.Status))
	}

	var t Group
	if err := json.Unmarshal(body, &t); err != nil {
		return Group{}, err
	}

	return t, nil
}

func (sdk mfSDK) UpdateGroup(t Group, token string) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("%s/%s", groupsEndpoint, t.ID)
	url := createURL(sdk.baseURL, sdk.groupsPrefix, endpoint)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
	if err != nil {
		return err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return Wrap(ErrFailedUpdate, New(resp.Status))
	}

	return nil
}
