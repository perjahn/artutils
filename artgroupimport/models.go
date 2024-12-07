package main

type ArtifactoryLDAPGroupSettings struct {
	Name                 string `json:"name"`
	EnabledLdap          string `json:"enabled_ldap"`
	GroupBaseDn          string `json:"group_base_dn"`
	GroupNameAttribute   string `json:"group_name_attribute"`
	GroupMemberAttribute string `json:"group_member_attribute"`
	SubTree              bool   `json:"sub_tree"`
	ForceAttributeSearch bool   `json:"force_attribute_search"`
	Filter               string `json:"filter"`
	DescriptionAttribute string `json:"description_attribute"`
	Strategy             string `json:"strategy"`
}

type ArtifactoryGroupImport struct {
	ImportGroups      []ArtifactoryImportGroups    `json:"importGroups"`
	LdapGroupSettings ArtifactoryLDAPGroupSettings `json:"ldapGroupSettings"`
}

type ArtifactoryImportGroups struct {
	GroupName      string `json:"groupName"`
	Description    string `json:"description"`
	GroupDn        string `json:"groupDn"`
	RequiredUpdate string `json:"requiredUpdate"`
}

type ArtifactoryLDAPSettings struct {
	Key                      string                        `json:"key"`
	Enabled                  bool                          `json:"enabled"`
	LdapUrl                  string                        `json:"ldap_url"`
	Search                   ArtifactoryLDAPSettingsSearch `json:"search"`
	AutoCreateUser           bool                          `json:"auto_create_user"`
	EmailAttribute           string                        `json:"email_attribute"`
	LdapPositionProtection   bool                          `json:"ldap_position_protection"`
	AllowUserToAccessProfile bool                          `json:"allow_user_to_access_profile"`
	PagingSupportEnabled     bool                          `json:"paging_support_enabled"`
}

type ArtifactoryLDAPSettingsSearch struct {
	SearchFilter    string `json:"search_filter"`
	SearchBase      string `json:"search_base"`
	SearchSubtree   bool   `json:"search_subtree"`
	ManagerDn       string `json:"manager_dn"`
	ManagerPassword string `json:"manager_password"`
}

type ArtifactoryLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
