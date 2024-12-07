package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"strings"

	ldap "github.com/go-ldap/ldap/v3"
)

func main() {
	bindDN := flag.String("binddn", "", "Bind DN (if required)")
	bindPW := flag.String("bindpw", "", "Bind password (if required)")
	attrs := flag.String("attrs", "dn,cn,mail", "Comma-separated list of attributes to request")
	startTLS := flag.Bool("starttls", false, "Use StartTLS (for ldap:// URLs)")
	insecure := flag.Bool("insecure", false, "Skip TLS certificate verification (insecure)")
	flag.Parse()

	args := flag.Args()

	if len(args) != 3 || args[0] == "" || args[1] == "" || args[2] == "" {
		usage()
		log.Fatal("server, base and filter are required")
	}

	server := args[0]
	base := args[1]
	filter := args[2]

	attrList := []string{}
	for a := range strings.SplitSeq(*attrs, ",") {
		a = strings.TrimSpace(a)
		if a != "" {
			attrList = append(attrList, a)
		}
	}

	if err := queryldap(server, base, filter, *bindDN, *bindPW, attrList, *startTLS, *insecure); err != nil {
		log.Fatalf("ldap query failed: %v", err)
	}
}

func usage() {
	fmt.Println("LDAPQUERY - ldap query tool")
	fmt.Println()
	fmt.Println("Usage: ldapquery [flags] <server> <baseDN> <filter>")
	fmt.Println()
	flag.PrintDefaults()
	fmt.Println()
	fmt.Println("server:  LDAP server URL, e.g. ldap://localhost:389 or ldaps://ldap.example:636")
	fmt.Println("filter:  LDAP filter to run, e.g. (uid=john)")
	fmt.Println("base:    Base DN for the search, e.g. dc=example,dc=com")
}

func queryldap(server, baseDN, filter, bindDN, bindPW string, attrs []string, startTLS, insecure bool) error {
	var l *ldap.Conn
	var err error

	if strings.HasPrefix(strings.ToLower(server), "ldaps://") {
		l, err = ldap.DialURL(server, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: insecure}))
		if err != nil {
			return fmt.Errorf("failed to dial LDAPS: %w", err)
		}
	} else {
		l, err = ldap.DialURL(server)
		if err != nil {
			return fmt.Errorf("failed to dial LDAP: %w", err)
		}
		if startTLS {
			if err = l.StartTLS(&tls.Config{InsecureSkipVerify: insecure}); err != nil {
				l.Close()
				return fmt.Errorf("StartTLS failed: %w", err)
			}
		}
	}
	defer l.Close()

	if bindDN != "" {
		if err = l.Bind(bindDN, bindPW); err != nil {
			return fmt.Errorf("bind failed: %w", err)
		}
	}

	searchReq := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		filter,
		attrs,
		nil,
	)

	sr, err := l.Search(searchReq)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	fmt.Printf("Got %d entries\n", len(sr.Entries))
	fmt.Printf("Got %d controls\n", len(sr.Controls))
	fmt.Printf("Got %d referrals\n", len(sr.Referrals))

	for _, entry := range sr.Entries {
		fmt.Printf("DN: %s\n", entry.DN)
		for _, a := range attrs {
			vals := entry.GetAttributeValues(a)
			if len(vals) == 0 {
				continue
			}
			fmt.Printf("  %s: %v\n", a, vals)
		}
		fmt.Println()
	}

	return nil
}
