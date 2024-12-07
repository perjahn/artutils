package main

import (
	"slices"
	"testing"
)

func TestMergePermissions_BothEmpty(t *testing.T) {
	a := Permissions{}
	b := Permissions{}
	result := mergePermissions(a, b)

	if len(result.Users) != 0 || len(result.Groups) != 0 {
		t.Errorf("Expected empty, got Users: %v, Groups: %v", result.Users, result.Groups)
	}
}

func TestMergePermissions_AEmpty(t *testing.T) {
	a := Permissions{}
	b := Permissions{Users: []string{"user1"}, Groups: []string{"group1"}}
	result := mergePermissions(a, b)

	if !slices.Equal(result.Users, []string{"user1"}) || !slices.Equal(result.Groups, []string{"group1"}) {
		t.Errorf("Mismatch. Got Users: %v, Groups: %v", result.Users, result.Groups)
	}
}

func TestMergePermissions_BEmpty(t *testing.T) {
	a := Permissions{Users: []string{"user1"}, Groups: []string{"group1"}}
	b := Permissions{}
	result := mergePermissions(a, b)

	if !slices.Equal(result.Users, []string{"user1"}) || !slices.Equal(result.Groups, []string{"group1"}) {
		t.Errorf("Mismatch. Got Users: %v, Groups: %v", result.Users, result.Groups)
	}
}

func TestMergePermissions_NoDuplicates(t *testing.T) {
	a := Permissions{Users: []string{"user1"}, Groups: []string{"group1"}}
	b := Permissions{Users: []string{"user2"}, Groups: []string{"group2"}}
	result := mergePermissions(a, b)

	if !slices.Equal(result.Users, []string{"user1", "user2"}) || !slices.Equal(result.Groups, []string{"group1", "group2"}) {
		t.Errorf("Mismatch. Got Users: %v, Groups: %v", result.Users, result.Groups)
	}
}

func TestMergePermissions_WithDuplicates(t *testing.T) {
	a := Permissions{Users: []string{"user1", "user2"}, Groups: []string{"group1"}}
	b := Permissions{Users: []string{"user2", "user3"}, Groups: []string{"group1", "group2"}}
	result := mergePermissions(a, b)

	expected_users := []string{"user1", "user2", "user3"}
	expected_groups := []string{"group1", "group2"}

	if !slices.Equal(result.Users, expected_users) || !slices.Equal(result.Groups, expected_groups) {
		t.Errorf("Mismatch. Got Users: %v (expected %v), Groups: %v (expected %v)", result.Users, expected_users, result.Groups, expected_groups)
	}
}

func TestMergePermissions_Sorted(t *testing.T) {
	a := Permissions{Users: []string{"zebra", "apple"}}
	b := Permissions{Users: []string{"monkey"}}
	result := mergePermissions(a, b)

	expected := []string{"apple", "monkey", "zebra"}
	if !slices.Equal(result.Users, expected) {
		t.Errorf("Not sorted correctly. Got %v, expected %v", result.Users, expected)
	}

	if !slices.IsSorted(result.Users) {
		t.Errorf("Result not sorted: %v", result.Users)
	}
}

func TestMergePermissions_NoDuplicatesInResult(t *testing.T) {
	a := Permissions{Users: []string{"user1", "user1", "user2"}}
	b := Permissions{Users: []string{"user1", "user2", "user2"}}
	result := mergePermissions(a, b)

	expected := []string{"user1", "user2"}
	if !slices.Equal(result.Users, expected) {
		t.Errorf("Duplicates not removed. Got %v, expected %v", result.Users, expected)
	}

	// Check no duplicates in result
	seen := make(map[string]bool)
	for _, u := range result.Users {
		if seen[u] {
			t.Errorf("Found duplicate user in result: %s", u)
		}
		seen[u] = true
	}
}
