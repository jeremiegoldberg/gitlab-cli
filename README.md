# gitlab-cli

A command-line interface tool for managing GitLab resources (issues, merge requests, and milestones) with built-in GitLab CI support.

## Features

- Full CRUD operations for issues, merge requests, and milestones
- Automatic issue linking in merge requests
- GitLab CI integration
- Structured data output (both human-readable and JSON formats)
- Detailed description viewing

## Installation

```bash
# Clone the repository
git clone https://gitlab.com/your-username/gitlab-cli.git

# Build the binary
cd gitlab-cli
go build -o gitlab-cli
```

## Authentication

The tool requires authentication with GitLab. Set one of these environment variables:

```bash
# For GitLab CI (automatically used in GitLab CI/CD pipelines)
export CI_JOB_TOKEN="your-ci-job-token"

# For local development (Personal Access Token)
export GITLAB_TOKEN="your-personal-access-token"
```

## Usage

### Resource Descriptions

You can view the full description of any resource:

```bash
# Get issue description
gitlab-cli issues get-description -i 123

# Get merge request description
gitlab-cli mr get-description -m 456
```

### Issue Linking

The tool automatically detects issue references in merge request descriptions. Supported formats:

- Simple reference: `#123`
- Action keywords: `fixes #123`, `closes #456`, `resolves #789`
- Other references: `references #101`, `refs #202`, `see #303`

```bash
# Create MR with linked issues
gitlab-cli mr create \
  --source feature \
  --target main \
  --title "New Feature" \
  --description "This MR fixes #123 and closes #456"

# View linked issues
gitlab-cli mr get-issues -m 789
gitlab-cli mr get-issues -m 789 --json
```

### Issues

```bash
# List issues
gitlab-cli issues list
gitlab-cli issues list --state opened
gitlab-cli issues list --project 123

# Get issue details
gitlab-cli issues get -i 123
gitlab-cli issues get -i 123 --json

# Get issue description
gitlab-cli issues get-description -i 123

# Create issue
gitlab-cli issues create \
  --title "New Issue" \
  --description "Issue description" \
  --project 123 \
  --labels "bug,urgent"

# Update issue
gitlab-cli issues update \
  --issue 123 \
  --title "Updated Title" \
  --state "close"

# Delete issue
gitlab-cli issues delete --issue 123
```

### Merge Requests

```bash
# List merge requests
gitlab-cli mr list
gitlab-cli mr list --state opened --target main

# Get MR details
gitlab-cli mr get -m 123

# Get MR description
gitlab-cli mr get-description -m 123

# Get issues linked to an MR
gitlab-cli mr get-issues -m 123
gitlab-cli mr get-issues -m 123 --json

# Create MR
gitlab-cli mr create \
  --source feature-branch \
  --target main \
  --title "New Feature" \
  --description "Feature description" \
  --remove-source

# Create MR with linked issues
gitlab-cli mr create \
  --source feature-branch \
  --target main \
  --title "New Feature" \
  --description "This MR fixes #123 and closes #456"

# Update MR
gitlab-cli mr update \
  --mr 123 \
  --title "Updated Title" \
  --target develop

# Merge MR
gitlab-cli mr merge \
  --mr 123 \
  --message "Merge commit message"

# Close MR
gitlab-cli mr close --mr 123
```

### Milestones

```bash
# List milestones
gitlab-cli milestones list
gitlab-cli milestones list --state active

# Get milestone details
gitlab-cli milestones get -m 123

# Create milestone
gitlab-cli milestones create \
  --title "v1.0.0" \
  --description "First release" \
  --due-date "2024-12-31"

# Update milestone
gitlab-cli milestones update \
  --milestone 123 \
  --title "v1.0.1" \
  --state "close"

# Delete milestone
gitlab-cli milestones delete --milestone 123
```

## GitLab CI Integration

The tool automatically detects when it's running in GitLab CI and uses the appropriate configuration. Example `.gitlab-ci.yml`:

```yaml
create_issue:
  script:
    - gitlab-cli issues create \
        --title "Pipeline Issue" \
        --description "Created from pipeline"
        # Project ID is automatically detected in CI

list_mrs:
  script:
    - gitlab-cli mr list --state opened
```

### CI Features

- Automatic authentication using `CI_JOB_TOKEN`
- Project ID detection using `CI_PROJECT_ID`
- CI metadata added to created issues
- GitLab API URL configuration using `CI_API_V4_URL`

## Command Reference

### Global Flags

- `--project, -p`: Project ID (optional in CI environment)

### Issues Commands

```bash
gitlab-cli issues [command] [flags]
```

Commands:
- `list`: List issues
  - `--state`: Filter by state (opened/closed)
  - `--project`: Project ID
  - `--json, -j`: Output in JSON format
- `get`: Get issue details
  - `--issue, -i`: Issue IID (required)
  - `--project`: Project ID
  - `--json, -j`: Output in JSON format
- `get-description`: Get issue description
  - `--issue, -i`: Issue IID (required)
- `create`: Create new issue
  - `--title, -t`: Issue title (required)
  - `--description, -d`: Issue description
  - `--labels, -l`: Comma-separated labels
  - `--project`: Project ID
- `update`: Update issue
  - `--issue, -i`: Issue IID (required)
  - `--title`: New title
  - `--description`: New description
  - `--state`: New state (close/reopen)
- `delete`: Delete issue
  - `--issue, -i`: Issue IID (required)

### Merge Requests Commands

```bash
gitlab-cli mr [command] [flags]
```

Commands:
- `list`: List merge requests
  - `--state`: Filter by state (opened/closed/merged/all)
  - `--target`: Filter by target branch
  - `--json, -j`: Output in JSON format
- `get`: Get MR details
  - `--mr, -m`: MR IID (required)
  - `--json, -j`: Output in JSON format
- `get-description`: Get MR description
  - `--mr, -m`: MR IID (required)
- `get-issues`: Get issues linked to an MR
  - `--mr, -m`: MR IID (required)
  - `--json`: Get JSON output
- `create`: Create new MR
  - `--source, -s`: Source branch (required)
  - `--target, -t`: Target branch (default: main)
  - `--title, -T`: MR title (required)
  - `--description, -d`: MR description
  - `--remove-source, -r`: Remove source branch when merged
- `merge`: Merge an MR
  - `--mr, -m`: MR IID (required)
  - `--message, -M`: Merge commit message
- `close`: Close an MR
  - `--mr, -m`: MR IID (required)

### Milestones Commands

```bash
gitlab-cli milestones [command] [flags]
```

Commands:
- `list`: List milestones
  - `--state`: Filter by state (active/closed)
  - `--json, -j`: Output in JSON format
- `get`: Get milestone details
  - `--milestone, -m`: Milestone ID (required)
  - `--json, -j`: Output in JSON format
- `create`: Create milestone
  - `--title, -t`: Milestone title (required)
  - `--description, -d`: Description
  - `--due-date, -D`: Due date (YYYY-MM-DD)
- `update`: Update milestone
  - `--milestone, -m`: Milestone ID (required)
  - `--title`: New title
  - `--description`: New description
  - `--due-date`: New due date
  - `--state`: State event (close/activate)

## JSON Output

Most read operations support JSON output format. This is useful for scripting or when you need to process the data programmatically.

```bash
# Get issue in JSON format
gitlab-cli issues get -i 123 --json

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
gitlab-cli mr get-issues -m 456 --json

# List all open issues in JSON
gitlab-cli issues list --state opened --json
```

The JSON output includes all relevant fields from the GitLab API, making it easy to integrate with other tools or scripts.
