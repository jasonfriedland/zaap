package scanner

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/jasonfriedland/zaap/internal/catalog"
)

type ControlPanels struct {
	Config Config
}

func (s ControlPanels) Scan(ctx context.Context, app catalog.App) Iterator {
	return scanNameContains(app, CategoryControlPane, []string{
		filepath.Join(s.Config.HomeDir, "Library/PreferencePanes"),
		filepath.Join(s.Config.SystemLibraryDir, "PreferencePanes"),
	}, "application name")
}

type StartupItems struct {
	Config Config
}

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

type QuickLook struct {
	Config Config
}

func (s QuickLook) Scan(ctx context.Context, app catalog.App) Iterator {
	return scanNameContains(app, CategoryQuickLook, []string{
		filepath.Join(s.Config.HomeDir, "Library/QuickLook"),
		filepath.Join(s.Config.SystemLibraryDir, "QuickLook"),
	}, "application name")
}

type ScreenSavers struct {
	Config Config
}

func (s ScreenSavers) Scan(ctx context.Context, app catalog.App) Iterator {
	return scanNameContains(app, CategoryScreenSaver, []string{
		filepath.Join(s.Config.HomeDir, "Library/Screen Savers"),
	}, "application name")
}

type InputMethods struct {
	Config Config
}

func (s InputMethods) Scan(ctx context.Context, app catalog.App) Iterator {
	return scanNameContains(app, CategoryInputMethod, []string{
		filepath.Join(s.Config.HomeDir, "Library/Input Methods"),
		filepath.Join(s.Config.SystemLibraryDir, "Input Methods"),
	}, "application name")
}

type Fonts struct {
	Config Config
}

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
