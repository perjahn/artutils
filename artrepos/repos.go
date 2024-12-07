package main

import (
	"fmt"
	"slices"
	"strings"
)

func ShowRepos(repos []ArtifactoryRepoResponse, permissiondetails []ArtifactoryPermissionDetails) {
	var allrepos []string

	for _, repo := range repos {
		found := false
		for _, pd := range permissiondetails {
			if repo.Key == pd.Name {
				found = true
				break
			}
		}
		if !found {
			var allpermissions []string

			for _, pd := range permissiondetails {
				_, ok := pd.Resources.Artifact.Targets[repo.Key]
				if ok {
					allrepos = append(allrepos, pd.Name)
				}
			}

			if len(allpermissions) == 0 {
				allrepos = append(allrepos, repo.Key+" <None>")
			} else {
				s := repo.Key + ": " + strings.Join(allpermissions, ", ")
				allrepos = append(allrepos, s)
			}
		}
	}

	slices.Sort(allrepos)

	for _, repo := range allrepos {
		fmt.Println(repo)
	}
}
