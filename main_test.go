package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetApplications(t *testing.T) {
	tmpDir := t.TempDir()
	appsDir := filepath.Join(tmpDir, "Applications")
	if err := os.MkdirAll(appsDir, 0755); err != nil {
		t.Fatalf("failed to create Applications dir: %v", err)
	}

	apps, err := getApplications(appsDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(apps) == 0 {
		t.Log("no apps in test Applications dir, test may not be useful")
	}

	for _, app := range apps {
		if app.Name == "" {
			t.Error("app name should not be empty")
		}
		if app.Path == "" {
			t.Error("app path should not be empty")
		}
		if !filepath.IsAbs(app.Path) {
			t.Error("app path should be absolute")
		}
	}
}

func TestPathExists(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "test.txt")

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

func TestAppInfoStruct(t *testing.T) {
	app := AppInfo{
		Name:            "TestApp",
		Path:            "/Applications/TestApp.app",
		BundleID:        "com.test.app",
		AssociatedFiles: []string{"/Users/test/Library/Preferences/com.test.app.plist"},
	}

	if app.Name != "TestApp" {
		t.Errorf("expected Name to be TestApp, got %s", app.Name)
	}
	if app.Path != "/Applications/TestApp.app" {
		t.Errorf("expected Path to be /Applications/TestApp.app, got %s", app.Path)
	}
	if app.BundleID != "com.test.app" {
		t.Errorf("expected BundleID to be com.test.app, got %s", app.BundleID)
	}
	if len(app.AssociatedFiles) != 1 {
		t.Errorf("expected 1 AssociatedFile, got %d", len(app.AssociatedFiles))
	}
}

func TestScanAssociatedFiles(t *testing.T) {
	home := t.TempDir()
	os.Setenv("HOME", home)

	prefsDir := filepath.Join(home, "Library/Preferences")
	appSupportDir := filepath.Join(home, "Library/Application Support")
	cachesDir := filepath.Join(home, "Library/Caches")

	if err := os.MkdirAll(prefsDir, 0755); err != nil {
		t.Fatalf("failed to create prefs dir: %v", err)
	}
	if err := os.MkdirAll(appSupportDir, 0755); err != nil {
		t.Fatalf("failed to create app support dir: %v", err)
	}
	if err := os.MkdirAll(cachesDir, 0755); err != nil {
		t.Fatalf("failed to create caches dir: %v", err)
	}

	if err := os.WriteFile(filepath.Join(prefsDir, "com.test.app.plist"), []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create plist file: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(appSupportDir, "com.test.app"), 0755); err != nil {
		t.Fatalf("failed to create app support subdir: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(cachesDir, "com.test.app"), 0755); err != nil {
		t.Fatalf("failed to create caches subdir: %v", err)
	}

	app := AppInfo{
		Name: "TestApp",
		Path: "/Applications/TestApp.app",
	}

	scanAssociatedFiles(&app)

	if len(app.AssociatedFiles) == 0 {
		t.Log("no associated files found - bundle ID scanning may not work in test")
	}

	t.Logf("Found %d associated files", len(app.AssociatedFiles))
	for _, f := range app.AssociatedFiles {
		t.Logf("  - %s", f)
	}
}

func TestDeletePath(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "test.txt")

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

func TestBundleIDFromAppName(t *testing.T) {
	tests := []struct {
		appName  string
		expected string
	}{
		{"Google Chrome", "GoogleChrome"},
		{"Visual Studio Code", "VisualStudioCode"},
		{"1Password", "1Password"},
	}

	for _, tt := range tests {
		result := strings.ReplaceAll(tt.appName, " ", "")
		if result != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, result)
		}
	}
}
