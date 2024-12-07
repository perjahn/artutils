package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

func main() {
	useCommonNodes := flag.Bool("c", false, "Only use common ancestor folders, i.e. for all artifacts, make sure they traverse path to root and use common folders if possible.")
	dryRun := flag.Bool("d", false, "Enable dry run mode (read-only, no changes will be made).")

	flag.Parse()
	args := flag.Args()
	if len(args) != 4 || args[0] == "" || args[1] == "" || args[2] == "" || args[3] == "" {
		usage()
		os.Exit(1)
	}

	baseurl := getBaseURL(args[0])
	token := getToken(args[1])
	reponame := args[2]
	propertyname := args[3]

	client := &http.Client{}

	artifacts, err := GetStuff(client, baseurl, token, reponame)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	UpdateRepos(client, baseurl, token, reponame, artifacts, propertyname, *useCommonNodes, *dryRun)
}

func getBaseURL(arg string) string {
	var baseurl string

	if envBaseURL := os.Getenv("ARTPROPERTIES_BASEURL"); envBaseURL != "" {
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

	if envToken := os.Getenv("ARTPROPERTIES_TOKEN"); envToken != "" {
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
	fmt.Println("ARTPROPERTIES - Artifactory properties tool")
	fmt.Println()
	fmt.Println("This tool is used for adding the parent folder name as property to every artifact in a repo.")
	fmt.Println()
	fmt.Println("Usage: artproperties [-c] [-d] <baseurl> <tokenfile> <reponame> <propertyname>")
	fmt.Println()
	fmt.Println("baseurl:       Base URL of Artifactory instance, like https://artifactory.example.com")
	fmt.Println("tokenfile:     File with access token (aka bearer token).")
	fmt.Println("reponame:      The name of the repo. Isn't this a useless description btw? Yes, 1000 times yes.")
	fmt.Println("propertyname:  The name of the property. See above.")
	fmt.Println()
	flag.PrintDefaults()
	fmt.Println()
	fmt.Println("ARTPROPERTIES_BASEURL:  Environment variable that overrides the base URL value.")
	fmt.Println("ARTPROPERTIES_TOKEN:    Environment variable that overrides the token value.")
}
