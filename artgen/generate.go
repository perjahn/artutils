package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func Generete(client *http.Client, baseurl string, token string, dryRun bool) error {
	count := 1100

	createRepositories(count, client, baseurl, token, dryRun)
	createUsers(count, client, baseurl, token, dryRun)
	createGroups(count, client, baseurl, token, dryRun)
	createPermissionTargets(count, client, baseurl, token, dryRun)

	return nil
}

func createRepositories(count int, client *http.Client, baseurl string, token string, dryRun bool) error {
	fmt.Print("Creating repos")

	for i := range count {
		fmt.Print(".")

		url := fmt.Sprintf("%s/artifactory/api/repositories/repo%04d", baseurl, i+1)

		repotypes := []string{"generic", "alpine", "ansible", "bower", "chef", "cocoapods", "composer", "conan", "conda", "cran", "debian", "docker", "gems", "gitlfs", "go", "gradle", "helm", "helmoci", "huggingfaceml", "ivy", "maven", "npm", "nuget", "oci", "opkg", "pub", "puppet", "pypi", "rpm", "sbt", "swift", "terraform", "terraformbackend", "vagrant", "yum"}
		layouttypes := []string{"simple-default", "ansible-default", "bower-default", "build-default", "cargo-default", "composer-default", "conan-default", "go-default", "ivy-default", "maven-1-default", "maven-2-default", "npm-default", "nuget-default", "puppet-default", "sbt-default", "sbt-ivy-default", "swift-default", "terraform-module-default", "terraform-provider-default", "vcs-default"}

		repo := ArtifactoryRepoRequest{
			Key:           fmt.Sprintf("repo%04d", i+1),
			Description:   fmt.Sprintf("description%04d", i+1),
			Rclass:        "local",
			PackageType:   repotypes[i%len(repotypes)],
			RepoLayoutRef: layouttypes[i%len(layouttypes)],
		}

		json, err := json.Marshal(repo)
		if err != nil {
			return fmt.Errorf("error creating repo, error generating json: %w", err)
		}
		req, err := http.NewRequest("PUT", url, bytes.NewReader(json))
		if err != nil {
			return fmt.Errorf("error creating repo, error creating request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		if !dryRun {
			resp, err := client.Do(req)
			if err != nil {
				return fmt.Errorf("error creating repo: %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				fmt.Printf("\nKey: '%s'\n", repo.Key)
				fmt.Printf("Url: '%s'\n", url)
				fmt.Printf("Unexpected status: '%s'\n", resp.Status)
				fmt.Printf("Request body: '%s'\n", req.Body)
				body, _ := io.ReadAll(resp.Body)
				fmt.Printf("Response body: '%s'\n", body)
				return fmt.Errorf("error creating repo")
			}
		}
	}

	fmt.Println()

	return nil
}

func createUsers(count int, client *http.Client, baseurl string, token string, dryRun bool) error {
	fmt.Print("Creating users")
	url := fmt.Sprintf("%s/access/api/v2/users", baseurl)

	for i := range count {
		fmt.Print(".")

		user := ArtifactoryUserRequest{
			Username:                 fmt.Sprintf("user%04d", i+1),
			Password:                 fmt.Sprintf("passWORD%04d", i+1),
			Email:                    fmt.Sprintf("user%04d@example.com", i+1),
			Groups:                   []string{},
			Admin:                    false,
			ProfileUpdatable:         true,
			InternalPasswordDisabled: false,
			DisableUiAccess:          false,
		}

		json, err := json.Marshal(user)
		if err != nil {
			return fmt.Errorf("error creating user, error generating json: %w", err)
		}
		req, err := http.NewRequest("POST", url, bytes.NewReader(json))
		if err != nil {
			return fmt.Errorf("error creating user, error creating request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		if !dryRun {
			resp, err := client.Do(req)
			if err != nil {
				return fmt.Errorf("error creating user: %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != 201 {
				fmt.Printf("\nUsername: '%s'\n", user.Username)
				fmt.Printf("Url: '%s'\n", url)
				fmt.Printf("Unexpected status: '%s'\n", resp.Status)
				fmt.Printf("Request body: '%s'\n", req.Body)
				body, _ := io.ReadAll(resp.Body)
				fmt.Printf("Response body: '%s'\n", body)
				return fmt.Errorf("error creating user")
			}
		}
	}

	fmt.Println()

	return nil
}

func createGroups(count int, client *http.Client, baseurl string, token string, dryRun bool) error {
	fmt.Print("Creating groups")
	url := fmt.Sprintf("%s/access/api/v2/groups", baseurl)

	for i := range count {
		fmt.Print(".")

		group := ArtifactoryGroupRequest{
			Name:            fmt.Sprintf("group%04d", i+1),
			Description:     fmt.Sprintf("description%04d", i+1),
			AutoJoin:        false,
			AdminPrivileges: false,
			ExternalId:      "",
			Members:         []string{fmt.Sprintf("user%04d", i+1)},
		}

		json, err := json.Marshal(group)
		if err != nil {
			return fmt.Errorf("error creating group, error generating json: %w", err)
		}
		req, err := http.NewRequest("POST", url, bytes.NewReader(json))
		if err != nil {
			return fmt.Errorf("error creating group, error creating request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		if !dryRun {
			resp, err := client.Do(req)
			if err != nil {
				return fmt.Errorf("error creating group: %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != 201 {
				fmt.Printf("\nGroupname: '%s'\n", group.Name)
				fmt.Printf("Url: '%s'\n", url)
				fmt.Printf("Unexpected status: '%s'\n", resp.Status)
				fmt.Printf("Request body: '%s'\n", req.Body)
				body, _ := io.ReadAll(resp.Body)
				fmt.Printf("Response body: '%s'\n", body)
				return fmt.Errorf("error creating group")
			}
		}
	}

	fmt.Println()

	return nil
}

func createPermissionTargets(count int, client *http.Client, baseurl string, token string, dryRun bool) error {
	fmt.Print("Creating permission targets")
	url := fmt.Sprintf("%s/access/api/v2/permissions", baseurl)

	for i := range count {
		fmt.Print(".")

		users := generateUsers(i)
		groups := generateGroups(i)

		key := fmt.Sprintf("repo%04d", i+1)

		targets := make(map[string]ArtifactoryPermissionDetailsTarget)
		targets[key] = ArtifactoryPermissionDetailsTarget{
			IncludePatterns: []string{"**"},
			ExcludePatterns: []string{},
		}

		artifactorypermissiontarget := ArtifactoryPermissionDetails{
			Name: key,
			Resources: ArtifactoryPermissionDetailsResources{
				Artifact: ArtifactoryPermissionDetailsArtifact{
					Actions: ArtifactoryPermissionDetailsActions{
						Users:  users,
						Groups: groups,
					},
					Targets: targets,
				},
			},
		}

		json, err := json.Marshal(artifactorypermissiontarget)

		//fmt.Printf("json: '%s'\n", string(json))

		if err != nil {
			return fmt.Errorf("error creating permission target, error generating json: %w", err)
		}
		req, err := http.NewRequest("POST", url, bytes.NewReader(json))
		if err != nil {
			return fmt.Errorf("error creating permission target, error creating request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		if !dryRun {
			resp, err := client.Do(req)
			if err != nil {
				return fmt.Errorf("error creating permission target: %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != 201 {
				fmt.Printf("\nKey: '%s'\n", key)
				fmt.Printf("Url: '%s'\n", url)
				fmt.Printf("Unexpected status: '%s'\n", resp.Status)
				fmt.Printf("Request body: '%s'\n", req.Body)
				body, _ := io.ReadAll(resp.Body)
				fmt.Printf("Response body: '%s'\n", body)
				return fmt.Errorf("error creating permission target")
			}
		}
	}

	fmt.Println()

	return nil
}

func generateUsers(i int) map[string][]string {
	permissions := []string{}

	if i%2 != 0 {
		permissions = append(permissions, "READ")
	}
	if i/2%2 != 0 {
		permissions = append(permissions, "ANNOTATE")
	}
	if i/4%2 != 0 {
		permissions = append(permissions, "WRITE")
	}
	if i/8%2 != 0 {
		permissions = append(permissions, "DELETE")
	}
	if i/16%2 != 0 {
		permissions = append(permissions, "MANAGE")
	}
	if i/32%2 != 0 {
		permissions = append(permissions, "SCAN")
	}

	users := make(map[string][]string)
	users[fmt.Sprintf("user%04d", i+1)] = permissions

	return users
}

func generateGroups(i int) map[string][]string {
	permissions := []string{}

	if i%2 != 0 {
		permissions = append(permissions, "READ")
	}
	if i/2%2 != 0 {
		permissions = append(permissions, "ANNOTATE")
	}
	if i/4%2 != 0 {
		permissions = append(permissions, "WRITE")
	}
	if i/8%2 != 0 {
		permissions = append(permissions, "DELETE")
	}
	if i/16%2 != 0 {
		permissions = append(permissions, "MANAGE")
	}
	if i/32%2 != 0 {
		permissions = append(permissions, "SCAN")
	}

	groups := make(map[string][]string)
	groups[fmt.Sprintf("group%04d", i+1)] = permissions

	return groups
}
