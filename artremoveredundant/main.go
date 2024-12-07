package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
)

type ArtifactoryPermissionDetails struct {
	Name       string                                `json:"name"`
	Resources  ArtifactoryPermissionDetailsResources `json:"resources"`
	JsonSource string                                `json:"-"`
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
	IncludePatternsx []string `json:"include_patterns"`
	ExcludePatternsx []string `json:"exclude_patterns"`
}

type Repo struct {
	Permissions map[string]Permissions
}

type Permissions struct {
	Users  []string
	Groups []string
}

func main() {
	args := os.Args

	if len(args) != 2 {
		log.Fatalf("Usage: artremoveredundant <filename>")
	}

	filename := args[1]
	if filename == "" {
		log.Fatalf("Filename not provided")
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Fatalf("File not found: '%s'", filename)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	var permissionTargets []ArtifactoryPermissionDetails
	err = json.Unmarshal(data, &permissionTargets)
	if err != nil {
		log.Fatalf("Error parsing json: %v", err)
	}

	if len(permissionTargets) == 0 {
		log.Fatalf("No permission targets found in the file: '%s'", filename)
	}

	fmt.Printf("Loaded permission targets: %d\n", len(permissionTargets))

	allRepos := collectAllRepos(permissionTargets)
	fmt.Printf("Found unique repositories: %d\n", len(allRepos))

	effectivePermissions := calculateEffectivePermissions(permissionTargets, allRepos)
	/*
		for repo, perms := range effectivePermissions {
			fmt.Printf("Repo: %s\n", repo)
			for permType, principals := range perms.Permissions {
				fmt.Printf("  %s Users: %v, Groups: %v\n", permType, principals.Users, principals.Groups)
			}
		}
	*/

	ptNamesNotInRepos := findPTsNotMatchingRepos(permissionTargets, allRepos)
	fmt.Printf("Permission targets with names not matching any repo: %d\n", len(ptNamesNotInRepos))

	redundant := detectRedundantPermissionTargets(permissionTargets, allRepos, effectivePermissions, ptNamesNotInRepos)

	if len(redundant) == 0 {
		fmt.Println("No redundant permission targets detected.")
	} else {
		fmt.Printf("Detected redundant permission targets: %d\n", len(redundant))
		for _, name := range redundant {
			fmt.Printf("'%s'\n", name)
		}
	}
}

func collectAllRepos(targets []ArtifactoryPermissionDetails) []string {
	repoSet := make(map[string]bool)

	for _, pt := range targets {
		for targetName := range pt.Resources.Artifact.Targets {
			if targetName != "" {
				repoSet[targetName] = true
			}
		}
	}

	repos := make([]string, 0, len(repoSet))
	for repo := range repoSet {
		repos = append(repos, repo)
	}
	slices.Sort(repos)
	return repos
}

func calculateEffectivePermissions(targets []ArtifactoryPermissionDetails, repos []string) map[string]Repo {
	result := make(map[string]Repo)

	for _, reponame := range repos {
		repo := Repo{}

		for _, pt := range targets {
			_, exists := pt.Resources.Artifact.Targets[reponame]
			if exists {
				extracted := extractPermissionsFromActions(pt.Resources.Artifact.Actions)

				for permname, perms := range extracted.Permissions {
					if repo.Permissions == nil {
						repo.Permissions = make(map[string]Permissions)
					}
					if existingPerms, ok := repo.Permissions[permname]; ok {
						repo.Permissions[permname] = mergePermissions(existingPerms, perms)
					} else {
						repo.Permissions[permname] = perms
					}
				}
			}
		}

		result[reponame] = repo
	}

	return result
}

func calculateEffectivePermissionsExcluding(targets []ArtifactoryPermissionDetails, repos []string, excludeIndex int) map[string]Repo {
	result := make(map[string]Repo)

	for _, reponame := range repos {
		repo := Repo{}

		for i, pt := range targets {
			if i == excludeIndex {
				continue
			}

			_, exists := pt.Resources.Artifact.Targets[reponame]
			if exists {
				extracted := extractPermissionsFromActions(pt.Resources.Artifact.Actions)

				for permname, perms := range extracted.Permissions {
					if repo.Permissions == nil {
						repo.Permissions = make(map[string]Permissions)
					}
					if existingPerms, ok := repo.Permissions[permname]; ok {
						repo.Permissions[permname] = mergePermissions(existingPerms, perms)
					} else {
						repo.Permissions[permname] = perms
					}
				}
			}
		}

		result[reponame] = repo
	}

	return result
}

func extractPermissionsFromActions(actions ArtifactoryPermissionDetailsActions) Repo {
	repo := Repo{
		Permissions: make(map[string]Permissions),
	}

	for user, userperms := range actions.Users {
		for _, perm := range userperms {
			perms := repo.Permissions[perm]
			perms.Users = append(perms.Users, user)
			repo.Permissions[perm] = perms
		}
	}

	for group, groupPerms := range actions.Groups {
		for _, perm := range groupPerms {
			perms := repo.Permissions[perm]
			perms.Groups = append(perms.Groups, group)
			repo.Permissions[perm] = perms
		}
	}

	for permName, perms := range repo.Permissions {
		slices.Sort(perms.Users)
		slices.Sort(perms.Groups)
		repo.Permissions[permName] = perms
	}

	return repo
}

func mergePermissions(a, b Permissions) Permissions {
	result := Permissions{
		Users:  append(slices.Clone(a.Users), b.Users...),
		Groups: append(slices.Clone(a.Groups), b.Groups...),
	}

	userSet := make(map[string]bool)
	var uniqueUsers []string
	for _, u := range result.Users {
		if !userSet[u] {
			userSet[u] = true
			uniqueUsers = append(uniqueUsers, u)
		}
	}
	slices.Sort(uniqueUsers)
	result.Users = uniqueUsers

	groupSet := make(map[string]bool)
	var uniqueGroups []string
	for _, g := range result.Groups {
		if !groupSet[g] {
			groupSet[g] = true
			uniqueGroups = append(uniqueGroups, g)
		}
	}
	slices.Sort(uniqueGroups)
	result.Groups = uniqueGroups

	return result
}

func findPTsNotMatchingRepos(targets []ArtifactoryPermissionDetails, repos []string) []string {
	repoNameSet := make(map[string]bool)
	for _, repo := range repos {
		repoNameSet[repo] = true
	}

	var result []string
	for _, pt := range targets {
		if !repoNameSet[pt.Name] {
			result = append(result, pt.Name)
		}
	}

	slices.Sort(result)
	return result
}

func detectRedundantPermissionTargets(targets []ArtifactoryPermissionDetails, repos []string,
	effectivePerms map[string]Repo, ptToTest []string) []string {

	var redundant []string
	workingTargets := append([]ArtifactoryPermissionDetails{}, targets...)

	changed := true
	for changed {
		changed = false

		for _, ptName := range ptToTest {
			alreadyMarked := slices.Contains(redundant, ptName)
			if alreadyMarked {
				continue
			}

			var ptIndex int = -1
			for i, pt := range workingTargets {
				if pt.Name == ptName {
					ptIndex = i
					break
				}
			}
			if ptIndex == -1 {
				continue
			}

			permissionsWithout := calculateEffectivePermissionsExcluding(workingTargets, repos, ptIndex)

			isRedundant := true
			for _, repo := range repos {
				if !repoEqual(effectivePerms[repo], permissionsWithout[repo]) {
					isRedundant = false
					break
				}
			}

			if isRedundant {
				redundant = append(redundant, ptName)
				workingTargets = append(workingTargets[:ptIndex], workingTargets[ptIndex+1:]...)
				changed = true
				break
			}
		}
	}

	return redundant
}

func permissionsEqual(a, b Permissions) bool {
	if len(a.Users) != len(b.Users) || len(a.Groups) != len(b.Groups) {
		return false
	}
	for i := range a.Users {
		if a.Users[i] != b.Users[i] {
			return false
		}
	}
	for i := range a.Groups {
		if a.Groups[i] != b.Groups[i] {
			return false
		}
	}
	return true
}

func repoEqual(a, b Repo) bool {
	for permType, principalsA := range a.Permissions {
		principalsB, exists := b.Permissions[permType]
		if !exists || !permissionsEqual(principalsA, principalsB) {
			return false
		}
	}
	for permType, principalsB := range b.Permissions {
		principalsA, exists := a.Permissions[permType]
		if !exists || !permissionsEqual(principalsA, principalsB) {
			return false
		}
	}
	return true
}
