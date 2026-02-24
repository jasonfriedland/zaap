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
	ControlPanels   []string
	StartupItems    []string
	QuickLook       []string
	ScreenSavers    []string
	InputMethods    []string
	Fonts           []string
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "zaap",
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

	printCategory("Associated files", app.AssociatedFiles)
	printCategory("Control Panels", app.ControlPanels)
	printCategory("Startup Items", app.StartupItems)
	printCategory("QuickLook Plugins", app.QuickLook)
	printCategory("Screen Savers", app.ScreenSavers)
	printCategory("Input Methods", app.InputMethods)
	printCategory("Fonts", app.Fonts)

	if len(app.AssociatedFiles)+len(app.ControlPanels)+len(app.StartupItems)+len(app.QuickLook)+len(app.ScreenSavers)+len(app.InputMethods)+len(app.Fonts) == 0 {
		fmt.Println("\nNo associated items found.")
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

	allItems := app.AssociatedFiles
	allItems = append(allItems, app.ControlPanels...)
	allItems = append(allItems, app.StartupItems...)
	allItems = append(allItems, app.QuickLook...)
	allItems = append(allItems, app.ScreenSavers...)
	allItems = append(allItems, app.InputMethods...)
	allItems = append(allItems, app.Fonts...)

	if len(allItems) > 0 {
		fmt.Println("\nDelete associated items? (y/n/all): ")
		line, _ = reader.ReadString('\n')
		line = strings.TrimSpace(line)

		if line == "all" {
			for _, f := range allItems {
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
			for _, f := range allItems {
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

	printCategory("Associated files", target.AssociatedFiles)
	printCategory("Control Panels", target.ControlPanels)
	printCategory("Startup Items", target.StartupItems)
	printCategory("QuickLook Plugins", target.QuickLook)
	printCategory("Screen Savers", target.ScreenSavers)
	printCategory("Input Methods", target.InputMethods)
	printCategory("Fonts", target.Fonts)

	action := "Would delete"
	if !dryRun {
		if err := deletePath(target.Path); err != nil {
			fmt.Fprintf(os.Stderr, "Error deleting app: %v\n", err)
			os.Exit(1)
		}
		action = "Deleted"
	}
	fmt.Printf("%s: %s\n", action, target.Path)

	allItems := target.AssociatedFiles
	allItems = append(allItems, target.ControlPanels...)
	allItems = append(allItems, target.StartupItems...)
	allItems = append(allItems, target.QuickLook...)
	allItems = append(allItems, target.ScreenSavers...)
	allItems = append(allItems, target.InputMethods...)
	allItems = append(allItems, target.Fonts...)

	for _, f := range allItems {
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

	scanControlPanels(app, bundleID)
	scanStartupItems(app, bundleID)
	scanQuickLook(app, bundleID)
	scanScreenSavers(app, bundleID)
	scanInputMethods(app, bundleID)
	scanFonts(app, bundleID)
}

func scanControlPanels(app *AppInfo, bundleID string) {
	home := os.Getenv("HOME")
	locations := []string{
		filepath.Join(home, "Library/PreferencePanes"),
		"/Library/PreferencePanes",
	}
	for _, dir := range locations {
		if entries, err := os.ReadDir(dir); err == nil {
			for _, entry := range entries {
				if strings.Contains(strings.ToLower(entry.Name()), strings.ToLower(app.Name)) {
					app.ControlPanels = append(app.ControlPanels, filepath.Join(dir, entry.Name()))
				}
			}
		}
	}
}

func scanStartupItems(app *AppInfo, bundleID string) {
	home := os.Getenv("HOME")
	locations := []string{
		filepath.Join(home, "Library/LaunchAgents"),
		"/Library/LaunchAgents",
		filepath.Join(home, "Library/LaunchDaemons"),
		"/Library/LaunchDaemons",
	}
	for _, dir := range locations {
		if entries, err := os.ReadDir(dir); err == nil {
			for _, entry := range entries {
				name := strings.ToLower(entry.Name())
				appName := strings.ToLower(app.Name)
				if strings.Contains(name, bundleID) || strings.Contains(name, appName) {
					app.StartupItems = append(app.StartupItems, filepath.Join(dir, entry.Name()))
				}
			}
		}
	}
}

func scanQuickLook(app *AppInfo, bundleID string) {
	home := os.Getenv("HOME")
	locations := []string{
		filepath.Join(home, "Library/QuickLook"),
		"/Library/QuickLook",
	}
	for _, dir := range locations {
		if entries, err := os.ReadDir(dir); err == nil {
			for _, entry := range entries {
				if strings.Contains(strings.ToLower(entry.Name()), strings.ToLower(app.Name)) {
					app.QuickLook = append(app.QuickLook, filepath.Join(dir, entry.Name()))
				}
			}
		}
	}
}

func scanScreenSavers(app *AppInfo, bundleID string) {
	home := os.Getenv("HOME")
	dir := filepath.Join(home, "Library/Screen Savers")
	if entries, err := os.ReadDir(dir); err == nil {
		for _, entry := range entries {
			if strings.Contains(strings.ToLower(entry.Name()), strings.ToLower(app.Name)) {
				app.ScreenSavers = append(app.ScreenSavers, filepath.Join(dir, entry.Name()))
			}
		}
	}
}

func scanInputMethods(app *AppInfo, bundleID string) {
	home := os.Getenv("HOME")
	locations := []string{
		filepath.Join(home, "Library/Input Methods"),
		"/Library/Input Methods",
	}
	for _, dir := range locations {
		if entries, err := os.ReadDir(dir); err == nil {
			for _, entry := range entries {
				if strings.Contains(strings.ToLower(entry.Name()), strings.ToLower(app.Name)) {
					app.InputMethods = append(app.InputMethods, filepath.Join(dir, entry.Name()))
				}
			}
		}
	}
}

func scanFonts(app *AppInfo, bundleID string) {
	home := os.Getenv("HOME")
	locations := []string{
		filepath.Join(home, "Library/Fonts"),
		"/Library/Fonts",
	}
	for _, dir := range locations {
		if entries, err := os.ReadDir(dir); err == nil {
			for _, entry := range entries {
				if strings.Contains(strings.ToLower(entry.Name()), strings.ToLower(app.Name)) {
					app.Fonts = append(app.Fonts, filepath.Join(dir, entry.Name()))
				}
			}
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

func printCategory(name string, items []string) {
	if len(items) > 0 {
		fmt.Printf("\n%s:\n", name)
		for _, f := range items {
			fmt.Printf("  - %s\n", f)
		}
	}
}
