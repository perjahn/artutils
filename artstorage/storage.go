package main

import (
	"fmt"
	"sort"
)

func ShowStorage(repos []ArtifactoryRepoStorageInfo) {
	sort.Slice(repos, func(i, j int) bool {
		return repos[i].UsedSpaceInBytes < repos[j].UsedSpaceInBytes
	})

	for _, repo := range repos {
		fmt.Printf("%s: %d\n", repo.RepoKey, repo.UsedSpaceInBytes)
	}
}
