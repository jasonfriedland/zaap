package catalog

import (
	"os"
	"path/filepath"
	"strings"
)

// Applications returns .app bundles found directly under dir.
func Applications(dir string) ([]App, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var apps []App
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".app") {
			path := filepath.Join(dir, entry.Name())
			apps = append(apps, App{
				Name:     strings.TrimSuffix(entry.Name(), ".app"),
				Path:     path,
				BundleID: BundleID(path),
			})
		}
	}
	return apps, nil
}
