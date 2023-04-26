// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"net/http"
	"time"
)

type tokenRes struct {
	Token string `json:"token,omitempty"`
}

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

// Message represents a resolved (normalized) SenML record.
type Message struct {
	Channel     string   `json:"channel,omitempty" db:"channel" bson:"channel"`
	Subtopic    string   `json:"subtopic,omitempty" db:"subtopic" bson:"subtopic,omitempty"`
	Publisher   string   `json:"publisher,omitempty" db:"publisher" bson:"publisher"`
	Protocol    string   `json:"protocol,omitempty" db:"protocol" bson:"protocol"`
	Name        string   `json:"name,omitempty" db:"name" bson:"name,omitempty"`
	Unit        string   `json:"unit,omitempty" db:"unit" bson:"unit,omitempty"`
	Time        float64  `json:"time,omitempty" db:"time" bson:"time,omitempty"`
	UpdateTime  float64  `json:"update_time,omitempty" db:"update_time" bson:"update_time,omitempty"`
	Value       *float64 `json:"value,omitempty" db:"value" bson:"value,omitempty"`
	StringValue *string  `json:"string_value,omitempty" db:"string_value" bson:"string_value,omitempty"`
	DataValue   *string  `json:"data_value,omitempty" db:"data_value" bson:"data_value,omitempty"`
	BoolValue   *bool    `json:"bool_value,omitempty" db:"bool_value" bson:"bool_value,omitempty"`
	Sum         *float64 `json:"sum,omitempty" db:"sum" bson:"sum,omitempty"`
}

// MessagesPage contains list of messages in a page with proper metadata.
type MessagesPage struct {
	Messages []Message `json:"messages,omitempty"`
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
	Members []Member `json:"members"`
	pageRes
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

type retrieveKeyRes struct {
	ID        string     `json:"id,omitempty"`
	IssuerID  string     `json:"issuer_id,omitempty"`
	Subject   string     `json:"subject,omitempty"`
	IssuedAt  time.Time  `json:"issued_at,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

func (res retrieveKeyRes) Code() int {
	return http.StatusOK
}

func (res retrieveKeyRes) Headers() map[string]string {
	return map[string]string{}
}

func (res retrieveKeyRes) Empty() bool {
	return false
}
