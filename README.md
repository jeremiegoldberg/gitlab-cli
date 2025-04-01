# mpg-gitlab

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
go install mpg-gitlab@latest
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
mpg-gitlab mr list
mpg-gitlab mr list --state opened --target main

# Get merge request details
mpg-gitlab mr get -m 123 [--json]

# Create merge request
mpg-gitlab mr create \
  --source feature-branch \
  --target main \
  --title "Add new feature" \
  --description "Implements..."

# Get linked issues
mpg-gitlab mr get-issues -m 123 [--json]

# Get MR description
mpg-gitlab mr get-description -m 123
```

### Milestone Management

```bash
# List milestones
mpg-gitlab milestones list -p PROJECT_ID
mpg-gitlab milestones list --state active

# Get milestone details
mpg-gitlab milestones get -p PROJECT_ID -m MILESTONE_ID [--json]

# Create milestone
mpg-gitlab milestones create \
  -p PROJECT_ID \
  --title "Release 1.0" \
  --description "First major release" \
  --due-date "2024-03-01"

# Update milestone
mpg-gitlab milestones update \
  -p PROJECT_ID \
  -m MILESTONE_ID \
  --title "Release 1.1" \
  --state "closed"

# Delete milestone
mpg-gitlab milestones delete -p PROJECT_ID -m MILESTONE_ID

# Add changelog entries to milestone
mpg-gitlab milestones add-changelog -p PROJECT_ID -m MILESTONE_ID
```

Milestone features include:
- Create, read, update, and delete milestones
- List milestones with filtering options
- Set due dates and start dates
- Track milestone progress
- Automatic changelog collection from merge requests
- JSON output support for automation

Example JSON output:
```json
{
  "id": 1,
  "title": "Release 1.0",
  "description": "First major release",
  "state": "active",
  "due_date": "2024-03-01",
  "start_date": "2024-01-15",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z",
  "web_url": "https://gitlab.com/..."
}
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
mpg-gitlab mr check-changelog -m 123

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
mpg-gitlab mr block -m 123 -r "Needs security review"

# Unblock when ready
mpg-gitlab mr unblock -m 123
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
      if ! mpg-gitlab mr check-changelog -m $CI_MERGE_REQUEST_IID; then
        echo "Missing changelog entry"
        mpg-gitlab mr block -m $CI_MERGE_REQUEST_IID -r "Missing changelog entry"
        exit 1
      fi

    # Other validations
    - mpg-gitlab mr get-issues -m $CI_MERGE_REQUEST_IID
    
  rules:
    - if: $CI_PIPELINE_SOURCE == "merge_request_event"
```

### JSON Output

Many commands support JSON output for integration with other tools:

```bash
# Get MR details in JSON
mpg-gitlab mr get -m 123 --json

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
mpg-gitlab mr get-issues -m 123 --json
```

## Development

```bash
# Clone repository
git clone https://github.com/username/mpg-gitlab.git

# Run tests
go test ./...

# Build binary
go build -o mpg-gitlab
```

## Contributing

1. Create a merge request with your changes
2. Add appropriate changelog entry
3. Ensure tests pass
4. Request review

## License

MIT License - see LICENSE file for details
