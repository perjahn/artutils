package main

import (
	"fmt"
	"sort"
)

func ShowUsers(userdetails []ArtifactoryUserDetails) {
	sort.Slice(userdetails, func(i, j int) bool {
		return userdetails[i].Username < userdetails[j].Username
	})

	headers := []string{"Username", "Email", "Realm", "Status", "LastLoggedIn"}
	lengths := make([]int, len(headers)-1)
	for i := range len(headers) - 1 {
		lengths[i] = len(headers[i])
	}

	for _, user := range userdetails {
		if len(user.Username) > lengths[0] {
			lengths[0] = len(user.Username)
		}
		if len(user.Email) > lengths[1] {
			lengths[1] = len(user.Email)
		}
		if len(user.Realm) > lengths[2] {
			lengths[2] = len(user.Realm)
		}
		if len(user.Status) > lengths[3] {
			lengths[3] = len(user.Status)
		}
	}

	fmt.Printf("%-*s  %-*s  %-*s  %-*s  %s\n",
		lengths[0], headers[0],
		lengths[1], headers[1],
		lengths[2], headers[2],
		lengths[3], headers[3],
		headers[4])

	for _, user := range userdetails {
		fmt.Printf("%-*s  %-*s  %-*s  %-*s  %s\n",
			lengths[0], user.Username,
			lengths[1], user.Email,
			lengths[2], user.Realm,
			lengths[3], user.Status,
			user.LastLoggedIn)
	}

	fmt.Printf("Count: %d\n", len(userdetails))
}
