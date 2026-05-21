package scanner

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jasonfriedland/zaap/internal/catalog"
)

type testFS struct {
	rootDir   string
	homeDir   string
	systemDir string
}

func newTestFS(t *testing.T) *testFS {
	t.Helper()
	tmpDir := t.TempDir()
	fs := &testFS{
		rootDir:   tmpDir,
		homeDir:   filepath.Join(tmpDir, "home"),
		systemDir: filepath.Join(tmpDir, "Library"),
	}
	dirs := []string{
		filepath.Join(fs.homeDir, "Applications"),
		filepath.Join(fs.homeDir, "Library", "Preferences"),
		filepath.Join(fs.homeDir, "Library", "Application Support"),
		filepath.Join(fs.homeDir, "Library", "Caches"),
		filepath.Join(fs.homeDir, "Library", "Logs"),
		filepath.Join(fs.homeDir, "Library", "Saved Application State"),
		filepath.Join(fs.homeDir, "Library", "Containers"),
		filepath.Join(fs.homeDir, "Library", "PreferencePanes"),
		filepath.Join(fs.homeDir, "Library", "LaunchAgents"),
		filepath.Join(fs.homeDir, "Library", "LaunchDaemons"),
		filepath.Join(fs.homeDir, "Library", "QuickLook"),
		filepath.Join(fs.homeDir, "Library", "Screen Savers"),
		filepath.Join(fs.homeDir, "Library", "Input Methods"),
		filepath.Join(fs.homeDir, "Library", "Fonts"),
		filepath.Join(fs.systemDir, "PreferencePanes"),
		filepath.Join(fs.systemDir, "LaunchAgents"),
		filepath.Join(fs.systemDir, "LaunchDaemons"),
		filepath.Join(fs.systemDir, "QuickLook"),
		filepath.Join(fs.systemDir, "Input Methods"),
		filepath.Join(fs.systemDir, "Fonts"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("failed to create dir %s: %v", dir, err)
		}
	}
	return fs
}

func (fs *testFS) config() Config {
	return Config{HomeDir: fs.homeDir, SystemLibraryDir: fs.systemDir}
}

func (fs *testFS) createApp(t *testing.T, name, bundleID string) catalog.App {
	t.Helper()
	appsDir := filepath.Join(fs.homeDir, "Applications")
	appPath := filepath.Join(appsDir, name+".app")
	contentsPath := filepath.Join(appPath, "Contents")
	if err := os.MkdirAll(contentsPath, 0755); err != nil {
		t.Fatalf("failed to create app contents: %v", err)
	}
	infoPlist := `<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0">
<dict>
	<key>CFBundleIdentifier</key>
	<string>` + bundleID + `</string>
	<key>CFBundleName</key>
	<string>` + name + `</string>
</dict>
</plist>`
	if err := os.WriteFile(filepath.Join(contentsPath, "Info.plist"), []byte(infoPlist), 0644); err != nil {
		t.Fatalf("failed to write Info.plist: %v", err)
	}
	return catalog.App{Name: name, Path: appPath, BundleID: bundleID}
}

func (fs *testFS) writeFile(t *testing.T, parts ...string) string {
	t.Helper()
	path := filepath.Join(parts...)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("failed to create parent dir: %v", err)
	}
	if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to write file %s: %v", path, err)
	}
	return path
}

func (fs *testFS) mkdir(t *testing.T, parts ...string) string {
	t.Helper()
	path := filepath.Join(parts...)
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatalf("failed to create dir %s: %v", path, err)
	}
	return path
}

func TestAssociatedFilesScanner(t *testing.T) {
	fs := newTestFS(t)
	app := fs.createApp(t, "TestApp", "com.test.app")
	pref := fs.writeFile(t, fs.homeDir, "Library", "Preferences", "com.test.app.plist")
	appSupport := fs.mkdir(t, fs.homeDir, "Library", "Application Support", "com.test.app")
	caches := fs.mkdir(t, fs.homeDir, "Library", "Caches", "com.test.app")

	items, err := Collect(AssociatedFiles{Config: fs.config()}.Scan(context.Background(), app))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertHasItem(t, items, pref, CategoryAssociated)
	assertHasItem(t, items, appSupport, CategoryAssociated)
	assertHasItem(t, items, caches, CategoryAssociated)
}

func TestExtensionScanners(t *testing.T) {
	fs := newTestFS(t)
	app := fs.createApp(t, "TestApp", "com.test.app")
	paths := map[Category]string{
		CategoryControlPane: fs.mkdir(t, fs.homeDir, "Library", "PreferencePanes", "TestApp.prefPane"),
		CategoryStartupItem: fs.writeFile(t, fs.homeDir, "Library", "LaunchAgents", "com.test.app.plist"),
		CategoryQuickLook:   fs.mkdir(t, fs.homeDir, "Library", "QuickLook", "TestApp.qlgenerator"),
		CategoryScreenSaver: fs.mkdir(t, fs.homeDir, "Library", "Screen Savers", "TestApp.saver"),
		CategoryInputMethod: fs.mkdir(t, fs.homeDir, "Library", "Input Methods", "TestApp.app"),
		CategoryFont:        fs.writeFile(t, fs.homeDir, "Library", "Fonts", "TestApp.ttf"),
	}

	items, err := Collect(Default(fs.config()).Scan(context.Background(), app))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for category, path := range paths {
		assertHasItem(t, items, path, category)
	}
}

func TestCompositeDeduplicatesByPath(t *testing.T) {
	item := Item{Path: "/tmp/example", Category: CategoryAssociated}
	duplicate := Item{Path: "/tmp/example", Category: CategoryStartupItem}
	composite := NewComposite(staticScanner{items: []Item{item}}, staticScanner{items: []Item{duplicate}})

	items, err := Collect(composite.Scan(context.Background(), catalog.App{Name: "Example"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 deduplicated item, got %d", len(items))
	}
}

type staticScanner struct {
	items []Item
}

func (s staticScanner) Scan(ctx context.Context, app catalog.App) Iterator {
	return NewIterator(s.items)
}

func assertHasItem(t *testing.T, items []Item, path string, category Category) {
	t.Helper()
	for _, item := range items {
		if item.Path == path && item.Category == category {
			return
		}
	}
	t.Fatalf("expected %s item %s in %#v", category, path, items)
}
