package main

type ArtifactoryRepo struct {
	Key         string `json:"key"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Url         string `json:"url"`
	PackageType string `json:"packageType"`
}

type ArtifactoryRepoRequest struct {
	Key           string `json:"key,omitempty"`
	Description   string `json:"description,omitempty"`
	Rclass        string `json:"rclass"`
	PackageType   string `json:"packageType,omitempty"`
	RepoLayoutRef string `json:"repoLayoutRef,omitempty"`
}

type ArtifactoryPermissions struct {
	Permissions []ArtifactoryPermission `json:"permissions"`
	Cursor      string                  `json:"cursor"`
}

type ArtifactoryPermission struct {
	Name string `json:"name"`
	Uri  string `json:"uri"`
}

type ArtifactoryPermissionDetails struct {
	Name      string                                `json:"name"`
	Resources ArtifactoryPermissionDetailsResources `json:"resources"`
}

type ArtifactoryPermissionDetailsResources struct {
	Artifact ArtifactoryPermissionDetailsArtifact `json:"artifact"`
}

type ArtifactoryPermissionDetailsArtifact struct {
	Actions ArtifactoryPermissionDetailsActions           `json:"actions"`
	Targets map[string]ArtifactoryPermissionDetailsTarget `json:"targets"`
}

type ArtifactoryPermissionDetailsActions struct {
	Users  map[string][]string `json:"users"`
	Groups map[string][]string `json:"groups"`
}

type ArtifactoryPermissionDetailsTarget struct {
	IncludePatterns []string `json:"include_patterns"`
	ExcludePatterns []string `json:"exclude_patterns"`
}

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

type ArtifactoryGroups struct {
	Groups []ArtifactoryGroup `json:"groups"`
	Cursor string             `json:"cursor"`
}

type ArtifactoryGroup struct {
	GroupName string `json:"group_name"`
	Uri       string `json:"uri"`
}

type ArtifactoryUserRequest struct {
	Username                 string   `json:"username"`
	Password                 string   `json:"password"`
	Email                    string   `json:"email"`
	Groups                   []string `json:"groups"`
	Admin                    bool     `json:"admin"`
	ProfileUpdatable         bool     `json:"profile_updatable"`
	InternalPasswordDisabled bool     `json:"internal_password_disabled"`
	DisableUiAccess          bool     `json:"disable_ui_access"`
}

type ArtifactoryUserResponse struct {
	Username                 string   `json:"username"`
	Email                    string   `json:"email"`
	Groups                   []string `json:"groups"`
	Realm                    string   `json:"realm"`
	Status                   string   `json:"status"`
	Admin                    bool     `json:"admin"`
	ProfileUpdatable         bool     `json:"profile_updatable"`
	InternalPasswordDisabled bool     `json:"internal_password_disabled"`
	DisableUiAccess          bool     `json:"disable_ui_access"`
}

type ArtifactoryGroupRequest struct {
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	AutoJoin        bool     `json:"auto_join"`
	AdminPrivileges bool     `json:"admin_privileges"`
	ExternalId      string   `json:"external_id"`
	Members         []string `json:"members"`
}

type ArtifactoryGroupResponse struct {
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	AutoJoin        bool     `json:"auto_join"`
	AdminPrivileges bool     `json:"admin_privileges"`
	Realm           string   `json:"realm"`
	RealmAttributes string   `json:"realm_attributes"`
	ExternalId      string   `json:"external_id"`
	Members         []string `json:"members"`
}

type Repo struct {
	Name        string   `json:"name"`
	PackageType string   `json:"packageType,omitempty"`
	Description string   `json:"description,omitempty"`
	Rclass      string   `json:"rclass,omitempty"`
	Read        []string `json:"read,omitempty"`
	Annotate    []string `json:"annotate,omitempty"`
	Write       []string `json:"write,omitempty"`
	Delete      []string `json:"delete,omitempty"`
	Manage      []string `json:"manage,omitempty"`
}
