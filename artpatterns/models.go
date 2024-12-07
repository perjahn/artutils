package main

type ArtifactoryRepoResponse struct {
	Key         string `json:"key"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Url         string `json:"url"`
	PackageType string `json:"packageType"`
}

type ArtifactoryRepoRequest struct {
	Key         string `json:"key,omitempty"`
	Description string `json:"description,omitempty"`
	Rclass      string `json:"rclass"`
	PackageType string `json:"packageType,omitempty"`
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
