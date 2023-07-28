// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package sdk

// updateClientSecretReq is used to update the client secret
type updateClientSecretReq struct {
	OldSecret string `json:"old_secret,omitempty"`
	NewSecret string `json:"new_secret,omitempty"`
}

type updateThingSecretReq struct {
	Secret string `json:"key,omitempty"`
}

// updateClientIdentityReq is used to update the client identity
type updateClientIdentityReq struct {
	token    string
	id       string
	Identity string `json:"identity,omitempty"`
}

// UserPasswordReq contains old and new passwords
type UserPasswordReq struct {
	OldPassword string `json:"old_password,omitempty"`
	Password    string `json:"password,omitempty"`
}

// ConnectionIDs contains ID lists of things and channels to be connected
type ConnectionIDs struct {
	ChannelIDs []string `json:"group_ids"`
	ThingIDs   []string `json:"client_ids"`
	Actions    []string `json:"actions,omitempty"`
}

type tokenReq struct {
	Identity string `json:"identity"`
	Secret   string `json:"secret"`
}

type identifyThingReq struct {
	Token string `json:"token,omitempty"`
}
