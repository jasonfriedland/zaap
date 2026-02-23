package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	verbose    bool
	listOnly   bool
	deleteName string
	dryRun     bool
)

type AppInfo struct {
	Name            string
	Path            string
	BundleID        string
	AssociatedFiles []string
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "appcleaner",
		Short: "macOS application cleanup utility",
		Run:   run,
	}

	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.Flags().BoolVarP(&listOnly, "list", "l", false, "list applications only")
	rootCmd.Flags().StringVarP(&deleteName, "delete", "d", "", "delete specific app by name")
	rootCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "show what would be deleted without actually deleting")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	if listOnly {
		listApplications()
		return
	}

	if deleteName != "" {
		deleteApp(deleteName)
		return
	}

	interactiveMode()
}

func listApplications() {
	apps, err := getApplications("/Applications")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Installed Applications:")
	fmt.Println("----------------------")
	for i, app := range apps {
		fmt.Printf("%d. %s\n", i+1, app.Name)
	}
}

func interactiveMode() {
	apps, err := getApplications("/Applications")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Select an application to delete:")
	fmt.Println("---------------------------------")
	for i, app := range apps {
		fmt.Printf("%d. %s\n", i+1, app.Name)
	}
	fmt.Println("0. Exit")

	fmt.Print("\nEnter number: ")
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)

	var selection int
	if _, err := fmt.Sscanf(line, "%d", &selection); err != nil {
		fmt.Println("Invalid input. Exiting.")
		os.Exit(1)
	}

	if selection == 0 || selection < 0 || selection > len(apps) {
		fmt.Println("Exiting.")
		os.Exit(0)
	}

	app := apps[selection-1]
	scanAssociatedFiles(&app)

	fmt.Printf("\nSelected: %s\n", app.Name)
	fmt.Printf("Location: %s\n", app.Path)

	if len(app.AssociatedFiles) > 0 {
		fmt.Println("\nAssociated files found:")
		for _, f := range app.AssociatedFiles {
			fmt.Printf("  - %s\n", f)
		}
	} else {
		fmt.Println("\nNo associated files found.")
	}

	fmt.Print("\nDelete this application? (y/n): ")
	line, _ = reader.ReadString('\n')
	line = strings.TrimSpace(line)

	if strings.ToLower(line) != "y" {
		fmt.Println("Cancelled.")
		os.Exit(0)
	}

	if dryRun {
		fmt.Printf("Would delete: %s\n", app.Path)
	} else {
		if err := deletePath(app.Path); err != nil {
			fmt.Fprintf(os.Stderr, "Error deleting app: %v\n", err)
		} else {
			fmt.Printf("Deleted: %s\n", app.Path)
		}
	}

	if len(app.AssociatedFiles) > 0 {
		fmt.Println("\nDelete associated files? (y/n/all): ")
		line, _ = reader.ReadString('\n')
		line = strings.TrimSpace(line)

		if line == "all" {
			for _, f := range app.AssociatedFiles {
				if dryRun {
					fmt.Printf("Would delete: %s\n", f)
				} else {
					if err := deletePath(f); err != nil {
						fmt.Fprintf(os.Stderr, "Error deleting %s: %v\n", f, err)
					} else {
						fmt.Printf("Deleted: %s\n", f)
					}
				}
			}
		} else if strings.ToLower(line) == "y" {
			for _, f := range app.AssociatedFiles {
				fmt.Printf("Delete %s? (y/n): ", filepath.Base(f))
				line, _ := reader.ReadString('\n')
				line = strings.TrimSpace(line)
				if strings.ToLower(line) == "y" {
					if dryRun {
						fmt.Printf("Would delete: %s\n", f)
					} else {
						if err := deletePath(f); err != nil {
							fmt.Fprintf(os.Stderr, "Error deleting %s: %v\n", f, err)
						} else {
							fmt.Printf("Deleted: %s\n", f)
						}
					}
				}
			}
		}
	}

	if dryRun {
		fmt.Println("\nDry run complete. No files were actually deleted.")
	}

	fmt.Println("\nDone!")
}

func deleteApp(name string) {
	apps, err := getApplications("/Applications")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	var target AppInfo
	found := false
	for _, app := range apps {
		if strings.EqualFold(app.Name, name) {
			target = app
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("Application not found: %s\n", name)
		os.Exit(1)
	}

	scanAssociatedFiles(&target)

	fmt.Printf("Selected: %s\n", target.Name)
	fmt.Printf("Location: %s\n", target.Path)

	if len(target.AssociatedFiles) > 0 {
		fmt.Println("\nAssociated files found:")
		for _, f := range target.AssociatedFiles {
			fmt.Printf("  - %s\n", f)
		}
	}

	action := "Would delete"
	if !dryRun {
		if err := deletePath(target.Path); err != nil {
			fmt.Fprintf(os.Stderr, "Error deleting app: %v\n", err)
			os.Exit(1)
		}
		action = "Deleted"
	}
	fmt.Printf("%s: %s\n", action, target.Path)

	for _, f := range target.AssociatedFiles {
		if dryRun {
			fmt.Printf("Would delete: %s\n", f)
		} else {
			if err := deletePath(f); err != nil {
				fmt.Fprintf(os.Stderr, "Error deleting %s: %v\n", f, err)
			} else {
				fmt.Printf("Deleted: %s\n", f)
			}
		}
	}

	if dryRun {
		fmt.Println("\nDry run complete. No files were actually deleted.")
	}
}

func getApplications(dir string) ([]AppInfo, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var apps []AppInfo
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".app") {
			apps = append(apps, AppInfo{
				Name: strings.TrimSuffix(entry.Name(), ".app"),
				Path: filepath.Join(dir, entry.Name()),
			})
		}
	}
	return apps, nil
}

func scanAssociatedFiles(app *AppInfo) {
	bundleID := getBundleID(app.Path)
	if bundleID == "" {
		bundleID = strings.ReplaceAll(app.Name, " ", "")
	}

	if verbose {
		fmt.Printf("Bundle ID: %s\n", bundleID)
	}

	home := os.Getenv("HOME")
	locations := []string{
		filepath.Join(home, "Library/Preferences", bundleID+".plist"),
		filepath.Join(home, "Library/Preferences", bundleID),
		filepath.Join(home, "Library/Application Support", bundleID),
		filepath.Join(home, "Library/Caches", bundleID),
		filepath.Join(home, "Library/Logs", bundleID),
		filepath.Join(home, "Library/Saved Application State", bundleID+".savedState"),
		filepath.Join(home, "Library/Containers", bundleID),
	}

	appSupportDir := filepath.Join(home, "Library/Application Support")
	if entries, err := os.ReadDir(appSupportDir); err == nil {
		for _, entry := range entries {
			if strings.Contains(strings.ToLower(entry.Name()), strings.ToLower(app.Name)) {
				locations = append(locations, filepath.Join(appSupportDir, entry.Name()))
			}
		}
	}

	prefsDir := filepath.Join(home, "Library/Preferences")
	if entries, err := os.ReadDir(prefsDir); err == nil {
		for _, entry := range entries {
			if strings.HasPrefix(entry.Name(), bundleID) {
				locations = append(locations, filepath.Join(prefsDir, entry.Name()))
			}
		}
	}

	cachesDir := filepath.Join(home, "Library/Caches")
	if entries, err := os.ReadDir(cachesDir); err == nil {
		for _, entry := range entries {
			if strings.Contains(strings.ToLower(entry.Name()), strings.ToLower(app.Name)) {
				locations = append(locations, filepath.Join(cachesDir, entry.Name()))
			}
		}
	}

	for _, loc := range locations {
		if exists, _ := pathExists(loc); exists {
			app.AssociatedFiles = append(app.AssociatedFiles, loc)
		}
	}
}

func getBundleID(appPath string) string {
	cmd := exec.Command("defaults", "read", filepath.Join(appPath, "Contents/Info"), "CFBundleIdentifier")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func deletePath(path string) error {
	return os.RemoveAll(path)
}
