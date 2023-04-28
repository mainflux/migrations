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

const thingsEndpoint = "things"
const connectEndpoint = "connect"

func (sdk mfSDK) CreateThing(t Thing, token string) (string, error) {
	data, err := json.Marshal(t)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/%s", sdk.thingsURL, thingsEndpoint)

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

	id := strings.TrimPrefix(resp.Header.Get("Location"), fmt.Sprintf("/%s/", thingsEndpoint))
	return id, nil
}

func (sdk mfSDK) CreateThings(things []Thing, token string) ([]Thing, error) {
	data, err := json.Marshal(things)
	if err != nil {
		return []Thing{}, err
	}

	url := fmt.Sprintf("%s/%s/%s", sdk.thingsURL, thingsEndpoint, "bulk")

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return []Thing{}, err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return []Thing{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return []Thing{}, Wrap(ErrFailedCreation, New(resp.Status))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []Thing{}, err
	}

	var ctr createThingsRes
	if err := json.Unmarshal(body, &ctr); err != nil {
		return []Thing{}, err
	}

	return ctr.Things, nil
}

func (sdk mfSDK) Things(token string, offset, limit uint64, name string) (ThingsPage, error) {
	endpoint := fmt.Sprintf("%s?offset=%d&limit=%d&name=%s", thingsEndpoint, offset, limit, name)
	url := fmt.Sprintf("%s/%s", sdk.thingsURL, endpoint)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return ThingsPage{}, err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return ThingsPage{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ThingsPage{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return ThingsPage{}, Wrap(ErrFailedFetch, New(resp.Status))
	}

	var tp ThingsPage
	if err := json.Unmarshal(body, &tp); err != nil {
		return ThingsPage{}, err
	}

	return tp, nil
}

func (sdk mfSDK) ThingsByChannel(token, chanID string, offset, limit uint64, disconn bool) (ThingsPage, error) {
	url := fmt.Sprintf("%s/channels/%s/things?offset=%d&limit=%d&disconnected=%t", sdk.thingsURL, chanID, offset, limit, disconn)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return ThingsPage{}, err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return ThingsPage{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ThingsPage{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return ThingsPage{}, Wrap(ErrFailedFetch, New(resp.Status))
	}

	var tp ThingsPage
	if err := json.Unmarshal(body, &tp); err != nil {
		return ThingsPage{}, err
	}

	return tp, nil
}

func (sdk mfSDK) Thing(id, token string) (Thing, error) {
	url := fmt.Sprintf("%s/%s/%s", sdk.thingsURL, thingsEndpoint, id)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Thing{}, err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return Thing{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Thing{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return Thing{}, Wrap(ErrFailedFetch, New(resp.Status))
	}

	var t Thing
	if err := json.Unmarshal(body, &t); err != nil {
		return Thing{}, err
	}

	return t, nil
}

func (sdk mfSDK) UpdateThing(t Thing, token string) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/%s/%s", sdk.thingsURL, thingsEndpoint, t.ID)

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

func (sdk mfSDK) DeleteThing(id, token string) error {
	url := fmt.Sprintf("%s/%s/%s", sdk.thingsURL, thingsEndpoint, id)

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

func (sdk mfSDK) Connect(connIDs ConnectionIDs, token string) error {
	data, err := json.Marshal(connIDs)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/%s", sdk.thingsURL, connectEndpoint)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return Wrap(ErrFailedConnect, New(resp.Status))
	}

	return nil
}

func (sdk mfSDK) DisconnectThing(thingID, chanID, token string) error {
	url := fmt.Sprintf("%s/%s/%s/%s/%s", sdk.thingsURL, channelsEndpoint, chanID, thingsEndpoint, thingID)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	resp, err := sdk.sendRequest(req, token, string(CTJSON))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return Wrap(ErrFailedDisconnect, New(resp.Status))
	}

	return nil
}
