package catalog

import (
	"os"
	"path/filepath"
	"testing"
)

func createApp(t *testing.T, appsDir, name, bundleID string) string {
	t.Helper()
	appPath := filepath.Join(appsDir, name+".app")
	contentsPath := filepath.Join(appPath, "Contents")
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
</dict>
</plist>`
	if err := os.WriteFile(filepath.Join(contentsPath, "Info.plist"), []byte(infoPlist), 0644); err != nil {
		t.Fatalf("failed to write Info.plist: %v", err)
	}
	return appPath
}

func TestApplications(t *testing.T) {
	appsDir := filepath.Join(t.TempDir(), "Applications")
	if err := os.MkdirAll(appsDir, 0755); err != nil {
		t.Fatalf("failed to create apps dir: %v", err)
	}
	appPath := createApp(t, appsDir, "TestApp", "com.test.app")
	if err := os.WriteFile(filepath.Join(appsDir, "README.txt"), []byte("ignore"), 0644); err != nil {
		t.Fatalf("failed to create non-app file: %v", err)
	}

	apps, err := Applications(appsDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(apps) != 1 {
		t.Fatalf("expected 1 app, got %d", len(apps))
	}
	if apps[0].Name != "TestApp" {
		t.Errorf("expected TestApp, got %s", apps[0].Name)
	}
	if apps[0].Path != appPath {
		t.Errorf("expected %s, got %s", appPath, apps[0].Path)
	}
	if apps[0].BundleID != "com.test.app" {
		t.Errorf("expected bundle ID com.test.app, got %s", apps[0].BundleID)
	}
}

func TestBundleIDParsesInfoPlist(t *testing.T) {
	appPath := createApp(t, t.TempDir(), "Chrome", "com.google.Chrome")

	if got := BundleID(appPath); got != "com.google.Chrome" {
		t.Errorf("expected com.google.Chrome, got %s", got)
	}
	if got := BundleID("/nonexistent/app.app"); got != "" {
		t.Errorf("expected empty bundle ID for nonexistent app, got %s", got)
	}
}
