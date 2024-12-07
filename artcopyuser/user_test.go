package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func mockHTTPClient(fn func(*http.Request) (*http.Response, error)) *http.Client {
	return &http.Client{
		Transport: roundTripperFunc(fn),
	}
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestCopyUser(t *testing.T) {
	tests := []struct {
		sourceusers       []string
		targetuser        string
		permissiondetails []ArtifactoryPermissionDetails
		dryRun            bool
		payload           ArtifactoryPermissionDetailsArtifact
		HTTPResponse      string
		consoleOutput     string
	}{
		{[]string{}, "", []ArtifactoryPermissionDetails{}, false, ArtifactoryPermissionDetailsArtifact{}, "", "Update count: 0\n"},
		{[]string{}, "", []ArtifactoryPermissionDetails{}, true, ArtifactoryPermissionDetailsArtifact{}, "", "Update count: 0\n"},
		{
			sourceusers: []string{"test-user1", "test-user2"},
			targetuser:  "test-user3",
			permissiondetails: []ArtifactoryPermissionDetails{
				{
					Name: "test-perm",
					Resources: ArtifactoryPermissionDetailsResources{
						Artifact: ArtifactoryPermissionDetailsArtifact{
							Actions: ArtifactoryPermissionDetailsActions{
								Users: map[string][]string{
									"test-user1": {"READ"},
									"test-user2": {"WRITE"},
									"test-user3": {"OTHER"},
								},
								Groups: map[string][]string{
									"test-group": {"READ"},
								},
							},
							Targets: map[string]ArtifactoryPermissionDetailsTarget{
								"test-repo": {
									IncludePatterns: []string{"**"},
									ExcludePatterns: []string{},
								},
							},
						},
					},
				},
			},
			dryRun: false,
			payload: ArtifactoryPermissionDetailsArtifact{
				Actions: ArtifactoryPermissionDetailsActions{
					Users: map[string][]string{
						"test-user1": {"READ"},
						"test-user2": {"WRITE"},
						"test-user3": {"OTHER", "READ", "WRITE"},
					},
					Groups: map[string][]string{
						"test-group": {"READ"},
					},
				},
				Targets: map[string]ArtifactoryPermissionDetailsTarget{
					"test-repo": {
						IncludePatterns: []string{"**"},
						ExcludePatterns: []string{},
					},
				},
			},
			HTTPResponse: `{"ok":true}`,
			consoleOutput: `Permission updating: 'test-perm': 'test-user1' (READ), 'test-user2' (WRITE): 'test-user3' (OTHER -> OTHER, READ, WRITE)
'test-perm': Updated permission target successfully.
Update count: 1
`},
	}
	for i, tc := range tests {
		var client *http.Client
		client = mockHTTPClient(func(req *http.Request) (*http.Response, error) {
			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("read request body: %v", err)
			}
			_ = req.Body.Close()

			if len(bodyBytes) == 0 {
				var empty ArtifactoryPermissionDetailsArtifact
				if !reflect.DeepEqual(tc.payload, empty) {
					t.Fatalf("unexpected empty request body, want payload: %+v", tc.payload)
				}
			} else {
				var got any
				if err := json.Unmarshal(bodyBytes, &got); err != nil {
					t.Fatalf("invalid JSON in request body: %v\nbody: %s", err, string(bodyBytes))
				}

				expectedBytes, err := json.Marshal(tc.payload)
				if err != nil {
					t.Fatalf("marshal expected payload: %v", err)
				}
				var want any
				if err := json.Unmarshal(expectedBytes, &want); err != nil {
					t.Fatalf("invalid JSON for expected payload: %v\nexpected: %s", err, string(expectedBytes))
				}

				if !reflect.DeepEqual(got, want) {
					t.Fatalf("unexpected request payload:\n got:  %s\n want: %s", string(bodyBytes), string(expectedBytes))
				}
			}

			var reponse *http.Response
			reponse = &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(tc.HTTPResponse)),
				Header:     make(http.Header),
			}

			return reponse, nil
		})
		var out bytes.Buffer

		err := CopyUser(&out, client, "", "", tc.sourceusers, tc.targetuser, tc.permissiondetails, tc.dryRun)
		if err != nil {
			t.Errorf("CopyUser (%d/%d): error = %v", i, len(tests), err)
		}

		got := out.String()
		if strings.TrimSpace(got) != strings.TrimSpace(tc.consoleOutput) {
			t.Errorf("CopyUser (%d/%d): console output mismatch\n got: %q\nwant: %q", i, len(tests), got, tc.consoleOutput)
		}
	}
}
