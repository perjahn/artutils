package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func GetStuff(
	client *http.Client,
	baseurl string,
	token string,
	reponame string) ([]ArtifactoryArtifact, error) {

	artifacts, err := getProperties(client, baseurl, token, reponame)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Artifact count: %d\n", len(artifacts))

	totalProps := 0
	for _, a := range artifacts {
		totalProps += len(a.Properties)
	}
	fmt.Printf("Total artifact properties: %d\n", totalProps)

	return artifacts, nil
}

func getProperties(
	client *http.Client,
	baseurl string,
	token string,
	reponame string) ([]ArtifactoryArtifact, error) {

	fmt.Println("Getting properties...")

	url := fmt.Sprintf("%s/artifactory/api/search/aql", baseurl)
	query := `items.find({"repo":"` + reponame + `","type":"file"}).include("name","path","property")`

	req, err := http.NewRequest("POST", url, strings.NewReader(query))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "text/plain")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error posting query: %w", err)
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

	err = os.WriteFile("properties.json", []byte(body), 0600)
	if err != nil {
		return nil, fmt.Errorf("error saving response body: %w", err)
	}

	var properties ArtifactoryProperties

	err = json.Unmarshal(body, &properties)
	if err != nil {
		return nil, fmt.Errorf("error parsing response body: %w", err)
	}

	return properties.Results, nil
}
