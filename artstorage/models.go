package main

type ArtifactoryStorageInfo struct {
	RepoStorageInfo []ArtifactoryRepoStorageInfo `json:"repositoriesSummaryList"`
}

type ArtifactoryRepoStorageInfo struct {
	RepoKey          string `json:"repoKey"`
	RepoType         string `json:"repoType"`
	FoldersCount     int    `json:"foldersCount"`
	FilesCount       int    `json:"filesCount"`
	UsedSpace        string `json:"usedSpace"`
	UsedSpaceInBytes int    `json:"usedSpaceInBytes"`
	ItemsCount       int    `json:"itemsCount"`
	PackageType      string `json:"packageType"`
	ProjectKey       string `json:"projectKey"`
	Percentage       string `json:"percentage"`
}
