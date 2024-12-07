package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func GetTasks(client *http.Client, baseurl string, token string) ([]ArtifactoryTask, error) {
	url := fmt.Sprintf("%s/artifactory/api/tasks", baseurl)

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

	var tasks ArtifactoryTasks
	err = json.Unmarshal(body, &tasks)
	if err != nil {
		return nil, fmt.Errorf("error parsing response body: %w", err)
	}

	err = os.WriteFile("alltasks.json", body, 0600)
	if err != nil {
		return nil, fmt.Errorf("error saving groups: %w", err)
	}

	return tasks.Tasks, nil
}
