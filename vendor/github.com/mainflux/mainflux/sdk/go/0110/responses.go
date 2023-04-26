// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package sdk

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
	Channel     string   `json:"channel,omitempty"`
	Subtopic    string   `json:"subtopic,omitempty"`
	Publisher   string   `json:"publisher,omitempty"`
	Protocol    string   `json:"protocol,omitempty"`
	Name        string   `json:"name,omitempty"`
	Unit        string   `json:"unit,omitempty"`
	Time        float64  `json:"time,omitempty"`
	UpdateTime  float64  `json:"update_time,omitempty"`
	Value       *float64 `json:"value,omitempty"`
	StringValue *string  `json:"string_value,omitempty"`
	DataValue   *string  `json:"data_value,omitempty"`
	BoolValue   *bool    `json:"bool_value,omitempty"`
	Sum         *float64 `json:"sum,omitempty"`
}

// MessagesPage contains list of messages in a page with proper metadata.
type MessagesPage struct {
	Messages []Message `json:"messages,omitempty"`
	pageRes
}
