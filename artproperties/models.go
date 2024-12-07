package main

type ArtifactoryProperties struct {
	Results []ArtifactoryArtifact `json:"results"`
}

type ArtifactoryArtifact struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	Properties []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"properties,omitempty"`
}
