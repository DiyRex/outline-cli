# Outline CLI - Usage Guide

## Table of Contents

1. [Installation](#installation)
2. [Configuration](#configuration)
3. [Auth & Status](#auth--status)
4. [Push & Pull Workflows](#push--pull-workflows)
5. [Diff - Compare Local vs Remote](#diff---compare-local-vs-remote)
6. [Document Management](#document-management)
7. [Collection Management](#collection-management)
8. [Changelog & Release Notes](#changelog--release-notes)
9. [Publish Command](#publish-command)
10. [Revisions & Attachments](#revisions--attachments)
11. [User & Group Management](#user--group-management)
12. [Search](#search)
13. [Backup](#backup)
14. [Interactive TUI](#interactive-tui)
15. [CI/CD Pipeline Usage](#cicd-pipeline-usage)
16. [Global Flags](#global-flags)
17. [Exit Codes](#exit-codes)
18. [API Key Setup](#api-key-setup)
19. [Troubleshooting](#troubleshooting)

---

## Installation

### Build from Source

```bash
cd Services/outline-stage-instance/outline-cli
make build
```

This produces an `outline` binary in the current directory.

### Install System-wide

```bash
make install
```

Copies the binary to `/usr/local/bin/outline`.

---

## Configuration

### Initial Setup

```bash
# Set your Outline instance URL
outline config --url="https://outline-stage.zuselab.dev"

# Set your API key
outline config --api-key="ol_api_your_key_here"
```

### View Current Configuration

```bash
outline config
```

Output:
```
Config file: /Users/you/.outline-cli/config.yaml
URL:         https://outline-stage.zuselab.dev
API Key:     ol_api_EQj...cqKK
```

### Environment Variables

You can also set configuration via environment variables:

```bash
export OUTLINE_URL="https://outline-stage.zuselab.dev"
export OUTLINE_API_KEY="ol_api_your_key_here"
```

Environment variables override the config file.

### Config File Location

`~/.outline-cli/config.yaml`

---

## Auth & Status

### Test Credentials

Validate your API key and display user info:

```bash
outline auth test
```

Output:
```
✓ Authentication successful
  User:  John Doe
  Email: john@example.com
  Role:  admin
  Team:  My Team
```

### Quick Identity Check

```bash
outline auth whoami
```

### Connection Status

```bash
outline status
```

Shows URL, API key (masked), connection status, user, role, team, and collection count.

---

## Push & Pull Workflows

### Pushing Local Content

The `push` command uploads a local folder structure to Outline as a collection with nested documents.

**Folder structure:**
```
my-docs/
  getting-started/
    introduction.md
    installation.md
  api-reference/
    authentication.md
    endpoints.md
  changelog.md
```

**Push command:**
```bash
outline push "My Documentation" ./my-docs/
```

**Result in Outline:**
```
Collection: My Documentation
  Getting Started (folder document)
    Introduction
    Installation
  Api Reference (folder document)
    Authentication
    Endpoints
  Changelog
```

#### How Titles are Determined

1. If the markdown file starts with a `# Heading`, that heading becomes the document title
2. Otherwise, the filename is used (hyphens/underscores become spaces, title-cased)

#### Dry Run

Preview what would be created without making changes:

```bash
outline push --dry-run "My Documentation" ./my-docs/
```

#### Delete Remote Orphans

Remove documents from the remote collection that don't exist locally:

```bash
outline push --delete "My Documentation" ./my-docs/
```

#### Updating Existing Content

Running push again on an existing collection updates documents that match by title at the same hierarchy level:

```bash
# First push creates everything
outline push "My Docs" ./docs/

# Edit local files, then push again to update
outline push "My Docs" ./docs/
```

### Pulling Content

Download a collection to a local folder:

```bash
outline pull "My Documentation" ./output/
```

This creates the folder structure with markdown files matching the document hierarchy.

Documents with children become directories. The parent document content is saved as `index.md` within the directory.

#### Dry Run

```bash
outline pull --dry-run "My Documentation" ./output/
```

---

## Diff - Compare Local vs Remote

Compare a local folder against a remote collection to see what would change:

```bash
outline diff "My Documentation" ./my-docs/
```

Output shows documents in categories:
- **Added** - Local files not yet in remote
- **Modified** - Files that differ from remote
- **Unchanged** - Files matching remote
- **Remote only** - Documents in remote but not in local folder

---

## Document Management

### List Documents

```bash
# List all documents
outline documents list

# List documents in a specific collection
outline documents list --collection <collection-id>

# JSON output
outline documents list -f json
```

### View Document

```bash
outline documents info <document-id>
```

### Create Document

```bash
# From inline text
outline documents create --title "My Document" --text "# Hello\n\nContent here" --collection <id>

# From file
outline documents create --title "My Document" --file ./content.md --collection <id>

# From stdin (for piping in CI/CD)
cat doc.md | outline documents create --title "Piped Doc" --collection <id> --stdin

# From a template
outline documents create --title "From Template" --collection <id> --template <template-id>

# As child of another document
outline documents create --title "Sub Doc" --collection <id> --parent <parent-doc-id>

# As draft (unpublished)
outline documents create --title "Draft" --collection <id> --publish=false
```

### Update Document

```bash
# Update title
outline documents update <id> --title "New Title"

# Update content from file
outline documents update <id> --file ./updated-content.md

# Update from stdin
echo "New content" | outline documents update <id> --stdin

# Update both
outline documents update <id> --title "New Title" --text "New content"
```

### Delete Document

```bash
# Soft delete (moves to trash)
outline documents delete <id>

# Permanent delete (asks for confirmation, skip with --yes)
outline documents delete <id> --permanent
```

### Archive & Restore

```bash
outline documents archive <id>
outline documents restore <id>
```

### Move Document

```bash
# Move to different collection
outline documents move <id> --collection <target-collection-id>

# Move under a parent document
outline documents move <id> --parent <parent-doc-id>
```

### Export Document

```bash
# Print to stdout
outline documents export <id>

# Save to file
outline documents export <id> --output ./exported.md
```

### Duplicate Document

```bash
# Duplicate single document
outline documents duplicate <id>

# Duplicate with all children
outline documents duplicate <id> --recursive
```

### Search Documents

```bash
# Full-text search
outline documents search "search query"

# Search within collection
outline documents search "query" --collection <id>
```

### Drafts & Recently Viewed

```bash
# List your draft documents
outline documents drafts

# List recently viewed documents
outline documents viewed
```

### Unpublish

```bash
outline documents unpublish <id>
```

---

## Collection Management

### List Collections

```bash
outline collections list
```

### View Collection Details

```bash
outline collections info <id>
```

### Create Collection

```bash
outline collections create --name "New Collection"
outline collections create --name "Docs" --description "Documentation" --color "#7C3AED"
```

### Update Collection

```bash
outline collections update <id> --name "Updated Name"
outline collections update <id> --description "New description"
```

### Delete Collection

```bash
outline collections delete <id>
```

Destructive operation - asks for confirmation. Use `--yes` to skip confirmation.

### View Document Tree

```bash
outline collections tree <id>
```

Shows the hierarchical document structure:
```
Getting Started (abc123)
  Introduction (def456)
  Installation (ghi789)
API Reference (jkl012)
  Authentication (mno345)
```

### Archive & Restore

```bash
outline collections archive <id>
outline collections restore <id>
```

---

## Changelog & Release Notes

Generate changelogs from git commit history using conventional commit format.

### Generate Changelog

```bash
# Generate markdown changelog to stdout
outline changelog generate --from v1.0 --to v1.1

# Include commit authors
outline changelog generate --from v1.0 --to v1.1 --include-authors

# Specify repo path
outline changelog generate --from v1.0 --to v1.1 --repo /path/to/repo
```

Output groups commits by type:
```markdown
# Changelog v1.0..v1.1

## Features
- **auth:** add login endpoint (`abc1234`)

## Bug Fixes
- resolve null pointer (`def6789`)

## Documentation
- update README (`ghi1111`)
```

### Push Changelog to Outline

```bash
outline changelog push --from v1.0 --to v1.1 \
  --collection <id> --title "Release v1.1"

# Nest under a parent document
outline changelog push --from v1.0 --to v1.1 \
  --collection <id> --title "Release v1.1" --parent <parent-id>
```

If a document with the same title already exists, it updates the existing document (upsert behavior).

---

## Publish Command

Quick way to upsert a local markdown file to Outline:

```bash
# Publish (creates new or updates existing by title match)
outline publish ./release-notes.md --collection <id>

# Override the document title
outline publish ./doc.md --collection <id> --title "Custom Title"

# Nest under a parent
outline publish ./doc.md --collection <id> --parent <parent-id>
```

Title is extracted from the first `# Heading` in the file, or derived from the filename.

---

## Revisions & Attachments

### Revisions

Browse and manage document revisions:

```bash
# List revisions for a document
outline revisions list --document <document-id>

# Get revision details
outline revisions info <revision-id>

# Delete a revision
outline revisions delete <revision-id>
```

### Attachments

```bash
# List attachments (optionally filter by document)
outline attachments list [--document <document-id>]

# Delete an attachment
outline attachments delete <attachment-id>
```

---

## User & Group Management

### Users

```bash
# List all users
outline users list

# Search users
outline users list --query "john"

# Get current user info
outline users info

# Get specific user
outline users info <user-id>
```

### Groups

```bash
# List groups
outline groups list

# Create group
outline groups create --name "Engineering"

# View members
outline groups members <group-id>

# Add user to group
outline groups add-user <group-id> <user-id>

# Remove user from group
outline groups remove-user <group-id> <user-id>

# Delete group
outline groups delete <group-id>
```

---

## Search

### Full-text Search

```bash
outline search "deployment process"
```

### Title-only Search (faster)

```bash
outline search "deployment" --titles
```

### Scoped Search

```bash
outline search "API" --collection <collection-id>
```

---

## Backup

Download all collections as local markdown folders:

```bash
# Backup to default directory (outline-backup-YYYYMMDD-HHMMSS)
outline backup

# Backup to specific directory
outline backup --output ./my-backup/
```

Each collection becomes a subfolder containing its document hierarchy as markdown files.

---

## Interactive TUI

Launch by running `outline` without arguments:

```bash
outline
```

### Navigation

| Key | Action |
|-----|--------|
| `j` / `↓` | Move down |
| `k` / `↑` | Move up |
| `Enter` | Select / Open |
| `Esc` / `Backspace` | Go back |
| `/` | Open search |
| `Tab` | Open search result |
| `q` | Quit |

### Views

1. **Collections** - Browse all collections
2. **Documents** - Tree view of documents within a collection
3. **Document Detail** - View document content with scrolling
4. **Search** - Interactive search with results

> Note: TUI is disabled in `--non-interactive` mode.

---

## CI/CD Pipeline Usage

The CLI is designed to work in automated pipelines (GitHub Actions, GitLab CI, etc.).

### Pipeline Setup

```bash
# Use environment variables for config
export OUTLINE_URL="https://outline.example.com"
export OUTLINE_API_KEY="$OUTLINE_API_KEY_SECRET"

# Validate credentials first
outline auth test --non-interactive
```

### Example: Publish Release Notes

```bash
outline changelog push \
  --from "$PREV_TAG" --to "$NEW_TAG" \
  --collection <id> --title "Release $NEW_TAG" \
  --non-interactive --quiet
```

### Example: Push Documentation

```bash
outline push "API Docs" ./docs/ --non-interactive --quiet
```

### Example: Pipe Content

```bash
echo "Build $BUILD_ID completed at $(date)" | \
  outline documents create \
    --title "Build Report $BUILD_ID" \
    --collection <id> --stdin --non-interactive
```

### Example: Check Exit Codes

```bash
outline auth test --non-interactive
if [ $? -eq 3 ]; then
  echo "Authentication failed"
  exit 1
fi
```

### Key Flags for CI/CD

| Flag | Purpose |
|------|---------|
| `--non-interactive` | Prevents TUI launch, errors instead of prompting |
| `--quiet` / `-q` | Only shows errors |
| `--yes` / `-y` | Auto-confirms destructive operations |
| `--no-color` | Plain text output (also respects `NO_COLOR` env) |
| `--format json` | Machine-readable output |

---

## Global Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--format` | `-f` | `table` | Output format: `table` or `json` |
| `--quiet` | `-q` | `false` | Suppress non-error output |
| `--verbose` | `-v` | `false` | Enable debug output |
| `--no-color` | | `false` | Disable colored output |
| `--non-interactive` | | `false` | Disable TUI and interactive prompts |
| `--yes` | `-y` | `false` | Auto-confirm destructive operations |
| `--timeout` | | `30` | HTTP timeout in seconds |

---

## Exit Codes

Standardized exit codes for scripting and CI/CD:

| Code | Meaning | Example |
|------|---------|---------|
| 0 | Success | Command completed normally |
| 1 | General error | Unexpected error |
| 2 | Configuration error | Missing URL or API key |
| 3 | Authentication error | Invalid/expired API key (401/403) |
| 4 | Not found | Document/collection not found (404) |
| 5 | Rate limited | Too many requests (429) |
| 6 | Validation error | Invalid input parameters |

---

## API Key Setup

### Creating an API Key in Outline

1. Go to your Outline instance Settings
2. Navigate to **API** section
3. Click **New API Key**
4. Set a name (e.g., "CLI Access")
5. Select scopes:
   - For full access: `read`, `write`
   - For read-only: `read`
6. Set expiration (optional)
7. Copy the generated key (starts with `ol_api_`)

### Required Scopes by Feature

| Feature | Required Scopes |
|---------|----------------|
| List/view | `read` or `documents:read`, `collections:read` |
| Create/update | `write` or `documents:write`, `collections:write` |
| Push/Pull | `read`, `write` |
| User management | `users:read` |
| Group management | `groups:read`, `groups:write` |
| Share management | `shares:read`, `shares:write` |
| Comments | `comments:read`, `comments:write` |

---

## Troubleshooting

### "URL not configured"
Run `outline config --url="https://your-outline-instance.com"`

### "API key not configured"
Run `outline config --api-key="ol_api_your_key"`

### "API error 401"
Your API key is invalid or expired. Create a new one in Outline Settings > API.

### "API error 403"
Your API key lacks the required scopes for the operation. Check the scopes table above.

### "connection refused"
Verify the Outline URL is correct and the instance is running.

### Push creates duplicate documents
Push matches by document title at each hierarchy level. If you renamed a document, it will be created as new rather than updating the old one.

### Pull missing content
Some documents may have empty content (e.g., folder-type documents). These are expected to have no content.

### Non-interactive mode errors
If running in a pipeline and getting "no subcommand specified", make sure you're providing a command (e.g., `outline collections list`, not just `outline`).

### Changelog shows no commits
Ensure the git refs (`--from` and `--to`) exist in the repository. Use `git tag -l` to list available tags.
