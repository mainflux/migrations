// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	// CTJSON represents JSON content type.
	CTJSON ContentType = "application/json"

	// CTJSONSenML represents JSON SenML content type.
	CTJSONSenML ContentType = "application/senml+json"

	// CTBinary represents binary content type.
	CTBinary ContentType = "application/octet-stream"

	// EnabledStatus represents enable status for a client
	EnabledStatus = "enabled"

	// DisabledStatus represents disabled status for a client
	DisabledStatus = "enabled"
)

// ContentType represents all possible content types.
type ContentType string

var _ SDK = (*mfSDK)(nil)

var (
	// ErrFailedCreation indicates that entity creation failed.
	ErrFailedCreation = New("failed to create entity in the db")

	// ErrFailedList indicates that entities list failed.
	ErrFailedList = New("failed to list entities")

	// ErrFailedUpdate indicates that entity update failed.
	ErrFailedUpdate = New("failed to update entity")

	// ErrFailedFetch indicates that fetching of entity data failed.
	ErrFailedFetch = New("failed to fetch entity")

	// ErrFailedRemoval indicates that entity removal failed.
	ErrFailedRemoval = New("failed to remove entity")

	// ErrFailedEnable indicates that client enable failed.
	ErrFailedEnable = New("failed to enable client")

	// ErrFailedDisable indicates that client disable failed.
	ErrFailedDisable = New("failed to disable client")
)

type PageMetadata struct {
	Total        uint64   `json:"total"`
	Offset       uint64   `json:"offset"`
	Limit        uint64   `json:"limit"`
	Level        uint64   `json:"level,omitempty"`
	Email        string   `json:"email,omitempty"`
	Name         string   `json:"name,omitempty"`
	Type         string   `json:"type,omitempty"`
	Disconnected bool     `json:"disconnected,omitempty"`
	Metadata     Metadata `json:"metadata,omitempty"`
	Status       string   `json:"status,omitempty"`
	Action       string   `json:"action,omitempty"`
	Subject      string   `json:"subject,omitempty"`
	Object       string   `json:"object,omitempty"`
	Tag          string   `json:"tag,omitempty"`
	Owner        string   `json:"owner,omitempty"`
	SharedBy     string   `json:"shared_by,omitempty"`
	Visibility   string   `json:"visibility,omitempty"`
	OwnerID      string   `json:"owner_id,omitempty"`
}

// Credentials represent client credentials: it contains
// "identity" which can be a username, email, generated name;
// and "secret" which can be a password or access token.
type Credentials struct {
	Identity string `json:"identity,omitempty"` // username or generated login ID
	Secret   string `json:"secret,omitempty"`   // password or token
}

// SDK contains Mainflux API.
type SDK interface {
	// CreateUser registers mainflux user.
	CreateUser(user User, token string) (User, SDKError)

	// User returns user object by id.
	User(id, token string) (User, SDKError)

	// Users returns list of users.
	Users(pm PageMetadata, token string) (UsersPage, SDKError)

	// Members retrieves everything that is assigned to a group identified by groupID.
	Members(groupID string, meta PageMetadata, token string) (MembersPage, SDKError)

	// UserProfile returns user logged in.
	UserProfile(token string) (User, SDKError)

	// UpdateUser updates existing user.
	UpdateUser(user User, token string) (User, SDKError)

	// UpdateUserTags updates the user's tags.
	UpdateUserTags(user User, token string) (User, SDKError)

	// UpdateUserIdentity updates the user's identity
	UpdateUserIdentity(user User, token string) (User, SDKError)

	// UpdateUserOwner updates the user's owner.
	UpdateUserOwner(user User, token string) (User, SDKError)

	// UpdatePassword updates user password.
	UpdatePassword(id, oldPass, newPass, token string) (User, SDKError)

	// EnableUser changes the status of the user to enabled.
	EnableUser(id, token string) (User, SDKError)

	// DisableUser changes the status of the user to disabled.
	DisableUser(id, token string) (User, SDKError)

	// CreateToken receives credentials and returns user token.
	CreateToken(user User) (Token, SDKError)

	// RefreshToken receives credentials and returns user token.
	RefreshToken(token string) (Token, SDKError)

	// CreateThing registers new thing and returns its id.
	CreateThing(thing Thing, token string) (Thing, SDKError)

	// CreateThings registers new things and returns their ids.
	CreateThings(things []Thing, token string) ([]Thing, SDKError)

	// Things returns page of things.
	Things(pm PageMetadata, token string) (ThingsPage, SDKError)

	// ThingsByChannel returns page of things that are connected or not connected
	// to specified channel.
	ThingsByChannel(chanID string, pm PageMetadata, token string) (ThingsPage, SDKError)

	// Thing returns thing object by id.
	Thing(id, token string) (Thing, SDKError)

	// UpdateThing updates existing thing.
	UpdateThing(thing Thing, token string) (Thing, SDKError)

	// UpdateThingTags updates the client's tags.
	UpdateThingTags(thing Thing, token string) (Thing, SDKError)

	// UpdateThingSecret updates the client's secret
	UpdateThingSecret(id, secret, token string) (Thing, SDKError)

	// UpdateThingOwner updates the client's owner.
	UpdateThingOwner(thing Thing, token string) (Thing, SDKError)

	// EnableThing changes client status to enabled.
	EnableThing(id, token string) (Thing, SDKError)

	// DisableThing changes client status to disabled - soft delete.
	DisableThing(id, token string) (Thing, SDKError)

	// IdentifyThing validates thing's key and returns its ID
	IdentifyThing(key string) (string, SDKError)

	// CreateGroup creates new group and returns its id.
	CreateGroup(group Group, token string) (Group, SDKError)

	// Memberships
	Memberships(clientID string, pm PageMetadata, token string) (MembershipsPage, SDKError)

	// Groups returns page of groups.
	Groups(pm PageMetadata, token string) (GroupsPage, SDKError)

	// Parents returns page of users groups.
	Parents(id string, pm PageMetadata, token string) (GroupsPage, SDKError)

	// Children returns page of users groups.
	Children(id string, pm PageMetadata, token string) (GroupsPage, SDKError)

	// Group returns users group object by id.
	Group(id, token string) (Group, SDKError)

	// UpdateGroup updates existing group.
	UpdateGroup(group Group, token string) (Group, SDKError)

	// EnableGroup changes group status to enabled.
	EnableGroup(id, token string) (Group, SDKError)

	// DisableGroup changes group status to disabled - soft delete.
	DisableGroup(id, token string) (Group, SDKError)

	// CreateChannel creates new channel and returns its id.
	CreateChannel(channel Channel, token string) (Channel, SDKError)

	// CreateChannels registers new channels and returns their ids.
	CreateChannels(channels []Channel, token string) ([]Channel, SDKError)

	// Channels returns page of channels.
	Channels(pm PageMetadata, token string) (ChannelsPage, SDKError)

	// ChannelsByThing returns page of channels that are connected or not connected
	// to specified thing.
	ChannelsByThing(thingID string, pm PageMetadata, token string) (ChannelsPage, SDKError)

	// Channel returns channel data by id.
	Channel(id, token string) (Channel, SDKError)

	// UpdateChannel updates existing channel.
	UpdateChannel(channel Channel, token string) (Channel, SDKError)

	// EnableChannel changes channel status to enabled.
	EnableChannel(id, token string) (Channel, SDKError)

	// DisableChannel changes channel status to disabled - soft delete.
	DisableChannel(id, token string) (Channel, SDKError)

	// AddPolicy creates a policy for the given subject, so that, after
	// AddPolicy, `subject` has a `relation` on `object`. Returns a non-nil
	// error in case of failures.
	AddPolicy(p Policy, token string) SDKError

	// UpdatePolicy updates policies based on the given policy structure.
	UpdatePolicy(p Policy, token string) SDKError

	// ListPolicies lists policies based on the given policy structure.
	ListPolicies(pm PageMetadata, token string) (PolicyPage, SDKError)

	// DeletePolicy removes a policy.
	DeletePolicy(p Policy, token string) SDKError

	// Assign assigns member of member type (thing or user) to a group.
	Assign(memberType []string, memberID, groupID, token string) SDKError

	// Unassign removes member from a group.
	Unassign(memberType []string, groupID string, memberID string, token string) SDKError

	// Connect bulk connects things to channels specified by id.
	Connect(conns ConnectionIDs, token string) SDKError

	// Disconnect
	Disconnect(connIDs ConnectionIDs, token string) SDKError

	// ConnectThing
	ConnectThing(thingID, chanID, token string) SDKError

	// DisconnectThing disconnect thing from specified channel by id.
	DisconnectThing(thingID, chanID, token string) SDKError

	// AddBootstrap add bootstrap configuration
	AddBootstrap(cfg BootstrapConfig, token string) (string, SDKError)

	// View returns Thing Config with given ID belonging to the user identified by the given token.
	ViewBootstrap(id, token string) (BootstrapConfig, SDKError)

	// Update updates editable fields of the provided Config.
	UpdateBootstrap(cfg BootstrapConfig, token string) SDKError

	// Update boostrap config certificates
	UpdateBootstrapCerts(id string, clientCert, clientKey, ca string, token string) SDKError

	// Remove removes Config with specified token that belongs to the user identified by the given token.
	RemoveBootstrap(id, token string) SDKError

	// Whitelist updates Thing state Config with given ID belonging to the user identified by the given token.
	Whitelist(cfg BootstrapConfig, token string) SDKError

	// IssueCert issues a certificate for a thing required for mtls.
	IssueCert(thingID, valid, token string) (Cert, SDKError)

	// ViewCert returns a certificate given certificate ID
	ViewCert(certID, token string) (Cert, SDKError)

	// RevokeCert revokes certificate for thing with thingID
	RevokeCert(thingID, token string) (time.Time, SDKError)
}

type mfSDK struct {
	bootstrapURL   string
	certsURL       string
	httpAdapterURL string
	readerURL      string
	thingsURL      string
	usersURL       string

	msgContentType ContentType
	client         *http.Client
}

// Config contains sdk configuration parameters.
type Config struct {
	BootstrapURL   string
	CertsURL       string
	HTTPAdapterURL string
	ReaderURL      string
	ThingsURL      string
	UsersURL       string

	MsgContentType  ContentType
	TLSVerification bool
}

// NewSDK returns new mainflux SDK instance.
func NewSDK(conf Config) SDK {
	return &mfSDK{
		bootstrapURL:   conf.BootstrapURL,
		certsURL:       conf.CertsURL,
		httpAdapterURL: conf.HTTPAdapterURL,
		readerURL:      conf.ReaderURL,
		thingsURL:      conf.ThingsURL,
		usersURL:       conf.UsersURL,

		msgContentType: conf.MsgContentType,
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: !conf.TLSVerification,
				},
			},
		},
	}
}

// processRequest creates and send a new HTTP request, and checks for errors in the HTTP response.
// It then returns the response headers, the response body, and the associated error(s) (if any).
func (sdk mfSDK) processRequest(method, url, token, contentType string, data []byte, expectedRespCodes ...int) (http.Header, []byte, SDKError) {
	req, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return make(http.Header), []byte{}, NewSDKError(err)
	}

	if token != "" {
		if !strings.Contains(token, "Thing ") {
			token = "Bearer " + token
		}
		req.Header.Set("Authorization", token)
	}
	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}

	resp, err := sdk.client.Do(req)
	if err != nil {
		return make(http.Header), []byte{}, NewSDKError(err)
	}
	defer resp.Body.Close()

	sdkerr := CheckError(resp, expectedRespCodes...)
	if sdkerr != nil {
		return make(http.Header), []byte{}, sdkerr
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return make(http.Header), []byte{}, NewSDKError(err)
	}

	return resp.Header, body, nil
}

func (sdk mfSDK) withQueryParams(baseURL, endpoint string, pm PageMetadata) (string, error) {
	q, err := pm.query()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s?%s", baseURL, endpoint, q), nil
}

func (pm PageMetadata) query() (string, error) {
	q := url.Values{}
	if pm.Offset != 0 {
		q.Add("offset", strconv.FormatUint(pm.Offset, 10))
	}
	if pm.Limit != 0 {
		q.Add("limit", strconv.FormatUint(pm.Limit, 10))
	}
	if pm.Total != 0 {
		q.Add("total", strconv.FormatUint(pm.Total, 10))
	}
	if pm.Level != 0 {
		q.Add("level", strconv.FormatUint(pm.Level, 10))
	}
	if pm.Email != "" {
		q.Add("email", pm.Email)
	}
	if pm.Name != "" {
		q.Add("name", pm.Name)
	}
	if pm.Type != "" {
		q.Add("type", pm.Type)
	}
	if pm.Visibility != "" {
		q.Add("visibility", pm.Visibility)
	}
	if pm.Status != "" {
		q.Add("status", pm.Status)
	}
	if pm.Metadata != nil {
		md, err := json.Marshal(pm.Metadata)
		if err != nil {
			return "", NewSDKError(err)
		}
		q.Add("metadata", string(md))
	}
	return q.Encode(), nil
}
