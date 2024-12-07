package main

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
)

func UpdateRepos(
	client *http.Client,
	baseurl string,
	token string,
	reponame string,
	artifacts []ArtifactoryArtifact,
	propertyname string,
	useCommonNodes bool,
	dryRun bool) error {

	updatecount := 0
	var commonnodes []string

	for _, artifact := range artifacts {
		commonnodes = append(commonnodes, artifact.Path)
	}

	if useCommonNodes {
		commonnodes = reduceToCommonNodes(commonnodes)
	}

	for _, artifact := range artifacts {
		ancestornode := ""
		for _, commonnode := range commonnodes {
			if strings.HasPrefix(artifact.Path+"/", commonnode+"/") && len(commonnode) > len(ancestornode) {
				ancestornode = commonnode
			}
		}

		var oldpropvalues []string
		for _, prop := range artifact.Properties {
			if prop.Key == propertyname {
				oldpropvalues = append(oldpropvalues, sanitize(prop.Value))
			}
		}

		foldername := filepath.Base(ancestornode)

		if len(oldpropvalues) == 1 && oldpropvalues[0] == foldername {
			fmt.Printf("No diff: %s/%s: '%s': '%s'\n", artifact.Path, artifact.Name, propertyname, foldername)
			continue
		}

		url := fmt.Sprintf("%s/artifactory/api/storage/%s/%s/%s?properties=%s=%s", baseurl, reponame, artifact.Path, artifact.Name, propertyname, foldername)
		req, err := http.NewRequest("PUT", url, nil)
		if err != nil {
			return fmt.Errorf("error creating request: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+token)

		if len(oldpropvalues) > 0 {
			prev := strings.Join(oldpropvalues, "', '")
			reduceinfo := ""
			if len(oldpropvalues) > 1 {
				reduceinfo = fmt.Sprintf(" (%d -> 1)", len(oldpropvalues))
			}
			fmt.Printf("Updating%s: %s/%s: '%s': '%s' -> '%s'\n", reduceinfo, artifact.Path, artifact.Name, propertyname, prev, foldername)
		} else {
			fmt.Printf("Adding: %s/%s: '%s': '%s'\n", artifact.Path, artifact.Name, propertyname, foldername)
		}

		if !dryRun {
			resp, err := client.Do(req)
			if err != nil {
				return fmt.Errorf("error updating property: %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 && resp.StatusCode != 204 {
				fmt.Printf("ArtifactPath: '%s'\n", artifact.Path)
				fmt.Printf("ArtifactName: '%s'\n", artifact.Name)
				fmt.Printf("PropertyName: '%s'\n", propertyname)
				fmt.Printf("FolderName: '%s'\n", foldername)
				fmt.Printf("Url: '%s'\n", url)
				fmt.Printf("Unexpected status: '%s'\n", resp.Status)
				fmt.Printf("Request body: '%s'\n", req.Body)
				body, _ := io.ReadAll(resp.Body)
				fmt.Printf("Response body: '%s'\n", body)
				return fmt.Errorf("error updating property")
			} else {
				fmt.Printf("Updated property successfully.\n")
				updatecount++
			}
		} else {
			updatecount++
		}
	}

	fmt.Printf("Update count: %d\n", updatecount)

	return nil
}

func reduceToCommonNodes(nodes []string) []string {
	var reducedNodes []string

	reducedNodes = append(reducedNodes, nodes...)

	goon := true
	for goon {
		goon = false
		for i := 0; i < len(reducedNodes); i++ {
			for j := 0; j < len(reducedNodes); j++ {
				if i != j && strings.HasPrefix(reducedNodes[i]+"/", reducedNodes[j]+"/") && len(reducedNodes[j]) <= len(reducedNodes[i]) {
					reducedNodes = append(reducedNodes[:i], reducedNodes[i+1:]...)
					goon = true
					i--
					break
				}
			}
		}
	}

	return reducedNodes
}

func sanitize(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r <= 32 {
			b.WriteByte('.')
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}
