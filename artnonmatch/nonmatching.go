package main

import (
	"fmt"
	"slices"
)

func ShowNonmatching(repos []ArtifactoryRepoResponse, permissiondetails []ArtifactoryPermissionDetails) {
	var allnonmatching []string

	for _, pd := range permissiondetails {
		if len(pd.Resources.Artifact.Targets) != 1 {
			allnonmatching = append(allnonmatching, pd.Name)
			continue
		}

		var targetKey string
		for key := range pd.Resources.Artifact.Targets {
			targetKey = key
			break
		}

		var found bool
		for _, repo := range repos {
			if targetKey == repo.Key {
				found = true
				break
			}
		}
		if !found {
			allnonmatching = append(allnonmatching, pd.Name)
		}
	}

	slices.Sort(allnonmatching)

	for _, pattern := range allnonmatching {
		fmt.Println(pattern)
	}

	fmt.Printf("Repos: %d, Permissions: %d, Matching permissions: %d, Non-matching permissions: %d\n",
		len(repos), len(permissiondetails), len(permissiondetails)-len(allnonmatching), len(allnonmatching))
}
