package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

func main() {
	dryRun := flag.Bool("d", false, "Enable dry run mode (read-only, no changes will be made).")

	flag.Parse()
	args := flag.Args()
	if len(args) != 2 || args[0] == "" || args[1] == "" {
		usage()
		os.Exit(1)
	}

	baseurl := getBaseURL(args[0])
	token := getToken(args[1])

	client := &http.Client{}

	err := Generete(client, baseurl, token, *dryRun)
	if err != nil {
		fmt.Println("Error generating:", err)
		os.Exit(1)
	}
}

func getBaseURL(arg string) string {
	var baseurl string

	if envBaseURL := os.Getenv("ARTGEN_BASEURL"); envBaseURL != "" {
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

	if envToken := os.Getenv("ARTGEN_TOKEN"); envToken != "" {
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
	fmt.Println("ARTGEN - Artifactory repos/users/groups/permissions generating tool")
	fmt.Println()
	fmt.Println("This tool is used for generating repos/users/groups/permissions.")
	fmt.Println()
	fmt.Println("Usage: artgen [-d] <baseurl> <tokenfile>")
	fmt.Println()
	fmt.Println("baseurl:    Base URL of Artifactory instance, like https://artifactory.example.com")
	fmt.Println("tokenfile:  File with access token (aka bearer token).")
	fmt.Println()
	flag.PrintDefaults()
	fmt.Println()
	fmt.Println("ARTGEN_BASEURL:  Environment variable that overrides the base URL value.")
	fmt.Println("ARTGEN_TOKEN:    Environment variable that overrides the token value.")
}
