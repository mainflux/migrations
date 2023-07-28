// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"net/http"
	"time"
)

type createThingsRes struct {
	Things []Thing `json:"things"`
}

type createChannelsRes struct {
	Channels []Channel `json:"channels"`
}

type pageRes struct {
	Total  uint64 `json:"total"`
	Offset uint64 `json:"offset"`
	Limit  uint64 `json:"limit"`
}

// ThingsPage contains list of things in a page with proper metadata.
type ThingsPage struct {
	Things []Thing `json:"things"`
	pageRes
}

// ChannelsPage contains list of channels in a page with proper metadata.
type ChannelsPage struct {
	Channels []Channel `json:"channels"`
	pageRes
}

type GroupsPage struct {
	Groups []Group `json:"groups"`
	pageRes
}

type UsersPage struct {
	Users []User `json:"users"`
	pageRes
}

type MembersPage struct {
	Members []User `json:"members"`
	pageRes
}

// MembershipsPage contains page related metadata as well as list of memberships that
// belong to this page.
type MembershipsPage struct {
	pageRes
	Memberships []Group `json:"memberships"`
}

// PolicyPage contains page related metadata as well as list
// of Policies that belong to the page.
type PolicyPage struct {
	PageMetadata
	Policies []Policy
}
type KeyRes struct {
	ID        string     `json:"id,omitempty"`
	Value     string     `json:"value,omitempty"`
	IssuedAt  time.Time  `json:"issued_at,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

func (res KeyRes) Code() int {
	return http.StatusCreated
}

func (res KeyRes) Headers() map[string]string {
	return map[string]string{}
}

func (res KeyRes) Empty() bool {
	return res.Value == ""
}

type revokeCertsRes struct {
	RevocationTime time.Time `json:"revocation_time"`
}

type identifyThingResp struct {
	ID string `json:"id,omitempty"`
}
