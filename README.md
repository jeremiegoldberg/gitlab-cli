# gitlab-manager

A command-line interface tool for managing GitLab resources (issues, merge requests, and milestones) with built-in GitLab CI support.

## Features

- Full CRUD operations for issues, merge requests, and milestones
- Automatic issue linking in merge requests
- Changelog entry validation
- Merge request blocking/unblocking
- GitLab CI integration
- Structured data output (both human-readable and JSON formats)
- Detailed description viewing

## Installation

```bash
# Clone the repository
git clone https://gitlab.com/your-username/gitlab-manager.git

# Build the binary
cd gitlab-manager
go build -o gitlab-manager
```

## Authentication

The tool supports two authentication methods:
1. GitLab Personal Access Token (for local use)
2. GitLab CI Job Token (for CI/CD pipelines)

```bash
# Using personal access token
export GITLAB_TOKEN=your_personal_access_token

# In GitLab CI, CI_JOB_TOKEN is automatically available
```

## Usage

### Merge Requests

```bash
# List merge requests
gitlab-manager mr list
gitlab-manager mr list --state opened
gitlab-manager mr list --target main

# Get merge request details
gitlab-manager mr get -m 123

# Create merge request
gitlab-manager mr create \
  --source feature-branch \
  --target main \
  --title "Add new feature" \
  --description "Implements..." \
  --remove-source

# Update merge request
gitlab-manager mr update -m 123 --title "Updated title"

# Merge a request
gitlab-manager mr merge -m 123 --message "Custom merge commit message"

# Close a request
gitlab-manager mr close -m 123

# Get linked issues
gitlab-manager mr get-issues -m 123 [--json]

# Check changelog entries
gitlab-manager mr check-changelog -m 123

# Block/Unblock merge requests
gitlab-manager mr block -m 123 -r "Needs security review"
gitlab-manager mr unblock -m 123
```

### Issues

```bash
# List issues
gitlab-manager issues list
gitlab-manager issues list --state opened

# Get issue details
gitlab-manager issues get -i 123 [--json]

# Create issue
gitlab-manager issues create \
  --title "Bug report" \
  --description "Found a bug..." \
  --labels bug,priority::high

# Update issue
gitlab-manager issues update -i 123 --state close

# Delete issue
gitlab-manager issues delete -i 123
```

### Milestones

```bash
# List milestones
gitlab-manager milestones list
gitlab-manager milestones list --state active

# Get milestone details
gitlab-manager milestones get -m 123 [--json]

# Create milestone
gitlab-manager milestones create \
  --title "v1.0.0" \
  --description "First release" \
  --due-date 2024-12-31

# Update milestone
gitlab-manager milestones update -m 123 --state close

# Delete milestone
gitlab-manager milestones delete -m 123
```

## Changelog Validation

The tool enforces changelog entries in merge requests or their linked issues. Valid changelog entries must start with one of:
- [Feature]
- [Improvement]
- [Fix]
- [Infra]
- [No-Changelog-Entry]

Example usage in CI:
```yaml
validate_mr:
  script:
    - gitlab-manager mr check-changelog  # Will use CI_MERGE_REQUEST_IID
    - gitlab-manager mr get-issues       # Will use CI_MERGE_REQUEST_IID
```

## Merge Request Blocking

You can block merge requests to prevent them from being merged:

```bash
# Block a merge request
gitlab-manager mr block -m 123 -r "Needs security review"

# Unblock when ready
gitlab-manager mr unblock -m 123
```

When a merge request is blocked:
1. Its title is prefixed with [BLOCKED]
2. A note is added with the blocking reason
3. An unblock note is added when unblocked

## CI/CD Integration

The tool automatically detects when running in GitLab CI and:
1. Uses CI_JOB_TOKEN for authentication
2. Uses CI_PROJECT_ID for project context
3. Adds CI metadata to created resources
4. Can block merges based on validation rules

Example CI configuration:
```yaml
validate_mr:
  script:
    - gitlab-manager mr check-changelog  # Will use CI_MERGE_REQUEST_IID
    - gitlab-manager mr get-issues       # Will use CI_MERGE_REQUEST_IID
```

## JSON Output

Many commands support JSON output for integration with other tools:

```bash
# Get issue in JSON format
gitlab-manager issues get -i 123 --json

# Example output:
{
  "iid": 123,
  "title": "Bug in login form",
  "description": "Users cannot log in when...",
  "state": "opened",
  "labels": ["bug", "priority::high"],
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-16T15:45:00Z",
  "web_url": "https://gitlab.com/group/project/-/issues/123"
}

# Get merge request with linked issues
gitlab-manager mr get-issues -m 456 --json

# List all open issues in JSON
gitlab-manager issues list --state opened --json
```

The JSON output includes all relevant fields from the GitLab API, making it easy to integrate with other tools or scripts.
