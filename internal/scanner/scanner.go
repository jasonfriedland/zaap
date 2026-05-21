package scanner

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/jasonfriedland/zaap/internal/catalog"
)

// Scanner finds files related to an application.
type Scanner interface {
	Scan(ctx context.Context, app catalog.App) Iterator
}

// Config defines filesystem roots used by scanners.
type Config struct {
	HomeDir          string
	SystemLibraryDir string
}

// DefaultConfig returns scanner roots for the current machine.
func DefaultConfig() Config {
	return Config{
		HomeDir:          os.Getenv("HOME"),
		SystemLibraryDir: "/Library",
	}
}

// Composite runs multiple scanners and deduplicates by path.
type Composite struct {
	scanners []Scanner
}

// NewComposite combines scanners into one scanner.
func NewComposite(scanners ...Scanner) *Composite {
	return &Composite{scanners: scanners}
}

// Default returns the standard zaap scanner set.
func Default(config Config) *Composite {
	return NewComposite(
		AssociatedFiles{Config: config},
		ControlPanels{Config: config},
		StartupItems{Config: config},
		QuickLook{Config: config},
		ScreenSavers{Config: config},
		InputMethods{Config: config},
		Fonts{Config: config},
	)
}

// Scan runs child scanners and returns unique items.
func (s *Composite) Scan(ctx context.Context, app catalog.App) Iterator {
	seen := make(map[string]struct{})
	var items []Item

	for _, child := range s.scanners {
		it := child.Scan(ctx, app)
		for it.Next() {
			item := it.Item()
			if _, ok := seen[item.Path]; ok {
				continue
			}
			seen[item.Path] = struct{}{}
			items = append(items, item)
		}
		if it.Err() != nil {
			return NewIterator(items)
		}
	}

	return NewIterator(items)
}

func identity(app catalog.App) string {
	if app.BundleID != "" {
		return app.BundleID
	}
	if bundleID := catalog.BundleID(app.Path); bundleID != "" {
		return bundleID
	}
	return strings.ReplaceAll(app.Name, " ", "")
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func appendExisting(items []Item, category Category, paths []string, reason string) []Item {
	for _, path := range paths {
		if exists(path) {
			items = append(items, Item{Path: path, Category: category, MatchReason: reason})
		}
	}
	return items
}

func appendMatchingEntries(items []Item, dir string, category Category, match func(name string) bool, reason string) []Item {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return items
	}
	for _, entry := range entries {
		if match(entry.Name()) {
			items = append(items, Item{Path: filepath.Join(dir, entry.Name()), Category: category, MatchReason: reason})
		}
	}
	return items
}

func containsFold(value, needle string) bool {
	return strings.Contains(strings.ToLower(value), strings.ToLower(needle))
}
