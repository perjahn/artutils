package main

import (
	"encoding/xml"
)

type ArtifactoryFile struct {
	XMLName        xml.Name       `xml:"artifactory-file"`
	Size           int64          `xml:"size"`
	AdditionalInfo AdditionalInfo `xml:"additionalInfo"`
}

type AdditionalInfo struct {
	ChecksumsInfo ChecksumsInfo `xml:"checksumsInfo"`
}

type ChecksumsInfo struct {
	Checksums Checksums `xml:"checksums"`
}

type Checksums struct {
	Checksum []Checksum `xml:"checksum"`
}

type Checksum struct {
	Type   string `xml:"type"`
	Actual string `xml:"actual"`
}
