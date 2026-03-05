# Outline CLI - API Reference

This document covers all API endpoints supported by the Outline CLI.

## Outline API Overview

- **Base URL**: `https://<your-instance>/api/`
- **Method**: All endpoints use HTTP POST (RPC-style)
- **Auth**: `Authorization: Bearer <api_key>` header
- **Content-Type**: `application/json`
- **API Key Prefix**: `ol_api_`

## Response Format

```json
{
  "ok": true,
  "data": { ... },
  "pagination": {
    "offset": 0,
    "limit": 25
  }
}
```

## Pagination

All list endpoints accept:
- `offset` (int, default 0)
- `limit` (int, default 25, max 100)

The CLI defaults to `limit: 100` to fetch more results per request.

---

## Supported Endpoints

### Documents

| CLI Command | API Endpoint | Description |
|-------------|-------------|-------------|
| `documents list` | `documents.list` | List documents |
| `documents info` | `documents.info` | Get document details |
| `documents create` | `documents.create` | Create document |
| `documents update` | `documents.update` | Update document |
| `documents delete` | `documents.delete` | Delete document |
| `documents archive` | `documents.archive` | Archive document |
| `documents restore` | `documents.restore` | Restore document |
| `documents move` | `documents.move` | Move document |
| `documents export` | `documents.info` | Export content |
| `documents duplicate` | `documents.duplicate` | Duplicate document |
| `documents search` | `documents.search` | Search documents |
| `documents drafts` | `documents.drafts` | List draft documents |
| `documents viewed` | `documents.viewed` | Recently viewed documents |
| `documents unpublish` | `documents.unpublish` | Unpublish a document |

### Collections

| CLI Command | API Endpoint | Description |
|-------------|-------------|-------------|
| `collections list` | `collections.list` | List collections |
| `collections info` | `collections.info` | Get collection details |
| `collections create` | `collections.create` | Create collection |
| `collections update` | `collections.update` | Update collection |
| `collections delete` | `collections.delete` | Delete collection |
| `collections archive` | `collections.archive` | Archive collection |
| `collections restore` | `collections.restore` | Restore collection |
| `collections tree` | `collections.documents` | Get document tree |

### Users

| CLI Command | API Endpoint | Description |
|-------------|-------------|-------------|
| `users list` | `users.list` | List users |
| `users info` | `users.info` | Get user details |

### Groups

| CLI Command | API Endpoint | Description |
|-------------|-------------|-------------|
| `groups list` | `groups.list` | List groups |
| `groups create` | `groups.create` | Create group |
| `groups delete` | `groups.delete` | Delete group |
| `groups members` | `groups.memberships` | List members |
| `groups add-user` | `groups.add_user` | Add user to group |
| `groups remove-user` | `groups.remove_user` | Remove user from group |

### Comments

| CLI Command | API Endpoint | Description |
|-------------|-------------|-------------|
| `comments list` | `comments.list` | List comments |
| `comments create` | `comments.create` | Create comment |
| `comments delete` | `comments.delete` | Delete comment |
| `comments resolve` | `comments.resolve` | Resolve thread |

### Shares

| CLI Command | API Endpoint | Description |
|-------------|-------------|-------------|
| `shares list` | `shares.list` | List shares |
| `shares create` | `shares.create` | Create share |
| `shares revoke` | `shares.revoke` | Revoke share |

### Stars

| CLI Command | API Endpoint | Description |
|-------------|-------------|-------------|
| `stars list` | `stars.list` | List starred items |
| `stars create` | `stars.create` | Star item |
| `stars delete` | `stars.delete` | Unstar item |

### Search

| CLI Command | API Endpoint | Description |
|-------------|-------------|-------------|
| `search` | `documents.search` | Full-text search |
| `search --titles` | `documents.search_titles` | Title search |

### Events

| CLI Command | API Endpoint | Description |
|-------------|-------------|-------------|
| `events` | `events.list` | List events |

### Revisions

| CLI Command | API Endpoint | Description |
|-------------|-------------|-------------|
| `revisions list` | `revisions.list` | List document revisions |
| `revisions info` | `revisions.info` | Get revision details |
| `revisions delete` | `revisions.delete` | Delete a revision |

### Attachments

| CLI Command | API Endpoint | Description |
|-------------|-------------|-------------|
| `attachments list` | `attachments.list` | List attachments |
| `attachments delete` | `attachments.delete` | Delete attachment |

### Auth

| CLI Command | API Endpoint | Description |
|-------------|-------------|-------------|
| `auth test` | `auth.info` | Validate credentials |
| `auth whoami` | `auth.info` | Show current user |
| `status` | `auth.info` + `collections.list` | Connection status |

---

## API Client Architecture

The CLI uses an internal Go API client (`internal/api/`) with service-based organization:

```
Client
  ├── Documents   (documents.go)
  ├── Collections (collections.go)
  ├── Users       (users.go)
  ├── Groups      (groups.go)
  ├── Comments    (comments.go)
  ├── Shares      (shares.go)
  ├── Stars       (stars.go)
  ├── Events      (events.go)
  ├── Search      (search.go)
  ├── Attachments (attachments.go)
  ├── Revisions   (revisions.go)
  └── Auth        (auth.go)
```

Each service wraps the base `Client.Post()` method which handles:
- JSON serialization/deserialization
- Bearer token authentication
- Error response handling
- Configurable HTTP timeout (default 30s, override with `--timeout`)

## Exit Code Mapping

The `ExitCodeFromAPIError()` function maps API errors to exit codes:

| Error Pattern | Exit Code |
|--------------|-----------|
| 401, 403, "authentication" | 3 (Auth error) |
| 404, "not_found" | 4 (Not found) |
| 429, "rate" | 5 (Rate limited) |
| "not configured" | 2 (Config error) |
| Other errors | 1 (General error) |
