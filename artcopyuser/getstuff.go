package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func GetStuff(
	client *http.Client,
	baseurl string,
	token string) ([]ArtifactoryPermissionDetails, error) {

	permissions, err := getPermissions(client, baseurl, token)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Permissions count: %d\n", len(permissions))

	permissiondetails, err := getPermissionDetails(client, baseurl, token, permissions)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Permission details count: %d\n", len(permissiondetails))

	return permissiondetails, nil
}

func getPermissions(client *http.Client, baseurl string, token string) ([]ArtifactoryPermission, error) {
	var cursor string
	var allpermissions []ArtifactoryPermission

	for {
		var permissions []ArtifactoryPermission
		var err error

		permissions, cursor, err = getPermissionsPage(client, baseurl, token, cursor)
		if err != nil {
			return nil, err
		}

		allpermissions = append(allpermissions, permissions...)

		if cursor == "" {
			json, err := json.MarshalIndent(allpermissions, "", "  ")
			if err != nil {
				return nil, fmt.Errorf("error generating json: %w", err)
			}

			err = os.WriteFile("allpermissions.json", []byte(json), 0600)
			if err != nil {
				return nil, fmt.Errorf("error saving permissions: %w", err)
			}

			return allpermissions, nil
		}
	}
}

func getPermissionsPage(client *http.Client, baseurl string, token string, cursor string) ([]ArtifactoryPermission, string, error) {
	var url string
	if cursor == "" {
		url = fmt.Sprintf("%s/access/api/v2/permissions", baseurl)
	} else {
		url = fmt.Sprintf("%s/access/api/v2/permissions?cursor=%s", baseurl, cursor)
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

	var permissions ArtifactoryPermissions
	err = json.Unmarshal(body, &permissions)
	if err != nil {
		return nil, "", fmt.Errorf("error parsing response body: %w", err)
	}

	return permissions.Permissions, permissions.Cursor, nil
}

func getPermissionDetails(client *http.Client, baseurl string, token string, permissions []ArtifactoryPermission) ([]ArtifactoryPermissionDetails, error) {
	var allpermissiondetails []ArtifactoryPermissionDetails

	fmt.Println("Getting Permission details...")

	for _, permission := range permissions {
		fmt.Print(".")

		url := fmt.Sprintf("%s/access/api/v2/permissions/%s", baseurl, url.PathEscape(permission.Name))
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
			fmt.Printf("Url2: '%s'\n", permission.Uri)
			fmt.Printf("Unexpected status: '%s'\n", resp.Status)
			fmt.Printf("Response body: '%s'\n", body)
		}

		filename := fmt.Sprintf("allpermissiondetails-%s.json", permission.Name)
		err = os.WriteFile(filename, []byte(body), 0600)
		if err != nil {
			return nil, fmt.Errorf("error saving response body: %w", err)
		}

		var permissiondetails ArtifactoryPermissionDetails

		err = json.Unmarshal(body, &permissiondetails)
		if err != nil {
			return nil, fmt.Errorf("error parsing response body: %w", err)
		}
		permissiondetails.JsonSource = string(body)

		allpermissiondetails = append(allpermissiondetails, permissiondetails)
	}

	fmt.Println()

	json, err := json.MarshalIndent(allpermissiondetails, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error generating json: %w", err)
	}

	err = os.WriteFile("allpermissiondetails.json", []byte(string(json)), 0600)
	if err != nil {
		return nil, fmt.Errorf("error saving permission details: %w", err)
	}

	return allpermissiondetails, nil
}
