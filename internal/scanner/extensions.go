package scanner

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/jasonfriedland/zaap/internal/catalog"
)

// ControlPanels scans preference pane locations.
type ControlPanels struct {
	Config Config
}

// Scan finds preference panes matching the app name.
func (s ControlPanels) Scan(ctx context.Context, app catalog.App) Iterator {
	return scanNameContains(app, CategoryControlPane, []string{
		filepath.Join(s.Config.HomeDir, "Library/PreferencePanes"),
		filepath.Join(s.Config.SystemLibraryDir, "PreferencePanes"),
	}, "application name")
}

// StartupItems scans launch agent and daemon locations.
type StartupItems struct {
	Config Config
}

// Scan finds startup items matching the app name or bundle ID.
func (s StartupItems) Scan(ctx context.Context, app catalog.App) Iterator {
	bundleID := identity(app)
	items := []Item{}
	for _, dir := range []string{
		filepath.Join(s.Config.HomeDir, "Library/LaunchAgents"),
		filepath.Join(s.Config.SystemLibraryDir, "LaunchAgents"),
		filepath.Join(s.Config.HomeDir, "Library/LaunchDaemons"),
		filepath.Join(s.Config.SystemLibraryDir, "LaunchDaemons"),
	} {
		items = appendMatchingEntries(items, dir, CategoryStartupItem, func(name string) bool {
			lower := strings.ToLower(name)
			return strings.Contains(lower, strings.ToLower(bundleID)) || strings.Contains(lower, strings.ToLower(app.Name))
		}, "bundle identifier or application name")
	}
	return NewIterator(items)
}

// QuickLook scans QuickLook plugin locations.
type QuickLook struct {
	Config Config
}

// Scan finds QuickLook plugins matching the app name.
func (s QuickLook) Scan(ctx context.Context, app catalog.App) Iterator {
	return scanNameContains(app, CategoryQuickLook, []string{
		filepath.Join(s.Config.HomeDir, "Library/QuickLook"),
		filepath.Join(s.Config.SystemLibraryDir, "QuickLook"),
	}, "application name")
}

// ScreenSavers scans user screen saver locations.
type ScreenSavers struct {
	Config Config
}

// Scan finds screen savers matching the app name.
func (s ScreenSavers) Scan(ctx context.Context, app catalog.App) Iterator {
	return scanNameContains(app, CategoryScreenSaver, []string{
		filepath.Join(s.Config.HomeDir, "Library/Screen Savers"),
	}, "application name")
}

// InputMethods scans input method locations.
type InputMethods struct {
	Config Config
}

// Scan finds input methods matching the app name.
func (s InputMethods) Scan(ctx context.Context, app catalog.App) Iterator {
	return scanNameContains(app, CategoryInputMethod, []string{
		filepath.Join(s.Config.HomeDir, "Library/Input Methods"),
		filepath.Join(s.Config.SystemLibraryDir, "Input Methods"),
	}, "application name")
}

// Fonts scans font locations.
type Fonts struct {
	Config Config
}

// Scan finds fonts matching the app name.
func (s Fonts) Scan(ctx context.Context, app catalog.App) Iterator {
	return scanNameContains(app, CategoryFont, []string{
		filepath.Join(s.Config.HomeDir, "Library/Fonts"),
		filepath.Join(s.Config.SystemLibraryDir, "Fonts"),
	}, "application name")
}

func scanNameContains(app catalog.App, category Category, dirs []string, reason string) Iterator {
	items := []Item{}
	for _, dir := range dirs {
		items = appendMatchingEntries(items, dir, category, func(name string) bool {
			return containsFold(name, app.Name)
		}, reason)
	}
	return NewIterator(items)
}
