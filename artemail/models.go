package main

import "time"

type ArtifactoryUsers struct {
	Users  []ArtifactoryUser `json:"users"`
	Cursor string            `json:"cursor"`
}

type ArtifactoryUser struct {
	Username string `json:"username"`
	Uri      string `json:"uri"`
	Realm    string `json:"realm"`
	Status   string `json:"status"`
}

type ArtifactoryUserDetails struct {
	Username                 string    `json:"username"`
	Email                    string    `json:"email"`
	Admin                    bool      `json:"admin"`
	EffectiveAdmin           bool      `json:"effective_admin"`
	ProfileUpdatable         bool      `json:"profile_updatable"`
	DisableUiAccess          bool      `json:"disable_ui_access"`
	InternalPasswordDisabled bool      `json:"internal_password_disabled"`
	LastLoggedIn             time.Time `json:"last_logged_in"`
	Realm                    string    `json:"realm"`
	Groups                   []string  `json:"groups"`
	Status                   string    `json:"status"`
}

type ArtifactoryUserDetailsEmail struct {
	Email string `json:"email"`
}
