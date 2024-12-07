package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func GetLDAPSettings(client *http.Client, baseurl string, token string) ([]ArtifactoryLDAPSettings, error) {
	url := fmt.Sprintf("%s/access/api/v1/ldap/settings", baseurl)

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

	var settings []ArtifactoryLDAPSettings
	err = json.Unmarshal(body, &settings)
	if err != nil {
		return nil, fmt.Errorf("error parsing response body: %w", err)
	}

	json, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error generating json: %w", err)
	}

	err = os.WriteFile("allldapsettings.json", []byte(json), 0600)
	if err != nil {
		return nil, fmt.Errorf("error saving groups: %w", err)
	}

	return settings, nil
}

func GetLDAPGroupSettings(client *http.Client, baseurl string, token string) ([]ArtifactoryLDAPGroupSettings, error) {
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

	var groups []ArtifactoryLDAPGroupSettings
	err = json.Unmarshal(body, &groups)
	if err != nil {
		return nil, fmt.Errorf("error parsing response body: %w", err)
	}

	json, err := json.MarshalIndent(groups, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error generating json: %w", err)
	}

	err = os.WriteFile("allldapgroupsettings.json", []byte(json), 0600)
	if err != nil {
		return nil, fmt.Errorf("error saving groups: %w", err)
	}

	return groups, nil
}
