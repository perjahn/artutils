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
	if len(args) != 8 || args[0] == "" || args[1] == "" || args[2] == "" || args[3] == "" || args[4] == "" || args[5] == "" || args[6] == "" || args[7] == "" {
		usage()
		os.Exit(1)
	}

	baseurl := getBaseURL(args[0])
	token := getToken(args[1])
	username := args[2]
	password := getPassword(args[3])
	groupname := args[4]
	ldapgroupsettingsName := args[5]
	ldapusername := args[6]
	ldappassword := getLdapPassword(args[7])

	client := &http.Client{}

	ldapsettings, err := GetLDAPSettings(client, baseurl, token)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	ldapgroupsettings, err := GetLDAPGroupSettings(client, baseurl, token)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	err = ImportGroup(client, baseurl, username, password, groupname, ldapsettings, ldapgroupsettings, ldapgroupsettingsName, ldapusername, ldappassword, *dryRun)
	if err != nil {
		fmt.Printf("Error importing groups: %v\n", err)
		os.Exit(1)
	}
}

func getBaseURL(arg string) string {
	var baseurl string

	if envBaseURL := os.Getenv("ARTGROUPIMPORT_BASEURL"); envBaseURL != "" {
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

	if envToken := os.Getenv("ARTGROUPIMPORT_TOKEN"); envToken != "" {
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

func getPassword(arg string) string {
	var token string

	if envPassword := os.Getenv("ARTGROUPIMPORT_PASSWORD"); envPassword != "" {
		token = envPassword
	} else {
		data, err := os.ReadFile(arg)
		if err != nil {
			fmt.Printf("Error reading password file: %v\n", err)
			os.Exit(1)
		}
		token = string(data)
	}
	if token == "" {
		fmt.Println("Error: Password is empty.")
		os.Exit(1)
	}

	return token
}

func getLdapPassword(arg string) string {
	var token string

	if envLdapPassword := os.Getenv("ARTGROUPIMPORT_LDAPPASSWORD"); envLdapPassword != "" {
		token = envLdapPassword
	} else {
		data, err := os.ReadFile(arg)
		if err != nil {
			fmt.Printf("Error reading ldap password file: %v\n", err)
			os.Exit(1)
		}
		token = string(data)
	}
	if token == "" {
		fmt.Println("Error: Ldap password is empty.")
		os.Exit(1)
	}

	return token
}

func usage() {
	fmt.Println("ARTGROUPIMPORT - Artifactory group import tool")
	fmt.Println()
	fmt.Println("This tool is used to sync Artifactory groups from ldap.")
	fmt.Println("This tool requires Artifactory 7.79.")
	fmt.Println()
	fmt.Println("Usage: artgroupimport [-d] <baseurl> <tokenfile> <username> <passwordfile> <groupname> <ldapgroupsettingsname> <ldapusername> <ldappasswordfile>")
	fmt.Println()
	fmt.Println("baseurl:                Base URL of Artifactory instance, like https://artifactory.example.com")
	fmt.Println("tokenfile:              File with access token (aka bearer token).")
	fmt.Println("username:               Name of the artifactory user.")
	fmt.Println("passwordfile:           File with password for the artifactory user.")
	fmt.Println("groupname:              Name of the group in the ldap directory to import.")
	fmt.Println("ldapgroupsettingsname:  Name of the Artifactory ldap group settings, which contains the group to import.")
	fmt.Println("ldapusername:           DN of the ldap user to make queries.")
	fmt.Println("ldappassword:           File with password of the ldap user to make queries.")
	fmt.Println()
	flag.PrintDefaults()
	fmt.Println()
	fmt.Println("ARTGROUPIMPORT_BASEURL:       Environment variable that overrides the base URL value.")
	fmt.Println("ARTGROUPIMPORT_TOKEN:         Environment variable that overrides the token value.")
	fmt.Println("ARTGROUPIMPORT_PASSWORD:      Environment variable that overrides the password value.")
	fmt.Println("ARTGROUPIMPORT_LDAPPASSWORD:  Environment variable that overrides the ldappassword value.")
}
