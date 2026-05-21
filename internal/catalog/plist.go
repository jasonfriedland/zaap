package catalog

import (
	"os"
	"path/filepath"

	"howett.net/plist"
)

// BundleID reads CFBundleIdentifier from an app bundle Info.plist.
func BundleID(appPath string) string {
	values, err := readInfoPlist(appPath)
	if err != nil {
		return ""
	}
	return values["CFBundleIdentifier"]
}

func readInfoPlist(appPath string) (map[string]string, error) {
	data, err := os.ReadFile(filepath.Join(appPath, "Contents", "Info.plist"))
	if err != nil {
		return nil, err
	}

	values := make(map[string]string)
	if _, err := plist.Unmarshal(data, &values); err != nil {
		return nil, err
	}
	return values, nil
}
