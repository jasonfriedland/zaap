package scanner

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/jasonfriedland/zaap/internal/catalog"
)

type AssociatedFiles struct {
	Config Config
}

func (s AssociatedFiles) Scan(ctx context.Context, app catalog.App) Iterator {
	bundleID := identity(app)
	home := s.Config.HomeDir
	items := []Item{}

	items = appendExisting(items, CategoryAssociated, []string{
		filepath.Join(home, "Library/Preferences", bundleID+".plist"),
		filepath.Join(home, "Library/Preferences", bundleID),
		filepath.Join(home, "Library/Application Support", bundleID),
		filepath.Join(home, "Library/Caches", bundleID),
		filepath.Join(home, "Library/Logs", bundleID),
		filepath.Join(home, "Library/Saved Application State", bundleID+".savedState"),
		filepath.Join(home, "Library/Containers", bundleID),
	}, "bundle identifier")

	items = appendMatchingEntries(items, filepath.Join(home, "Library/Application Support"), CategoryAssociated, func(name string) bool {
		return containsFold(name, app.Name)
	}, "application name")

	items = appendMatchingEntries(items, filepath.Join(home, "Library/Preferences"), CategoryAssociated, func(name string) bool {
		return strings.HasPrefix(name, bundleID)
	}, "bundle identifier prefix")

	items = appendMatchingEntries(items, filepath.Join(home, "Library/Caches"), CategoryAssociated, func(name string) bool {
		return containsFold(name, app.Name)
	}, "application name")

	return NewIterator(items)
}
