package main

import (
	"fmt"
	"slices"
)

func ShowPatterns(repos []ArtifactoryRepoResponse, permissiondetails []ArtifactoryPermissionDetails) {
	var allpatterns []string
	totalincludecount := 0
	totalexcludecount := 0

	for _, pd := range permissiondetails {
		for target := range pd.Resources.Artifact.Targets {
			patterns := pd.Resources.Artifact.Targets[target]

			includePatternsString := ""
			includeCount := 0
			if len(patterns.IncludePatterns) > 0 {
				includeCount = len(patterns.IncludePatterns)
				includePatternsString = fmt.Sprintf("'%s'", patterns.IncludePatterns[0])
				for _, pattern := range patterns.IncludePatterns[1:] {
					includePatternsString += fmt.Sprintf(", '%s'", pattern)
				}
			} else {
				//fmt.Printf("Missing include_patterns in permission target: %s\n", target)
			}

			excludePatternsString := ""
			excludeCount := 0
			if len(patterns.ExcludePatterns) > 0 {
				excludeCount = len(patterns.ExcludePatterns)
				excludePatternsString = fmt.Sprintf("'%s'", patterns.ExcludePatterns[0])
				for _, pattern := range patterns.ExcludePatterns[1:] {
					excludePatternsString += fmt.Sprintf(", '%s'", pattern)
				}
			} else {
				//fmt.Printf("Missing exclude_patterns in permission target: '%s'\n", pd.Name)
			}

			s1 := ""
			if includePatternsString != "'**'" {
				s1 = fmt.Sprintf(" Include (%d): %s", includeCount, includePatternsString)
			}

			s2 := ""
			if excludePatternsString != "" && excludePatternsString != "''" {
				s2 = fmt.Sprintf(" Exclude (%d): %s", excludeCount, excludePatternsString)
			}

			if s1 != "" || s2 != "" {
				separator := ","
				if s1 == "" || s2 == "" {
					separator = ""
				}
				if s1 != "" {
					totalincludecount += includeCount
				}
				if s2 != "" {
					totalexcludecount += excludeCount
				}

				p := fmt.Sprintf("'%s' / '%s':%s%s%s", pd.Name, target, s1, separator, s2)
				allpatterns = append(allpatterns, p)
			}
		}
	}

	slices.Sort(allpatterns)

	for _, pattern := range allpatterns {
		fmt.Println(pattern)
	}

	fmt.Printf("Patterns: %d, Includes: %d, Excludes: %d\n", len(allpatterns), totalincludecount, totalexcludecount)
}
