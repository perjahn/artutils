package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func GetStuff(client *http.Client, baseurl string, token string) ([]ArtifactoryRepoDetailsResponse, error) {
	repos, err := getRepos(client, baseurl, token)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Repo count: %d\n", len(repos))

	repodetails, err := getRepoDetails(client, baseurl, token, repos)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Repo details count: %d\n", len(repodetails))

	return repodetails, nil
}

func getRepos(client *http.Client, baseurl string, token string) ([]ArtifactoryRepoResponse, error) {
	url := fmt.Sprintf("%s/artifactory/api/repositories", baseurl)
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

	err = os.WriteFile("allrepos.json", []byte(body), 0600)
	if err != nil {
		return nil, fmt.Errorf("error saving response body: %w", err)
	}

	var repos []ArtifactoryRepoResponse
	err = json.Unmarshal(body, &repos)
	if err != nil {
		return nil, fmt.Errorf("error parsing response body: %w", err)
	}

	return repos, nil
}

func getRepoDetails(client *http.Client, baseurl string, token string, repos []ArtifactoryRepoResponse) ([]ArtifactoryRepoDetailsResponse, error) {
	var allrepodetails []ArtifactoryRepoDetailsResponse

	for _, repo := range repos {
		fmt.Print(".")

		url := fmt.Sprintf("%s/artifactory/api/repositories/%s", baseurl, url.PathEscape(repo.Key))
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
			fmt.Printf("Url2: '%s'\n", repo.Url)
			fmt.Printf("Unexpected status: '%s'\n", resp.Status)
			fmt.Printf("Response body: '%s'\n", body)
		}

		err = os.WriteFile("allrepodetails1.json", []byte(body), 0600)
		if err != nil {
			return nil, fmt.Errorf("error saving response body: %w", err)
		}

		var repodetails ArtifactoryRepoDetailsResponse

		err = json.Unmarshal(body, &repodetails)
		if err != nil {
			return nil, fmt.Errorf("error parsing response body: %w", err)
		}

		allrepodetails = append(allrepodetails, repodetails)
	}

	fmt.Println()

	json, err := json.MarshalIndent(allrepodetails, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error generating json: %w", err)
	}

	err = os.WriteFile("allrepodetails.json", []byte(string(json)), 0600)
	if err != nil {
		return nil, fmt.Errorf("error saving repo details: %w", err)
	}

	return allrepodetails, nil
}
