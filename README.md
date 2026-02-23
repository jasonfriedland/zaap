# zaap

A CLI utility to delete macOS applications and identify/remove associated files outside the application directory.

## Installation

```bash
go install github.com/jasonfriedland/zaap@latest
```

## Usage

```bash
# List all installed applications
zaap --list

# Delete a specific application
zaap --delete "App Name"

# Dry run (show what would be deleted without actually deleting)
zaap --delete "App Name" --dry-run
```

## Features

- Lists all `.app` bundles in `/Applications`
- Scans common locations for associated files:
  - `~/Library/Preferences/`
  - `~/Library/Application Support/`
  - `~/Library/Caches/`
  - `~/Library/Logs/`
  - `~/Library/Saved Application State/`
  - `~/Library/Containers/`
- Interactive confirmation before deleting
- Handles permission errors gracefully
