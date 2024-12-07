# Artifactory Redundant Permission Target Detector

A Go application that analyzes Artifactory permission targets to identify and remove redundant ones based on effective repository permissions.

## Overview

This tool analyzes a JSON file containing Artifactory permission targets and detects which ones are redundant - meaning they don't change the effective permissions of any repository if removed (yes, that was an em dash!).

### How It Works

1. **Loads Permission Targets**: Reads a JSON file containing an array of Artifactory permission target objects
2. **Collects All Repositories**: Extracts all unique repository names mentioned across all permission targets
3. **Calculates Effective Permissions**: For each repository, combines permissions from all applicable permission targets (using OR logic - if any target grants a permission, it's granted)
4. **Identifies Non-Matching Targets**: Finds permission targets whose names don't match any repository name
5. **Detects Redundancy**: For each non-matching target, simulates removing it and checks if any repo's effective permissions would change
6. **Reports Results**: Lists all redundant permission targets that can be safely removed

## Building the Application

```bash
go build
```

## Running the Application

```bash
./artremoveredundant permission-targets.json
```

## JSON Input Format

The input JSON should be an array of permission target objects matching Artifactory's API structure:

```json
[
  {
    "name": "repo1",
    "resources": {
      "artifact": {
        "actions": {
          "users": {
            "user1": [
              "READ"
            ],
            "user2": [
              "WRITE"
            ]
          },
          "groups": {
            "group1": [
              "READ"
            ],
            "group2": [
              "ANNOTATE"
            ]
          }
        },
        "targets": {
          "repo1": {
            "include_patterns": [
              "**"
            ],
            "exclude_patterns": []
          },
          "repo2": {
            "include_patterns": [
              "**"
            ],
            "exclude_patterns": []
          }
        }
      }
    }
  },
```

In this example, the permission target is named `repo1`, i.e. has the same name as the one of the repos it's connected to.

## Output Example

```
Loaded permission targets: 5
Found unique repositories: 2
Permission targets with names not matching any repo: 2
Detected redundant permission targets: 2
´developers-read´
´developers-write´
´qa-permissions´
´redundant-duplicate´
```

## Features

- **Comprehensive Permission Analysis**: Considers the usual 5, and all other as well
- **Accurate Redundancy Detection**: Only flags targets that truly don't affect repository permissions
- **Clear Output**: Detailed breakdown of permissions and redundancy findings

## Example Use Case

A typical scenario where redundancy detection is useful:

1. You have three permission targets:
   - `dev-read`: grants read to `libs-release` and `libs-snapshot`
   - `dev-write`: grants write to `libs-release` only
   - `admin`: grants all permissions to all repos

2. The `dev-read` and `dev-write` targets might be redundant if their permissions are already covered by the `admin` target or if they grant identical permissions to multiple targets that could be consolidated.

3. This tool identifies and reports such redundancies, allowing you to clean up and simplify your permission structure.

## Sample Data

A sample `permission-targets.json` file is provided for testing. It demonstrates:
- Multiple permission targets with various scopes
- Different permission levels
- Scenario with actual and potentially redundant targets

## Algorithm Details

### Permission Combination
- Multiple permission targets can apply to the same repository
- Permissions are combined using logical OR (if any target grants read, the effective permission includes read)
- More restrictive targets cannot override more permissive ones

### Redundancy Detection
- A permission target is flagged as redundant if removing it doesn't change any repository's effective permissions
- Only checks targets whose names don't match any repository name (repo-specific targets are never flagged as redundant)
- Tests against all unique repositories found in the targets
- Uses an **iterative greedy algorithm**: Once a redundant target is identified, it's removed from the working set before testing the next target. This prevents false negatives with codependent targets.

#### Codependent Targets Handling
The redundancy detection algorithm uses an iterative approach to correctly handle cases where permission targets are codependent (i.e., multiple targets together provide permissions that no single target provides alone). 

For example:
- PT1 grants read permission to repo-A
- PT2 grants write permission to repo-A  
- PT3 grants both read and write to repo-A

Without iteration: All three might be marked as redundant (incorrect)
With iteration: PT3 is identified as redundant first and removed, then PT1 and PT2 are correctly identified as necessary for their independent permissions.

The algorithm repeats until no more redundancies are found, ensuring codependent targets are correctly classified.

## Limitations

- Requires permission target data in Artifactory API format
- Does not currently analyze include or exclude permissions (only repo)
- Does not currently analyze build or release_bundle permissions (only repo)
- Assumes permission data is valid and consistent

## Exit Codes

- `0`: Success
- `1`: File not found or invalid JSON

## Requirements

- Go 1.21 or later
