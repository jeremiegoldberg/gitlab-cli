# gitlab-manager

A command-line tool for managing GitLab merge requests with built-in changelog validation and merge blocking capabilities.

## Features

- Merge Request Management
  - Create, read, update, and close merge requests
  - Block/unblock merge requests
  - Validate changelog entries
  - Track linked issues
  - JSON output support
- GitLab CI Integration
  - Automatic project detection
  - Pipeline-based validations
  - Merge blocking
- Changelog Validation
  - Supports multiple entry types
  - Checks both MR and linked issues
  - CI/CD integration
- Merge Request Blocking
  - Block with custom reasons
  - Track block status
  - Automatic title updates

## Installation

```bash
go install gitlab.com/your-org/gitlab-manager@latest
```

## Configuration

The tool requires a GitLab access token:

```bash
# For local development
export GITLAB_TOKEN=your_personal_access_token

# In GitLab CI, CI_JOB_TOKEN is used automatically
```

## Usage

### Merge Request Management

```bash
# List merge requests
gitlab-manager mr list
gitlab-manager mr list --state opened --target main

# Get merge request details
gitlab-manager mr get -m 123 [--json]

# Create merge request
gitlab-manager mr create \
  --source feature-branch \
  --target main \
  --title "Add new feature" \
  --description "Implements..."

# Get linked issues
gitlab-manager mr get-issues -m 123 [--json]

# Get MR description
gitlab-manager mr get-description -m 123
```

### Changelog Validation

The tool enforces changelog entries in either the MR description or linked issues.
Valid changelog entries must start with one of:

- [Feature] - New features
- [Improvement] - Improvements to existing features
- [Fix] - Bug fixes
- [Infra] - Infrastructure changes
- [No-Changelog-Entry] - Skip changelog requirement

```bash
# Check changelog entries
gitlab-manager mr check-changelog -m 123

# Example valid entries:
[Feature] Add new user authentication system
[Fix] Resolve login page redirect issue
[Improvement] Enhance search performance
[Infra] Upgrade PostgreSQL to 14
[No-Changelog-Entry] Internal refactoring
```

### Merge Request Blocking

Control merge request status with blocking:

```bash
# Block a merge request
gitlab-manager mr block -m 123 -r "Needs security review"

# Unblock when ready
gitlab-manager mr unblock -m 123
```

When blocked:
- MR title is prefixed with [BLOCKED]
- A note is added with the blocking reason
- Merge operations are prevented

### CI/CD Integration

Example GitLab CI configuration:

```yaml
validate_merge_request:
  script:
    # Validate changelog
    - |
      if ! gitlab-manager mr check-changelog -m $CI_MERGE_REQUEST_IID; then
        echo "Missing changelog entry"
        gitlab-manager mr block -m $CI_MERGE_REQUEST_IID -r "Missing changelog entry"
        exit 1
      fi

    # Other validations
    - gitlab-manager mr get-issues -m $CI_MERGE_REQUEST_IID
    
  rules:
    - if: $CI_PIPELINE_SOURCE == "merge_request_event"
```

### JSON Output

Many commands support JSON output for integration with other tools:

```bash
# Get MR details in JSON
gitlab-manager mr get -m 123 --json

# Example output:
{
  "iid": 123,
  "title": "Add new feature",
  "description": "[Feature] Implement...",
  "state": "opened",
  "source_branch": "feature-branch",
  "target_branch": "main",
  "created_at": "2024-01-15T10:30:00Z",
  "web_url": "https://gitlab.com/..."
}

# Get linked issues in JSON
gitlab-manager mr get-issues -m 123 --json
```

## Development

```bash
# Clone repository
git clone https://gitlab.com/your-org/gitlab-manager.git

# Run tests
go test ./...

# Build binary
go build -o gitlab-manager
```

## Contributing

1. Create a merge request with your changes
2. Add appropriate changelog entry
3. Ensure tests pass
4. Request review

## License

MIT License - see LICENSE file for details
