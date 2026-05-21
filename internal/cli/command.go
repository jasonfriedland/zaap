package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jasonfriedland/zaap/internal/catalog"
	"github.com/jasonfriedland/zaap/internal/cleanup"
	"github.com/jasonfriedland/zaap/internal/scanner"
	"github.com/jasonfriedland/zaap/internal/terminal"
)

type Options struct {
	ApplicationsDir  string
	HomeDir          string
	SystemLibraryDir string
	Verbose          bool
	ListOnly         bool
	DeleteName       string
	DryRun           bool
	In               io.Reader
	Out              io.Writer
	Err              io.Writer
}

func Execute(ctx context.Context, args []string, in io.Reader, out io.Writer, errOut io.Writer) int {
	options := Options{
		ApplicationsDir:  "/Applications",
		HomeDir:          os.Getenv("HOME"),
		SystemLibraryDir: "/Library",
		In:               in,
		Out:              out,
		Err:              errOut,
	}
	cmd := NewRootCommand(ctx, &options)
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(errOut, "Error: %v\n", err)
		return 1
	}
	return 0
}

func NewRootCommand(ctx context.Context, options *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zaap",
		Short: "macOS application cleanup utility",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(ctx, options)
		},
	}

	cmd.Flags().BoolVarP(&options.Verbose, "verbose", "v", false, "verbose output")
	cmd.Flags().BoolVarP(&options.ListOnly, "list", "l", false, "list applications only")
	cmd.Flags().StringVarP(&options.DeleteName, "delete", "d", "", "delete specific app by name")
	cmd.Flags().BoolVarP(&options.DryRun, "dry-run", "n", false, "show what would be deleted without actually deleting")
	return cmd
}

func run(ctx context.Context, options *Options) error {
	if options.ListOnly {
		return listApplications(options)
	}
	if options.DeleteName != "" {
		return deleteApp(ctx, options, options.DeleteName)
	}
	return interactiveMode(ctx, options)
}

func listApplications(options *Options) error {
	apps, err := catalog.Applications(options.ApplicationsDir)
	if err != nil {
		return err
	}
	fmt.Fprintln(options.Out, "Installed Applications:")
	fmt.Fprintln(options.Out, "----------------------")
	for i, app := range apps {
		fmt.Fprintf(options.Out, "%d. %s\n", i+1, app.Name)
	}
	return nil
}

func deleteApp(ctx context.Context, options *Options, name string) error {
	apps, err := catalog.Applications(options.ApplicationsDir)
	if err != nil {
		return err
	}

	var target catalog.App
	found := false
	for _, app := range apps {
		if strings.EqualFold(app.Name, name) {
			target = app
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("application not found: %s", name)
	}

	items, err := scan(ctx, options, target)
	if err != nil {
		return err
	}
	printSelection(options.Out, target, items)

	remover := cleanup.Remover{DryRun: options.DryRun}
	printRemoveResult(options, remover.Remove(target.Path), "app")
	for _, item := range items {
		printRemoveResult(options, remover.Remove(item.Path), item.Path)
	}
	if options.DryRun {
		fmt.Fprintln(options.Out, "\nDry run complete. No files were actually deleted.")
	}
	return nil
}

func interactiveMode(ctx context.Context, options *Options) error {
	apps, err := catalog.Applications(options.ApplicationsDir)
	if err != nil {
		return err
	}

	prompter := terminal.NewPrompter(options.In, options.Out)
	app, ok, err := prompter.SelectApp(apps)
	if err != nil {
		fmt.Fprintln(options.Out, "Invalid input. Exiting.")
		return err
	}
	if !ok {
		fmt.Fprintln(options.Out, "Exiting.")
		return nil
	}

	items, err := scan(ctx, options, app)
	if err != nil {
		return err
	}
	printSelection(options.Out, app, items)

	if len(items) == 0 {
		fmt.Fprintln(options.Out, "\nNo associated items found.")
	}

	confirmed, err := prompter.Confirm("\nDelete this application? (y/n): ")
	if err != nil {
		return err
	}
	if !confirmed {
		fmt.Fprintln(options.Out, "Cancelled.")
		return nil
	}

	remover := cleanup.Remover{DryRun: options.DryRun}
	printRemoveResult(options, remover.Remove(app.Path), "app")

	if len(items) > 0 {
		mode, err := prompter.AssociatedMode()
		if err != nil {
			return err
		}
		if mode == "all" {
			for _, item := range items {
				printRemoveResult(options, remover.Remove(item.Path), item.Path)
			}
		} else if strings.ToLower(mode) == "y" {
			for _, item := range items {
				confirmed, err := prompter.ConfirmPath(item.Path)
				if err != nil {
					return err
				}
				if confirmed {
					printRemoveResult(options, remover.Remove(item.Path), item.Path)
				}
			}
		}
	}

	if options.DryRun {
		fmt.Fprintln(options.Out, "\nDry run complete. No files were actually deleted.")
	}
	fmt.Fprintln(options.Out, "\nDone!")
	return nil
}

func scan(ctx context.Context, options *Options, app catalog.App) ([]scanner.Item, error) {
	if options.Verbose {
		fmt.Fprintf(options.Out, "Bundle ID: %s\n", app.BundleID)
	}
	config := scanner.Config{
		HomeDir:          options.HomeDir,
		SystemLibraryDir: options.SystemLibraryDir,
	}
	return scanner.Collect(scanner.Default(config).Scan(ctx, app))
}

func printSelection(out io.Writer, app catalog.App, items []scanner.Item) {
	fmt.Fprintf(out, "Selected: %s\n", app.Name)
	fmt.Fprintf(out, "Location: %s\n", app.Path)
	terminal.PrintItems(out, items)
}

func printRemoveResult(options *Options, result cleanup.Result, label string) {
	if result.DryRun {
		fmt.Fprintf(options.Out, "Would delete: %s\n", result.Path)
		return
	}
	if result.Err != nil {
		fmt.Fprintf(options.Err, "Error deleting %s: %v\n", label, result.Err)
		return
	}
	fmt.Fprintf(options.Out, "Deleted: %s\n", result.Path)
}
