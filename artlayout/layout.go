package main

import (
	"fmt"
	"sort"
)

func ShowRepos(repos []ArtifactoryRepoDetailsResponse) {
	layouts := make(map[string]int)

	for _, repo := range repos {
		layout := repo.RepoLayoutRef
		layouts[layout]++
	}

	type layoutCount struct {
		layout string
		count  int
	}
	var layoutCounts []layoutCount
	for layout, count := range layouts {
		layoutCounts = append(layoutCounts, layoutCount{layout, count})
	}

	sort.Slice(layoutCounts, func(i, j int) bool {
		if layoutCounts[i].count != layoutCounts[j].count {
			return layoutCounts[i].count > layoutCounts[j].count
		}
		return layoutCounts[i].layout < layoutCounts[j].layout
	})

	for _, lc := range layoutCounts {
		fmt.Printf("%s: %d\n", lc.layout, lc.count)
	}
}
