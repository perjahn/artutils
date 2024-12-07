package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"slices"
	"sort"
	"strings"
)

func CopyUser(
	w io.Writer,
	client *http.Client,
	baseurl string,
	token string,
	sourceusers []string,
	targetuser string,
	permissions []ArtifactoryPermissionDetails,
	dryRun bool) error {

	updatecount := 0

	sort.Slice(permissions, func(i, j int) bool {
		return permissions[i].Name < permissions[j].Name
	})

	for _, permission := range permissions {
		anysourceexists := false
		for _, sourceuser := range sourceusers {
			_, sourceexists := permission.Resources.Artifact.Actions.Users[sourceuser]
			if sourceexists {
				anysourceexists = true
			}
		}
		if !anysourceexists {
			continue
		}

		url := fmt.Sprintf("%s/access/api/v2/permissions/%s/artifact", baseurl, permission.Name)

		users := make(map[string][]string)
		groups := make(map[string][]string)
		targets := make(map[string]ArtifactoryPermissionDetailsTarget)

		maps.Copy(users, permission.Resources.Artifact.Actions.Users)
		maps.Copy(groups, permission.Resources.Artifact.Actions.Groups)
		maps.Copy(targets, permission.Resources.Artifact.Targets)

		for _, sourceuser := range sourceusers {
			_, sourceexists := users[sourceuser]
			if !sourceexists {
				continue
			}

			_, targetexists := users[targetuser]
			if targetexists {
				users[targetuser] = mergeInto(users[sourceuser], users[targetuser])
			} else {
				users[targetuser] = users[sourceuser]
			}
		}

		permissiontarget := ArtifactoryPermissionDetailsArtifact{
			Actions: ArtifactoryPermissionDetailsActions{
				Users:  users,
				Groups: groups,
			},
			Targets: targets,
		}

		if areEqual(permission.Resources.Artifact.Actions.Users[targetuser], permissiontarget.Actions.Users[targetuser]) {
			continue
		}

		sourcestring := ""
		for _, sourceuser := range sourceusers {
			sourcepermissions, sourceexists := permission.Resources.Artifact.Actions.Users[sourceuser]
			if sourceexists {
				if sourcestring != "" {
					sourcestring += ", "
				}
				sourcestring += fmt.Sprintf("'%s' (%s)", sourceuser, strings.Join(sourcepermissions, ", "))
			}
		}
		_, targetexists := permission.Resources.Artifact.Actions.Users[targetuser]
		if !targetexists {
			fmt.Fprintf(w, "Permission adding: '%s': %s: '%s' ( -> %s)\n",
				permission.Name,
				sourcestring,
				targetuser,
				strings.Join(permissiontarget.Actions.Users[targetuser], ", "))
		} else {
			fmt.Fprintf(w, "Permission updating: '%s': %s: '%s' (%s -> %s)\n",
				permission.Name,
				sourcestring,
				targetuser,
				strings.Join(permission.Resources.Artifact.Actions.Users[targetuser], ", "),
				strings.Join(permissiontarget.Actions.Users[targetuser], ", "))
		}

		json, err := json.Marshal(permissiontarget)
		if err != nil {
			return fmt.Errorf("error updating permission target, error generating json: %w", err)
		}

		req, err := http.NewRequest("PUT", url, bytes.NewReader(json))
		if err != nil {
			return fmt.Errorf("error updating permission target, error creating request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		artifactjsonsource := GetArtifactJson(permission.JsonSource)
		PrintDiff(w, artifactjsonsource, string(json))

		if !dryRun {
			resp, err := client.Do(req)
			if err != nil {
				return fmt.Errorf("error updating permission target: %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				fmt.Fprintf(w, "Permission name: '%s'\n", permission.Name)
				fmt.Fprintf(w, "Url: '%s'\n", url)
				fmt.Fprintf(w, "Unexpected status: '%s'\n", resp.Status)
				fmt.Fprintf(w, "Request body: '%s'\n", req.Body)
				body, _ := io.ReadAll(resp.Body)
				fmt.Fprintf(w, "Response body: '%s'\n", body)
				return fmt.Errorf("error updating permission target")
			} else {
				fmt.Fprintf(w, "'%s': Updated permission target successfully.\n", permission.Name)
				updatecount++
			}
		} else {
			updatecount++
		}
	}

	fmt.Fprintf(w, "Update count: %d\n", updatecount)

	return nil
}

func areEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	aCopy := make([]string, len(a))
	copy(aCopy, a)
	bCopy := make([]string, len(b))
	copy(bCopy, b)

	sort.Strings(aCopy)
	sort.Strings(bCopy)

	for i := range aCopy {
		if aCopy[i] != bCopy[i] {
			return false
		}
	}

	return true
}

func mergeInto(src, dest []string) []string {
	for _, s := range src {
		if !slices.Contains(dest, s) {
			dest = append(dest, s)
		}
	}
	return dest
}
