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

Example session:

```bash
$ zaap

Select an application to delete:
---------------------------------
1. Dropbox
2. Google Chrome
3. Keynote
4. Maccy
5. Numbers
6. Pages
7. Safari
8. iMovie
0. Exit

Enter number: 1

Selected: Dropbox
Location: /Applications/Dropbox.app

Associated files:
  - /Users/test/Library/Preferences/com.getdropbox.dropbox.plist
  - /Users/test/Library/Application Support/Dropbox
  - /Users/test/Library/Preferences/com.getdropbox.dropbox.plist
  - /Users/test/Library/Caches/DropboxUpdateClient
  - /Users/test/Library/Caches/com.dropbox.DropboxUpdater
  - /Users/test/Library/Caches/com.getdropbox.DropboxMetaInstaller

Startup Items:
  - /Users/test/Library/LaunchAgents/com.dropbox.DropboxUpdater.wake.plist
  - /Users/test/Library/LaunchAgents/com.dropbox.dropboxmacupdate.agent.plist
  - /Users/test/Library/LaunchAgents/com.dropbox.dropboxmacupdate.xpcservice.plist

QuickLook Plugins:
  - /Users/test/Library/QuickLook/DropboxQL.qlgenerator

Delete this application? (y/n): n
Cancelled.
```
