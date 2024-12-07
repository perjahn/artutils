package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func GetLDAPGroups(client *http.Client, baseurl string, token string) ([]ArtifactoryLDAPGroup, error) {
	url := fmt.Sprintf("%s/access/api/v1/ldap/groups", baseurl)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != 200 {
		fmt.Printf("Url: '%s'\n", url)
		fmt.Printf("Unexpected status: '%s'\n", resp.Status)
		fmt.Printf("Response body: '%s'\n", body)
	}

	var groups []ArtifactoryLDAPGroup
	err = json.Unmarshal(body, &groups)
	if err != nil {
		return nil, fmt.Errorf("error parsing response body: %w", err)
	}

	json, err := json.MarshalIndent(groups, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error generating json: %w", err)
	}

	err = os.WriteFile("allldapgroups.json", []byte(json), 0600)
	if err != nil {
		return nil, fmt.Errorf("error saving groups: %w", err)
	}

	return groups, nil
}

func GetGroups(client *http.Client, baseurl string, token string) ([]ArtifactoryGroup, error) {
	var cursor string
	var allgroups []ArtifactoryGroup

	for {
		var groups []ArtifactoryGroup
		var err error

		groups, cursor, err = getGroupsPage(client, baseurl, token, cursor)
		if err != nil {
			return nil, err
		}

		allgroups = append(allgroups, groups...)

		if cursor == "" {
			json, err := json.MarshalIndent(allgroups, "", "  ")
			if err != nil {
				return nil, fmt.Errorf("error generating json: %w", err)
			}

			err = os.WriteFile("allgroups.json", []byte(json), 0600)
			if err != nil {
				return nil, fmt.Errorf("error saving groups: %w", err)
			}

			return allgroups, nil
		}
	}
}

func getGroupsPage(client *http.Client, baseurl string, token string, cursor string) ([]ArtifactoryGroup, string, error) {
	var url string
	if cursor == "" {
		url = fmt.Sprintf("%s/access/api/v2/groups", baseurl)
	} else {
		url = fmt.Sprintf("%s/access/api/v2/groups?cursor=%s", baseurl, cursor)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != 200 {
		fmt.Printf("Url: '%s'\n", url)
		fmt.Printf("Unexpected status: '%s'\n", resp.Status)
		fmt.Printf("Response body: '%s'\n", body)
	}

	var groups ArtifactoryGroups
	err = json.Unmarshal(body, &groups)
	if err != nil {
		return nil, "", fmt.Errorf("error parsing response body: %w", err)
	}

	return groups.Groups, groups.Cursor, nil
}
