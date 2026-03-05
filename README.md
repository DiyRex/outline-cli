# Outline CLI

A comprehensive CLI and TUI tool for managing [Outline](https://www.getoutline.com/) wiki instances. Built with Go, Cobra, and Bubble Tea.

## Features

- **Configuration management** - Store Outline URL and API key
- **Push/Pull sync** - Upload local folders as collections, download collections as markdown files
- **Full API coverage** - Documents, collections, users, groups, comments, shares, stars, search, events, revisions, attachments
- **Interactive TUI** - Browse collections and documents in the terminal
- **CI/CD ready** - Quiet mode, exit codes, stdin piping, non-interactive mode, auto-confirm
- **Changelog generation** - Generate release notes from git history and push to Outline
- **Diff & Backup** - Compare local vs remote, backup all collections
- **Multiple output formats** - Table and JSON output

## Installation

```bash
# Build from source
make build

# Install to /usr/local/bin
make install
```

## Quick Start

```bash
# Configure the CLI
outline config --url="https://outline.example.com" --api-key="ol_api_your_key_here"

# Verify connection
outline status

# List collections
outline collections list

# Push a folder of markdown files as a collection
outline push "My Docs" ./docs/

# Pull a collection to local folder
outline pull "My Docs" ./output/

# Launch interactive TUI
outline
```

## Global Flags

All commands support these flags:

| Flag | Short | Description |
|------|-------|-------------|
| `--format` | `-f` | Output format: `table` or `json` (default: table) |
| `--quiet` | `-q` | Suppress non-error output |
| `--verbose` | `-v` | Enable debug output |
| `--no-color` | | Disable colored output (also respects `NO_COLOR` env) |
| `--non-interactive` | | Disable TUI and interactive prompts |
| `--yes` | `-y` | Auto-confirm destructive operations |
| `--timeout` | | HTTP timeout in seconds (default: 30) |

## Configuration

Configuration is stored in `~/.outline-cli/config.yaml`.

```bash
outline config --url="https://outline.example.com"
outline config --api-key="ol_api_your_key_here"
```

Environment variables are also supported:
- `OUTLINE_URL`
- `OUTLINE_API_KEY`

### API Key Scopes

Create an API key in Outline Settings > API with the following scopes:
- **Full access**: `read`, `write`
- **Granular**: `documents:read`, `documents:write`, `collections:read`, `collections:write`, `users:read`, `groups:read`, `groups:write`, `comments:read`, `comments:write`, `shares:read`, `shares:write`

## Commands

### Auth & Status

```bash
# Test credentials and show user info
outline auth test

# Quick identity check
outline auth whoami

# Connection status, user, team, collection count
outline status
```

### Push & Pull

```bash
# Push a folder as a collection (creates collection if needed)
outline push "Collection Name" ./path/to/folder/

# Push with dry-run (preview only)
outline push --dry-run "Collection Name" ./path/to/folder/

# Push and delete remote docs not in local folder
outline push --delete "Collection Name" ./path/to/folder/

# Pull a collection to local folder
outline pull "Collection Name" ./output/

# Pull with dry-run
outline pull --dry-run "Collection Name" ./output/
```

The push command preserves folder hierarchy:
```
docs/
  getting-started/
    intro.md
    setup.md
  api/
    authentication.md
    endpoints.md
```
Becomes a collection with nested documents matching the structure.

### Diff

Compare local folder against a remote collection:

```bash
outline diff "Collection Name" ./local-folder/
```

Shows added, modified, unchanged, and remote-only documents.

### Documents

```bash
outline documents list [--collection <id>]
outline documents info <id>
outline documents create --title "Title" --text "Content" --collection <id>
outline documents create --title "Title" --file ./content.md --collection <id>
outline documents create --title "Title" --collection <id> --stdin     # read from stdin
outline documents create --title "Title" --collection <id> --template <id>
outline documents update <id> --title "New Title" --text "New content"
outline documents update <id> --file ./content.md
outline documents update <id> --stdin                                  # update from stdin
outline documents delete <id> [--permanent]
outline documents archive <id>
outline documents restore <id>
outline documents move <id> --collection <target-id> [--parent <parent-id>]
outline documents export <id> [--output file.md]
outline documents duplicate <id> [--recursive]
outline documents search <query> [--collection <id>]
outline documents drafts
outline documents viewed
outline documents unpublish <id>
```

#### Stdin Piping (CI/CD)

```bash
cat doc.md | outline documents create --title "Piped Doc" --collection <id> --stdin
echo "Updated content" | outline documents update <id> --stdin
```

### Collections

```bash
outline collections list
outline collections info <id>
outline collections create --name "Name" [--description "Desc"] [--color "#hex"]
outline collections update <id> --name "New Name"
outline collections delete <id>    # requires confirmation (skip with --yes)
outline collections archive <id>
outline collections restore <id>
outline collections tree <id>
```

### Changelog (CI/CD)

Generate release notes from git commits:

```bash
# Generate changelog to stdout
outline changelog generate --from v1.0 --to v1.1

# Include commit authors
outline changelog generate --from v1.0 --to v1.1 --include-authors

# Use a specific repo path
outline changelog generate --from v1.0 --to v1.1 --repo /path/to/repo

# Generate and push to Outline as a document
outline changelog push --from v1.0 --to v1.1 --collection <id> --title "Release v1.1"
outline changelog push --from v1.0 --to v1.1 --collection <id> --title "Release v1.1" --parent <parent-id>
```

Commits are grouped by conventional commit type (feat, fix, docs, chore, etc.).

### Publish

Quick upsert a local markdown file to Outline:

```bash
# Publish a file (creates or updates by title match)
outline publish ./release-notes.md --collection <id>

# Override the title
outline publish ./doc.md --collection <id> --title "Custom Title"

# Nest under a parent document
outline publish ./doc.md --collection <id> --parent <parent-id>
```

### Revisions

```bash
outline revisions list --document <id>
outline revisions info <id>
outline revisions delete <id>
```

### Attachments

```bash
outline attachments list [--document <id>]
outline attachments delete <id>
```

### Users

```bash
outline users list [--query "search"]
outline users info [id]
```

### Groups

```bash
outline groups list
outline groups create --name "Group Name"
outline groups delete <id>
outline groups members <id>
outline groups add-user <group-id> <user-id>
outline groups remove-user <group-id> <user-id>
```

### Comments

```bash
outline comments list [--document <id>]
outline comments create --document <id> --text "Comment text"
outline comments delete <id>
outline comments resolve <id>
```

### Shares

```bash
outline shares list
outline shares create --document <id>
outline shares revoke <id>
```

### Stars

```bash
outline stars list
outline stars create --document <id>
outline stars delete <id>
```

### Search

```bash
outline search "query" [--collection <id>] [--titles]
```

### Events

```bash
outline events [--document <id>] [--collection <id>] [--audit]
```

### Backup

Download all collections as local markdown folders:

```bash
outline backup [--output ./backup-dir/]
```

### Version

```bash
outline version
```

## Exit Codes

Standardized exit codes for CI/CD pipelines:

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Configuration error (missing URL/key) |
| 3 | Authentication error (401/403) |
| 4 | Resource not found (404) |
| 5 | Rate limited (429) |
| 6 | Validation error (bad input) |

## Interactive TUI

Run `outline` without arguments to launch the interactive TUI.

**Navigation:**
- `j`/`k` or arrow keys: Navigate up/down
- `Enter`: Select / drill into
- `Esc` / `Backspace`: Go back
- `/`: Open search
- `q`: Quit

## CI/CD Pipeline Usage

```bash
# In a GitHub Actions / GitLab CI pipeline:
export OUTLINE_URL="https://outline.example.com"
export OUTLINE_API_KEY="$OUTLINE_API_KEY_SECRET"

# Validate credentials
outline auth test --non-interactive

# Generate and publish release notes
outline changelog push --from $PREV_TAG --to $NEW_TAG \
  --collection <id> --title "Release $NEW_TAG" \
  --non-interactive

# Push documentation
outline push "API Docs" ./docs/ --non-interactive --quiet

# Pipe content from CI
echo "Build $BUILD_ID completed" | outline documents create \
  --title "Build Report $BUILD_ID" --collection <id> --stdin --non-interactive
```

## Testing

```bash
# Run all tests
go test ./... -v

# Run tests for a specific package
go test ./internal/api/ -v
go test ./internal/sync/ -v
go test ./internal/changelog/ -v
go test ./internal/cli/ -v
go test ./internal/config/ -v
```

## Development

```bash
make build    # Build binary
make fmt      # Format code
make vet      # Run go vet
make tidy     # Tidy modules
make clean    # Remove binary
```

## License

Internal tool - Zusetech.
