package main

import "time"

type ArtifactoryTasks struct {
	Tasks []ArtifactoryTask `json:"tasks"`
}

type ArtifactoryTask struct {
	Id          *string    `json:"id"`
	Type        *string    `json:"type"`
	State       *string    `json:"state"`
	Description *string    `json:"description"`
	Started     *time.Time `json:"started"`
}
