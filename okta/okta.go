package okta

import "time"

type OktaGroups struct {
	Groups []OktaGroup
}

type OktaGroup struct {
	// TODO: Flush this out from their API docs!
	Id                    string
	Created               time.Time
	LastUpdated           time.Time
	LastMembershipUpdated time.Time
	ObjectClass           []string
	Type                  string
	Profile               OktaGroupProfile
	Links                 OktaGroupLinks `json:"_links"`
}

type OktaGroupProfile struct {
	Name        string
	Description string
}

type OktaGroupLinks struct {
	Logo  []OktaLink
	Users OktaLink
	Apps  OktaLink
}

type OktaLink struct {
	Name string
	Href string
	Type string
}
