package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
)

func main() {
	args := os.Args
	if len(args) != 3 || args[1] == "" || args[2] == "" {
		usage()
		os.Exit(1)
	}

	baseurl := getBaseURL(args[1])
	token := getToken(args[2])

	client := &http.Client{}

	repos, permissiondetails, err := GetStuff(client, baseurl, token)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	ShowPatterns(repos, permissiondetails)
}

func getBaseURL(arg string) string {
	var baseurl string

	if envBaseURL := os.Getenv("ARTPATTERNS_BASEURL"); envBaseURL != "" {
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

	if envToken := os.Getenv("ARTPATTERNS_TOKEN"); envToken != "" {
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
	fmt.Println("ARTPATTERNS - Artifactory permission targets patterns tool")
	fmt.Println()
	fmt.Println("This tool is used for visualizing non-default permission target patterns.")
	fmt.Println()
	fmt.Println("Usage: artpatterns <baseurl> <tokenfile>")
	fmt.Println()
	fmt.Println("baseurl:    Base URL of Artifactory instance, like https://artifactory.example.com")
	fmt.Println("tokenfile:  File with access token (aka bearer token).")
	fmt.Println()
	fmt.Println("ARTPATTERNS_BASEURL:  Environment variable that overrides the base URL value.")
	fmt.Println("ARTPATTERNS_TOKEN:    Environment variable that overrides the token value.")
}
