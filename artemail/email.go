package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func FixEmail(userdetails []ArtifactoryUserDetails, client *http.Client, baseurl string, token string, dryRun bool) error {
	for _, user := range userdetails {
		if 1 == 2 {
			fmt.Printf("User: '%s' '%s' '%s' '%s'\n",
				user.Username, user.Email, user.Realm, user.Status)
		}
	}

	fmt.Printf("Count: %d\n", len(userdetails))

	updatecount := 0

	for _, user := range userdetails {
		if !strings.HasSuffix(user.Email, "\n") {
			continue
		}

		url := fmt.Sprintf("%s/access/api/v2/users/%s", baseurl, user.Username)

		artifactoryuseremail := ArtifactoryUserDetailsEmail{
			Email: user.Email[:len(user.Email)-1],
		}

		json, err := json.Marshal(artifactoryuseremail)
		if err != nil {
			return fmt.Errorf("error updating user: %w", err)
		}
		req, err := http.NewRequest("PATCH", url, bytes.NewReader(json))
		if err != nil {
			return fmt.Errorf("error updating user: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		fmt.Printf("Updating: '%s' -> '%s'\n", user.Email, artifactoryuseremail.Email)

		if !dryRun {
			resp, err := client.Do(req)
			if err != nil {
				return fmt.Errorf("error updating user: %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				fmt.Printf("Username: '%s'\n", user.Username)
				fmt.Printf("Url: '%s'\n", url)
				fmt.Printf("Unexpected status: '%s'\n", resp.Status)
				fmt.Printf("Request body: '%s'\n", req.Body)
				body, _ := io.ReadAll(resp.Body)
				fmt.Printf("Response body: '%s'\n", body)
				return fmt.Errorf("error updating user")
			} else {
				fmt.Printf("'%s': Updated user successfully.\n", user.Username)
				updatecount++
			}
		}
	}

	fmt.Printf("Updated: %d\n", updatecount)

	return nil
}
