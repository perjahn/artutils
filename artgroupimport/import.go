package main

import (
	"fmt"
	"io"
	"net/http"
)

func ImportGroups(
	groups []ArtifactoryLDAPGroup,
	client *http.Client,
	baseurl string,
	token string,
	dryRun bool) error {

	fmt.Printf("Groups to import: %d\n", len(groups))
	for _, group := range groups {
		err := importGroup(group, client, baseurl, token, dryRun)
		if err != nil {
			fmt.Printf("'%s': Warning: Ignoring group: %v\n", group.Name, err)
		}
	}

	return nil
}

func importGroup(
	group ArtifactoryLDAPGroup,
	client *http.Client,
	baseurl string,
	token string,
	dryRun bool) error {

	url := fmt.Sprintf("%s/access/api/v1/ldap/groups/%s/refresh", baseurl, group.Name)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("error importing group, error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	if !dryRun {
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("error importing group: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			fmt.Printf("Name: '%s'\n", group.Name)
			fmt.Printf("Url: '%s'\n", url)
			fmt.Printf("Unexpected status: '%s'\n", resp.Status)
			fmt.Printf("Request body: '%s'\n", req.Body)
			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("Response body: '%s'\n", body)
			return fmt.Errorf("error importing group")
		} else {
			fmt.Printf("'%s': Imported group successfully.\n", group.Name)
		}
	}

	return nil
}
