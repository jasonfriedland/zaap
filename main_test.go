package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

type testFS struct {
	rootDir string
	homeDir string
}

func newTestFS(t *testing.T) *testFS {
	tmpDir := t.TempDir()
	homeDir := filepath.Join(tmpDir, "home")
	userAppsDir := filepath.Join(homeDir, "Applications")

	structure := []string{
		homeDir,
		userAppsDir,
		filepath.Join(homeDir, "Library", "Preferences"),
		filepath.Join(homeDir, "Library", "Application Support"),
		filepath.Join(homeDir, "Library", "Caches"),
		filepath.Join(homeDir, "Library", "Logs"),
		filepath.Join(homeDir, "Library", "Saved Application State"),
		filepath.Join(homeDir, "Library", "Containers"),
		filepath.Join(homeDir, "Library", "PreferencePanes"),
		filepath.Join(homeDir, "Library", "LaunchAgents"),
		filepath.Join(homeDir, "Library", "LaunchDaemons"),
		filepath.Join(homeDir, "Library", "QuickLook"),
		filepath.Join(homeDir, "Library", "Screen Savers"),
		filepath.Join(homeDir, "Library", "Input Methods"),
		filepath.Join(homeDir, "Library", "Fonts"),
		filepath.Join(homeDir, "Library", "Fonts", "CustomFonts"),
	}

	for _, dir := range structure {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("failed to create dir %s: %v", dir, err)
		}
	}

	return &testFS{
		rootDir: tmpDir,
		homeDir: homeDir,
	}
}

func (fs *testFS) createApp(t *testing.T, name, bundleID string) string {
	appsDir := filepath.Join(fs.homeDir, "Applications")
	appPath := filepath.Join(appsDir, name+".app")
	contentsPath := filepath.Join(appPath, "Contents")
	infoPlistPath := filepath.Join(contentsPath, "Info.plist")

	if err := os.MkdirAll(contentsPath, 0755); err != nil {
		t.Fatalf("failed to create app contents: %v", err)
	}

	infoPlist := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleIdentifier</key>
	<string>` + bundleID + `</string>
	<key>CFBundleName</key>
	<string>` + name + `</string>
	<key>CFBundleExecutable</key>
	<string>` + name + `</string>
</dict>
</plist>`

	if err := os.WriteFile(infoPlistPath, []byte(infoPlist), 0644); err != nil {
		t.Fatalf("failed to write Info.plist: %v", err)
	}

	return appPath
}

func (fs *testFS) createPrefFile(t *testing.T, bundleID, suffix string) string {
	prefsDir := filepath.Join(fs.homeDir, "Library", "Preferences")
	filename := bundleID + suffix
	path := filepath.Join(prefsDir, filename)
	if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create pref file: %v", err)
	}
	return path
}

func (fs *testFS) createAppSupportDir(t *testing.T, bundleID string) string {
	appSupportDir := filepath.Join(fs.homeDir, "Library", "Application Support", bundleID)
	if err := os.MkdirAll(appSupportDir, 0755); err != nil {
		t.Fatalf("failed to create app support dir: %v", err)
	}
	return appSupportDir
}

func (fs *testFS) createCachesDir(t *testing.T, bundleID string) string {
	cachesDir := filepath.Join(fs.homeDir, "Library", "Caches", bundleID)
	if err := os.MkdirAll(cachesDir, 0755); err != nil {
		t.Fatalf("failed to create caches dir: %v", err)
	}
	return cachesDir
}

func (fs *testFS) createLaunchAgent(t *testing.T, bundleID string) string {
	agentsDir := filepath.Join(fs.homeDir, "Library", "LaunchAgents")
	filename := bundleID + ".plist"
	path := filepath.Join(agentsDir, filename)
	if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create launch agent: %v", err)
	}
	return path
}

func (fs *testFS) createPrefPane(t *testing.T, appName string) string {
	panesDir := filepath.Join(fs.homeDir, "Library", "PreferencePanes")
	path := filepath.Join(panesDir, appName+".prefPane")
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatalf("failed to create pref pane: %v", err)
	}
	return path
}

func (fs *testFS) createScreenSaver(t *testing.T, appName string) string {
	saverDir := filepath.Join(fs.homeDir, "Library", "Screen Savers")
	path := filepath.Join(saverDir, appName+".saver")
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatalf("failed to create screen saver: %v", err)
	}
	return path
}

func (fs *testFS) createInputMethod(t *testing.T, appName string) string {
	inputDir := filepath.Join(fs.homeDir, "Library", "Input Methods")
	path := filepath.Join(inputDir, appName+".app")
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatalf("failed to create input method: %v", err)
	}
	return path
}

func (fs *testFS) createFont(t *testing.T, appName string) string {
	fontsDir := filepath.Join(fs.homeDir, "Library", "Fonts")
	path := filepath.Join(fontsDir, appName+".ttf")
	if err := os.WriteFile(path, []byte("fontdata"), 0644); err != nil {
		t.Fatalf("failed to create font: %v", err)
	}
	return path
}

func (fs *testFS) createQuickLookPlugin(t *testing.T, appName string) string {
	qlDir := filepath.Join(fs.homeDir, "Library", "QuickLook")
	path := filepath.Join(qlDir, appName+".qlgenerator")
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatalf("failed to create quicklook plugin: %v", err)
	}
	return path
}

func TestGetApplications(t *testing.T) {
	fs := newTestFS(t)
	os.Setenv("HOME", fs.homeDir)

	appPath := fs.createApp(t, "TestApp", "com.test.app")

	apps, err := getApplications(filepath.Join(fs.homeDir, "Applications"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(apps) != 1 {
		t.Fatalf("expected 1 app, got %d", len(apps))
	}

	if apps[0].Name != "TestApp" {
		t.Errorf("expected name TestApp, got %s", apps[0].Name)
	}

	if apps[0].Path != appPath {
		t.Errorf("expected path %s, got %s", appPath, apps[0].Path)
	}
}

func TestGetBundleID(t *testing.T) {
	fs := newTestFS(t)

	appPath := fs.createApp(t, "Chrome", "com.google.Chrome")

	bundleID := getBundleID(appPath)
	if runtime.GOOS == "darwin" {
		if bundleID != "com.google.Chrome" {
			t.Errorf("expected bundle ID com.google.Chrome, got %s", bundleID)
		}
	} else {
		if bundleID != "" {
			t.Errorf("expected empty bundle ID on non-macOS, got %s", bundleID)
		}
	}

	nonexistent := getBundleID("/nonexistent/app.app")
	if nonexistent != "" {
		t.Errorf("expected empty string for nonexistent app, got %s", nonexistent)
	}
}

func TestScanAssociatedFilesWithBundleID(t *testing.T) {
	fs := newTestFS(t)
	os.Setenv("HOME", fs.homeDir)

	bundleID := "com.test.app"
	appPath := fs.createApp(t, "TestApp", bundleID)

	fs.createPrefFile(t, bundleID, ".plist")
	fs.createPrefFile(t, bundleID, "")
	fs.createAppSupportDir(t, bundleID)
	fs.createCachesDir(t, bundleID)

	fs.createPrefFile(t, "TestApp", ".plist")
	fs.createPrefFile(t, "TestApp", "")
	fs.createAppSupportDir(t, "TestApp")
	fs.createCachesDir(t, "TestApp")

	app := AppInfo{
		Name: "TestApp",
		Path: appPath,
	}

	scanAssociatedFiles(&app)

	if len(app.AssociatedFiles) == 0 {
		t.Error("expected associated files to be found")
	}

	foundPrefs := false
	foundAppSupport := false
	foundCaches := false

	for _, f := range app.AssociatedFiles {
		if filepath.Base(f) == bundleID+".plist" || filepath.Base(f) == "TestApp.plist" {
			foundPrefs = true
		}
		if filepath.Base(filepath.Dir(f)) == "Application Support" && (filepath.Base(f) == bundleID || filepath.Base(f) == "TestApp") {
			foundAppSupport = true
		}
		if filepath.Base(filepath.Dir(f)) == "Caches" && (filepath.Base(f) == bundleID || filepath.Base(f) == "TestApp") {
			foundCaches = true
		}
	}

	if !foundPrefs {
		t.Error("expected to find preferences plist")
	}
	if !foundAppSupport {
		t.Error("expected to find Application Support directory")
	}
	if !foundCaches {
		t.Error("expected to find Caches directory")
	}
}

func TestScanControlPanels(t *testing.T) {
	fs := newTestFS(t)
	os.Setenv("HOME", fs.homeDir)

	bundleID := "com.test.app"
	appPath := fs.createApp(t, "TestApp", bundleID)

	panePath := fs.createPrefPane(t, "TestApp")

	app := AppInfo{
		Name: "TestApp",
		Path: appPath,
	}

	scanControlPanels(&app, bundleID)

	if len(app.ControlPanels) != 1 {
		t.Fatalf("expected 1 control panel, got %d", len(app.ControlPanels))
	}

	if app.ControlPanels[0] != panePath {
		t.Errorf("expected %s, got %s", panePath, app.ControlPanels[0])
	}
}

func TestScanScreenSavers(t *testing.T) {
	fs := newTestFS(t)
	os.Setenv("HOME", fs.homeDir)

	bundleID := "com.test.app"
	appPath := fs.createApp(t, "TestApp", bundleID)

	saverPath := fs.createScreenSaver(t, "TestApp")

	app := AppInfo{
		Name: "TestApp",
		Path: appPath,
	}

	scanScreenSavers(&app, bundleID)

	if len(app.ScreenSavers) != 1 {
		t.Fatalf("expected 1 screen saver, got %d", len(app.ScreenSavers))
	}

	if app.ScreenSavers[0] != saverPath {
		t.Errorf("expected %s, got %s", saverPath, app.ScreenSavers[0])
	}
}

func TestScanInputMethods(t *testing.T) {
	fs := newTestFS(t)
	os.Setenv("HOME", fs.homeDir)

	bundleID := "com.test.app"
	appPath := fs.createApp(t, "TestApp", bundleID)

	inputPath := fs.createInputMethod(t, "TestApp")

	app := AppInfo{
		Name: "TestApp",
		Path: appPath,
	}

	scanInputMethods(&app, bundleID)

	if len(app.InputMethods) != 1 {
		t.Fatalf("expected 1 input method, got %d", len(app.InputMethods))
	}

	if app.InputMethods[0] != inputPath {
		t.Errorf("expected %s, got %s", inputPath, app.InputMethods[0])
	}
}

func TestScanFonts(t *testing.T) {
	fs := newTestFS(t)
	os.Setenv("HOME", fs.homeDir)

	bundleID := "com.test.app"
	appPath := fs.createApp(t, "TestApp", bundleID)

	fontPath := fs.createFont(t, "TestApp")

	app := AppInfo{
		Name: "TestApp",
		Path: appPath,
	}

	scanFonts(&app, bundleID)

	if len(app.Fonts) != 1 {
		t.Fatalf("expected 1 font, got %d", len(app.Fonts))
	}

	if app.Fonts[0] != fontPath {
		t.Errorf("expected %s, got %s", fontPath, app.Fonts[0])
	}
}

func TestScanQuickLook(t *testing.T) {
	fs := newTestFS(t)
	os.Setenv("HOME", fs.homeDir)

	bundleID := "com.test.app"
	appPath := fs.createApp(t, "TestApp", bundleID)

	qlPath := fs.createQuickLookPlugin(t, "TestApp")

	app := AppInfo{
		Name: "TestApp",
		Path: appPath,
	}

	scanQuickLook(&app, bundleID)

	if len(app.QuickLook) != 1 {
		t.Fatalf("expected 1 quicklook plugin, got %d", len(app.QuickLook))
	}

	if app.QuickLook[0] != qlPath {
		t.Errorf("expected %s, got %s", qlPath, app.QuickLook[0])
	}
}

func TestScanStartupItems(t *testing.T) {
	fs := newTestFS(t)
	os.Setenv("HOME", fs.homeDir)

	bundleID := "com.test.app"
	appPath := fs.createApp(t, "TestApp", bundleID)

	agentPath := fs.createLaunchAgent(t, bundleID)

	app := AppInfo{
		Name: "TestApp",
		Path: appPath,
	}

	scanStartupItems(&app, bundleID)

	if len(app.StartupItems) != 1 {
		t.Fatalf("expected 1 startup item, got %d", len(app.StartupItems))
	}

	if app.StartupItems[0] != agentPath {
		t.Errorf("expected %s, got %s", agentPath, app.StartupItems[0])
	}
}

func TestDeletePath(t *testing.T) {
	fs := newTestFS(t)

	tmpFile := filepath.Join(fs.rootDir, "test.txt")
	if err := os.WriteFile(tmpFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	if err := deletePath(tmpFile); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if exists, _ := pathExists(tmpFile); exists {
		t.Error("file should have been deleted")
	}

	if err := deletePath("/nonexistent/path"); err != nil {
		t.Logf("expected error for non-existent path: %v", err)
	}
}

func TestDeleteAppAndAssociatedFiles(t *testing.T) {
	fs := newTestFS(t)
	os.Setenv("HOME", fs.homeDir)

	bundleID := "com.test.app"
	appPath := fs.createApp(t, "TestApp", bundleID)

	prefPath := fs.createPrefFile(t, bundleID, ".plist")
	appSupportPath := fs.createAppSupportDir(t, bundleID)
	cachesPath := fs.createCachesDir(t, bundleID)

	app := AppInfo{
		Name:            "TestApp",
		Path:            appPath,
		AssociatedFiles: []string{prefPath, appSupportPath, cachesPath},
	}

	if err := deletePath(app.Path); err != nil {
		t.Fatalf("failed to delete app: %v", err)
	}

	if exists, _ := pathExists(app.Path); exists {
		t.Error("app should have been deleted")
	}

	for _, f := range app.AssociatedFiles {
		if err := deletePath(f); err != nil {
			t.Fatalf("failed to delete associated file %s: %v", f, err)
		}
	}

	for _, f := range app.AssociatedFiles {
		if exists, _ := pathExists(f); exists {
			t.Errorf("associated file should have been deleted: %s", f)
		}
	}
}

func TestFullScanWithMultipleAssociatedFiles(t *testing.T) {
	fs := newTestFS(t)
	os.Setenv("HOME", fs.homeDir)

	bundleID := "com.google.chrome"
	appPath := fs.createApp(t, "Google Chrome", bundleID)

	fs.createPrefFile(t, bundleID, ".plist")
	fs.createPrefFile(t, bundleID, "")
	fs.createAppSupportDir(t, bundleID)
	fs.createCachesDir(t, bundleID)
	fs.createLaunchAgent(t, bundleID)
	fs.createPrefPane(t, "Google Chrome")
	fs.createScreenSaver(t, "Google Chrome")
	fs.createInputMethod(t, "Google Chrome")
	fs.createFont(t, "Google Chrome")
	fs.createQuickLookPlugin(t, "Google Chrome")

	fs.createAppSupportDir(t, "GoogleChrome")
	fs.createCachesDir(t, "GoogleChrome")
	fs.createLaunchAgent(t, "GoogleChrome")
	fs.createLaunchAgent(t, "com.google.chrome")
	fs.createLaunchAgent(t, "Google Chrome")

	app := AppInfo{
		Name: "Google Chrome",
		Path: appPath,
	}

	scanAssociatedFiles(&app)

	if len(app.AssociatedFiles) == 0 {
		t.Error("expected associated files to be found")
	}

	if len(app.ControlPanels) == 0 {
		t.Error("expected control panels to be found")
	}

	if len(app.StartupItems) == 0 {
		t.Error("expected startup items to be found")
	}

	if len(app.ScreenSavers) == 0 {
		t.Error("expected screen savers to be found")
	}

	if len(app.InputMethods) == 0 {
		t.Error("expected input methods to be found")
	}

	if len(app.Fonts) == 0 {
		t.Error("expected fonts to be found")
	}

	if len(app.QuickLook) == 0 {
		t.Error("expected quicklook plugins to be found")
	}

	t.Logf("BundleID used: com.google.chrome (or fallback GoogleChrome on non-macOS)")
	t.Logf("AssociatedFiles: %d, ControlPanels: %d, StartupItems: %d, ScreenSavers: %d, InputMethods: %d, Fonts: %d, QuickLook: %d",
		len(app.AssociatedFiles), len(app.ControlPanels), len(app.StartupItems),
		len(app.ScreenSavers), len(app.InputMethods), len(app.Fonts), len(app.QuickLook))
}

func TestAppWithNoBundleID(t *testing.T) {
	fs := newTestFS(t)
	os.Setenv("HOME", fs.homeDir)

	appPath := fs.createApp(t, "TestApp", "")

	app := AppInfo{
		Name: "TestApp",
		Path: appPath,
	}

	scanAssociatedFiles(&app)

	t.Logf("Associated files: %v", app.AssociatedFiles)
	t.Logf("BundleID used for scanning: com.testapp (fallback from app name)")
}

func TestPathExists(t *testing.T) {
	fs := newTestFS(t)

	tmpFile := filepath.Join(fs.rootDir, "test.txt")

	exists, err := pathExists(tmpFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exists {
		t.Error("expected file to not exist")
	}

	if err := os.WriteFile(tmpFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	exists, err = pathExists(tmpFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !exists {
		t.Error("expected file to exist")
	}
}
