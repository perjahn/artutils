package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	dryRun := flag.Bool("d", false, "Enable dry run mode (read-only, no changes will be made).")

	flag.Parse()
	args := flag.Args()
	if len(args) != 4 || args[0] == "" || args[1] == "" || args[2] == "" || args[3] == "" {
		usage()
		os.Exit(1)
	}

	baseurl := getBaseURL(args[0])
	token := getToken(args[1])
	sourceusers := strings.Split(args[2], ",")
	targetuser := args[3]

	client := &http.Client{}

	permissions, err := GetStuff(client, baseurl, token)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	CopyUser(os.Stdout, client, baseurl, token, sourceusers, targetuser, permissions, *dryRun)
}

func getBaseURL(arg string) string {
	var baseurl string

	if envBaseURL := os.Getenv("ARTCOPYUSER_BASEURL"); envBaseURL != "" {
		baseurl = envBaseURL
	} else {
		baseurl = arg
	}
	if baseurl == "" {
		fmt.Println("Error: Base URL is empty.")
		os.Exit(1)
	}
	if _, err := url.ParseRequestURI(baseurl); err != nil {
		fmt.Printf("Error: Invalid base URL '%s': %v\n", baseurl, err)
		os.Exit(1)
	}

	return baseurl
}

func getToken(arg string) string {
	var token string

	if envToken := os.Getenv("ARTCOPYUSER_TOKEN"); envToken != "" {
		token = envToken
	} else {
		data, err := os.ReadFile(arg)
		if err != nil {
			fmt.Printf("Error reading token file: %v\n", err)
			os.Exit(1)
		}
		token = string(data)
	}
	if token == "" {
		fmt.Println("Error: Token is empty.")
		os.Exit(1)
	}

	return token
}

func usage() {
	fmt.Println("ARTCOPYUSER - Artifactory user tool")
	fmt.Println()
	fmt.Println("This tool is used for copying permissions for a user to another user.")
	fmt.Println()
	fmt.Println("Usage: artcopyuser [-d] <baseurl> <tokenfile> <sourceusers> <targetuser>")
	fmt.Println()
	fmt.Println("baseurl:      Base URL of Artifactory instance, like https://artifactory.example.com")
	fmt.Println("tokenfile:    File with access token (aka bearer token).")
	fmt.Println("sourceusers:  Comma separated list of source users.")
	fmt.Println("targetuser:   The name of the target user.")
	fmt.Println()
	flag.PrintDefaults()
	fmt.Println()
	fmt.Println("ARTCOPYUSER_BASEURL:  Environment variable that overrides the base URL value.")
	fmt.Println("ARTCOPYUSER_TOKEN:    Environment variable that overrides the token value.")
}
