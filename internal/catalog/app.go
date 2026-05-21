package catalog

// App describes an installed macOS application bundle.
type App struct {
	Name     string
	Path     string
	BundleID string
}
