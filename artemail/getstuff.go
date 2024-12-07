package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func GetStuff(client *http.Client, baseurl string, token string) ([]ArtifactoryUserDetails, error) {
	users, err := getUsers(client, baseurl, token)
	if err != nil {
		return nil, err
	}

	fmt.Printf("User count: %d\n", len(users))

	userdetails, err := getUserDetails(client, baseurl, token, users)
	if err != nil {
		return nil, err
	}

	fmt.Printf("User details count: %d\n", len(userdetails))

	return userdetails, nil
}

func getUsers(client *http.Client, baseurl string, token string) ([]ArtifactoryUser, error) {
	var cursor string
	var allusers []ArtifactoryUser

	for {
		var users []ArtifactoryUser
		var err error

		users, cursor, err = getUsersPage(client, baseurl, token, cursor)
		if err != nil {
			return nil, err
		}

		allusers = append(allusers, users...)

		if cursor == "" {
			json, err := json.MarshalIndent(allusers, "", "  ")
			if err != nil {
				return nil, fmt.Errorf("error generating json: %w", err)
			}

			err = os.WriteFile("allusers.json", []byte(json), 0600)
			if err != nil {
				return nil, fmt.Errorf("error saving users: %w", err)
			}

			return allusers, nil
		}
	}
}

func getUsersPage(client *http.Client, baseurl string, token string, cursor string) ([]ArtifactoryUser, string, error) {
	var url string
	if cursor == "" {
		url = fmt.Sprintf("%s/access/api/v2/users", baseurl)
	} else {
		url = fmt.Sprintf("%s/access/api/v2/users?cursor=%s", baseurl, cursor)
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

	var users ArtifactoryUsers
	err = json.Unmarshal(body, &users)
	if err != nil {
		return nil, "", fmt.Errorf("error parsing response body: %w", err)
	}

	return users.Users, users.Cursor, nil
}

func getUserDetails(client *http.Client, baseurl string, token string, users []ArtifactoryUser) ([]ArtifactoryUserDetails, error) {
	var alluserdetails []ArtifactoryUserDetails

	fmt.Println("Getting user details...")

	for _, user := range users {
		fmt.Print(".")

		url := fmt.Sprintf("%s/access/api/v2/users/%s", baseurl, url.PathEscape(user.Username))
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
			fmt.Printf("Url2: '%s'\n", user.Uri)
			fmt.Printf("Unexpected status: '%s'\n", resp.Status)
			fmt.Printf("Response body: '%s'\n", body)
		}

		err = os.WriteFile("alluserdetails1.json", []byte(body), 0600)
		if err != nil {
			return nil, fmt.Errorf("error saving response body: %w", err)
		}

		var userdetails ArtifactoryUserDetails

		err = json.Unmarshal(body, &userdetails)
		if err != nil {
			return nil, fmt.Errorf("error parsing response body: %w", err)
		}

		alluserdetails = append(alluserdetails, userdetails)
	}

	fmt.Println()

	json, err := json.MarshalIndent(alluserdetails, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error generating json: %w", err)
	}

	err = os.WriteFile("alluserdetails.json", []byte(string(json)), 0600)
	if err != nil {
		return nil, fmt.Errorf("error saving user details: %w", err)
	}

	return alluserdetails, nil
}
