<p align="center">
  <img src="docs/outline-icon.svg" width="80" alt="Outline" />
</p>

<h1 align="center">Outline CLI</h1>

<p align="center">
  <strong>A powerful CLI &amp; TUI for managing <a href="https://www.getoutline.com/">Outline</a> wiki instances.</strong><br/>
  Built with Go &bull; Cobra &bull; Bubble Tea
</p>

<p align="center">
  <a href="#-install">Install</a> &bull;
  <a href="#-quick-start">Quick Start</a> &bull;
  <a href="#-commands">Commands</a> &bull;
  <a href="#-cicd-pipeline-usage">CI/CD</a> &bull;
  <a href="#-interactive-tui">TUI</a> &bull;
  <a href="docs/usage.md">Full Docs</a>
</p>

---

## ✦ Features

| | Feature | Description |
|---|---------|-------------|
| 📄 | **Full API Coverage** | Documents, collections, users, groups, comments, shares, stars, search, events, revisions, attachments |
| 🔄 | **Push / Pull Sync** | Upload local folders as collections, download collections as markdown |
| 🖥️ | **Interactive TUI** | Browse and search your wiki from the terminal |
| 🚀 | **CI/CD Ready** | Quiet mode, exit codes, stdin piping, non-interactive mode, auto-confirm |
| 📋 | **Changelog** | Generate release notes from git history and push to Outline |
| 🔍 | **Diff** | Compare local folder vs remote collection |
| 💾 | **Backup** | Download all collections as local markdown |

---

## 📦 Install

### Quick Install (Linux / macOS)

```bash
curl -fsSL https://raw.githubusercontent.com/DiyRex/outline-cli/main/install.sh | sh
```

> Auto-detects your OS and architecture, downloads the latest release binary, and installs to `/usr/local/bin`.

### Download from Releases

Pre-built binaries for every platform are available on the [Releases page](https://github.com/DiyRex/outline-cli/releases):

| Platform | Binary |
|----------|--------|
| Linux (x86_64) | `outline-linux-amd64` |
| Linux (ARM64) | `outline-linux-arm64` |
| macOS (Intel) | `outline-darwin-amd64` |
| macOS (Apple Silicon) | `outline-darwin-arm64` |
| Windows (x86_64) | `outline-windows-amd64.exe` |

### Build from Source

```bash
git clone https://github.com/DiyRex/outline-cli.git
cd outline-cli
make build    # produces ./outline binary
make install  # copies to /usr/local/bin
```

---

## 🚀 Quick Start

```bash
# 1. Configure
outline config --url="https://outline.example.com"
outline config --api-key="ol_api_your_key_here"

# 2. Verify
outline status

# 3. Use
outline collections list
outline push "My Docs" ./docs/
outline pull "My Docs" ./output/

# 4. Interactive mode
outline
```

---

## 📖 Commands

### Auth & Status

```bash
outline auth test          # Validate credentials
outline auth whoami        # Show current user
outline status             # Connection status, user, team info
```

### Push & Pull

```bash
outline push "Docs" ./folder/              # Upload folder as collection
outline push "Docs" ./folder/ --dry-run    # Preview changes
outline push "Docs" ./folder/ --delete     # Remove remote orphans
outline pull "Docs" ./output/              # Download collection
outline pull "Docs" ./output/ --dry-run    # Preview download
```

**Folder hierarchy is preserved:**
```
docs/
├── getting-started/
│   ├── intro.md
│   └── setup.md
└── api/
    ├── auth.md
    └── endpoints.md
```

### Diff

```bash
outline diff "Docs" ./folder/    # Compare local vs remote
```

### Documents

```bash
outline documents list [--collection <id>]
outline documents info <id>
outline documents create --title "Title" --collection <id> --text "Content"
outline documents create --title "Title" --collection <id> --file ./doc.md
outline documents create --title "Title" --collection <id> --stdin        # pipe from stdin
outline documents create --title "Title" --collection <id> --template <id>
outline documents update <id> --title "New Title" --text "Updated"
outline documents update <id> --file ./doc.md
outline documents update <id> --stdin
outline documents delete <id> [--permanent]
outline documents archive <id>
outline documents restore <id>
outline documents move <id> --collection <id> [--parent <id>]
outline documents export <id> [--output file.md]
outline documents duplicate <id> [--recursive]
outline documents search <query> [--collection <id>]
outline documents drafts
outline documents viewed
outline documents unpublish <id>
```

### Collections

```bash
outline collections list
outline collections info <id>
outline collections create --name "Name" [--description "Desc"] [--color "#hex"]
outline collections update <id> --name "New Name"
outline collections delete <id>
outline collections archive <id>
outline collections restore <id>
outline collections tree <id>
```

### Changelog & Release Notes

```bash
# Generate changelog from git commits (conventional commits)
outline changelog generate --from v1.0 --to v1.1 [--include-authors] [--repo /path]

# Generate and push to Outline
outline changelog push --from v1.0 --to v1.1 --collection <id> --title "Release v1.1"
```

### Publish

```bash
# Upsert a local markdown file (creates or updates by title match)
outline publish ./doc.md --collection <id> [--title "Override"] [--parent <id>]
```

### Revisions & Attachments

```bash
outline revisions list --document <id>
outline revisions info <id>
outline revisions delete <id>

outline attachments list [--document <id>]
outline attachments delete <id>
```

### Users & Groups

```bash
outline users list [--query "search"]
outline users info [id]

outline groups list
outline groups create --name "Name"
outline groups delete <id>
outline groups members <id>
outline groups add-user <group-id> <user-id>
outline groups remove-user <group-id> <user-id>
```

### Comments, Shares, Stars

```bash
outline comments list [--document <id>]
outline comments create --document <id> --text "Text"
outline comments delete <id>
outline comments resolve <id>

outline shares list
outline shares create --document <id>
outline shares revoke <id>

outline stars list
outline stars create --document <id>
outline stars delete <id>
```

### Search & Events

```bash
outline search "query" [--collection <id>] [--titles]
outline events [--document <id>] [--collection <id>] [--audit]
```

### Backup

```bash
outline backup [--output ./backup-dir/]
```

---

## ⚙️ Global Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--format` | `-f` | `table` | Output format: `table` or `json` |
| `--quiet` | `-q` | `false` | Suppress non-error output |
| `--verbose` | `-v` | `false` | Enable debug output |
| `--no-color` | | `false` | Disable colored output (respects `NO_COLOR` env) |
| `--non-interactive` | | `false` | Disable TUI and interactive prompts |
| `--yes` | `-y` | `false` | Auto-confirm destructive operations |
| `--timeout` | | `30` | HTTP timeout in seconds |

---

## 🔢 Exit Codes

Standardized for scripting and CI/CD:

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | General error |
| `2` | Configuration error |
| `3` | Authentication error (401/403) |
| `4` | Not found (404) |
| `5` | Rate limited (429) |
| `6` | Validation error |

---

## 🚀 CI/CD Pipeline Usage

```bash
# Setup
export OUTLINE_URL="https://outline.example.com"
export OUTLINE_API_KEY="$OUTLINE_API_KEY_SECRET"

# Validate credentials
outline auth test --non-interactive

# Push documentation
outline push "API Docs" ./docs/ --non-interactive --quiet

# Generate & publish release notes
outline changelog push --from "$PREV_TAG" --to "$NEW_TAG" \
  --collection <id> --title "Release $NEW_TAG" --non-interactive

# Pipe content
echo "Build $BUILD_ID completed" | outline documents create \
  --title "Build Report" --collection <id> --stdin --non-interactive

# Check exit code
outline auth test --non-interactive
if [ $? -eq 3 ]; then echo "Auth failed"; exit 1; fi
```

---

## 🖥️ Interactive TUI

Run `outline` without arguments to launch the TUI.

| Key | Action |
|-----|--------|
| `j` / `↓` | Move down |
| `k` / `↑` | Move up |
| `Enter` | Select / Open |
| `Esc` / `Backspace` | Go back |
| `/` | Search |
| `q` | Quit |

**Views:** Collections → Documents → Document Detail → Search

---

## 🔧 Configuration

```bash
outline config --url="https://your-instance.com"
outline config --api-key="ol_api_your_key"
outline config    # view current config
```

Config file: `~/.outline-cli/config.yaml`

Environment variables (`OUTLINE_URL`, `OUTLINE_API_KEY`) override the config file.

### API Key Scopes

Create an API key in **Outline Settings → API**:

| Feature | Required Scopes |
|---------|----------------|
| Read operations | `read` |
| Write operations | `write` |
| Push / Pull | `read`, `write` |
| User management | `users:read` |
| Group management | `groups:read`, `groups:write` |

---

## 🧪 Testing

```bash
go test ./... -v              # Run all tests
go test ./internal/api/ -v    # API client tests
go test ./internal/sync/ -v   # Sync engine tests
go test ./internal/cli/ -v    # CLI output & exit code tests
```

---

## 🛠️ Development

```bash
make build    # Build binary
make install  # Install to /usr/local/bin
make fmt      # Format code
make vet      # Run go vet
make tidy     # Tidy modules
make clean    # Remove binary
```

---

## 📚 Documentation

- [Usage Guide](docs/usage.md) — Detailed walkthrough of every feature
- [API Reference](docs/api-reference.md) — Endpoint mapping and client architecture

---

## 📄 License

Internal tool — [Zusetech](https://zusetech.com).
