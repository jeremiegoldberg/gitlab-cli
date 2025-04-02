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

## Command Reference

### Merge Requests

```bash
# List merge requests
mpg-gitlab mr list [flags]
  --state string     Filter by state (opened/closed/merged)
  --target string    Filter by target branch
  --author string    Filter by author username
  --labels string    Filter by labels (comma-separated)
  --json            Output in JSON format

# Get merge request details
mpg-gitlab mr get [flags]
  -m, --mr int      Merge request IID (required)
  -p, --project int Project ID
  --json           Output in JSON format

# Get merge request IID from commit message
mpg-gitlab mr get-mr-from-commit [flags]
  -c, --commit string    Commit ID (SHA)
  -m, --message string   Commit message (optional)
  -p, --project int     Project ID

# Create merge request
mpg-gitlab mr create [flags]
  -s, --source string      Source branch (required)
  -t, --target string      Target branch (required)
  --title string          Title for the merge request (required)
  --description string    Description text
  --labels string         Labels to apply (comma-separated)
  --milestone int         Milestone ID to assign
  --assignee string       Assignee username

# Block/Unblock merge request
mpg-gitlab mr block [flags]
  -m, --mr int           Merge request IID (required)
  -r, --reason string    Blocking reason (required)

mpg-gitlab mr unblock [flags]
  -m, --mr int          Merge request IID (required)

# Validate changelog
mpg-gitlab mr check-changelog [flags]
  -m, --mr int          Merge request IID (required)
  --strict             Fail if no changelog entry found
```

### Issues

```bash
# List issues
mpg-gitlab issues list [flags]
  --state string      Filter by state (opened/closed)
  --labels string     Filter by labels (comma-separated)
  --milestone int     Filter by milestone ID
  --json            Output in JSON format

# Get issue details
mpg-gitlab issues get [flags]
  -i, --issue int    Issue IID (required)
  -p, --project int  Project ID
  --json           Output in JSON format

# Create issue
mpg-gitlab issues create [flags]
  --title string         Issue title (required)
  --description string   Issue description
  --labels string        Labels to apply (comma-separated)
  --milestone int        Milestone ID to assign
  --assignee string      Assignee username
```

### Milestones

```bash
# List milestones
mpg-gitlab milestones list [flags]
  -p, --project int   Project ID
  --state string     Filter by state (active/closed)
  --json            Output in JSON format

# Get milestone details
mpg-gitlab milestones get [flags]
  -p, --project int    Project ID
  -m, --milestone int  Milestone ID (required)
  --json             Output in JSON format

# Create milestone
mpg-gitlab milestones create [flags]
  -p, --project int       Project ID
  --title string         Milestone title (required)
  --description string   Description text
  --due-date string      Due date (YYYY-MM-DD)
  --start-date string    Start date (YYYY-MM-DD)

# Update milestone
mpg-gitlab milestones update [flags]
  -p, --project int       Project ID
  -m, --milestone int     Milestone ID (required)
  --title string         New title
  --description string   New description
  --due-date string      New due date
  --state string         New state (activate/close)

# Add changelog to milestone
mpg-gitlab mr add-changelog [flags]
  -p, --project int    Project ID
  -m, --mr int        Merge request IID (required)

Note: Requires the merge request to have a milestone assigned
```

### Notes and Comments

```bash
# List notes
mpg-gitlab notes list [flags]
  -p, --project int   Project ID
  -m, --mr int       Merge request IID
  -i, --issue int    Issue IID
  --json           Output in JSON format

# Add note
mpg-gitlab notes add [flags]
  -p, --project int    Project ID
  -m, --mr int        Merge request IID
  -i, --issue int     Issue IID
  --body string      Note content (required)
```

### Global Flags

Available for all commands:

```bash
  -h, --help        Show help for command
  --debug          Enable debug output
  --quiet          Suppress all output except errors
  --config string  Config file (default is $HOME/.mpg-gitlab.yaml)
```

### Environment Variables

- `GITLAB_TOKEN`: Personal access token for GitLab API
- `GITLAB_API_URL`: Custom GitLab API URL (for self-hosted instances)
- `CI_PROJECT_ID`: Project ID (automatically set in GitLab CI)
- `CI_MERGE_REQUEST_IID`: Merge request IID (automatically set in GitLab CI)

## Examples

### Typical Workflow

```bash
# Create a merge request
mpg-gitlab mr create \
  --source feature-branch \
  --target main \
  --title "Add new feature" \
  --description "[Feature] Implement new authentication system"

# Check changelog
mpg-gitlab mr check-changelog -m 123

# Block MR if needed
mpg-gitlab mr block -m 123 -r "Needs security review"

# Unblock when ready
mpg-gitlab mr unblock -m 123

# Add to milestone changelog
mpg-gitlab milestones add-changelog -m 45
```

### CI/CD Integration

```yaml
validate_merge_request:
  script:
    - |
      if ! mpg-gitlab mr check-changelog -m $CI_MERGE_REQUEST_IID; then
        echo "Missing changelog entry"
        mpg-gitlab mr block -m $CI_MERGE_REQUEST_IID -r "Missing changelog entry"
        exit 1
      fi
     # Check milestone
    - |
      if ! mpg-gitlab mr check-milestone -m $CI_MERGE_REQUEST_IID; then
        echo "Missing milestone"
        mpg-gitlab mr block -m $CI_MERGE_REQUEST_IID -r "Missing milestone assignment"
        exit 1
      fi
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
