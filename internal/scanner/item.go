package scanner

// Category identifies the kind of related file found.
type Category string

const (
	// CategoryAssociated covers preferences, caches, logs, and app support.
	CategoryAssociated Category = "Associated files"
	// CategoryControlPane covers preference pane bundles.
	CategoryControlPane Category = "Control Panels"
	// CategoryStartupItem covers launch agents and daemons.
	CategoryStartupItem Category = "Startup Items"
	// CategoryQuickLook covers QuickLook generators.
	CategoryQuickLook Category = "QuickLook Plugins"
	// CategoryScreenSaver covers screen saver bundles.
	CategoryScreenSaver Category = "Screen Savers"
	// CategoryInputMethod covers input method bundles.
	CategoryInputMethod Category = "Input Methods"
	// CategoryFont covers font files.
	CategoryFont Category = "Fonts"
)

// CategoryOrder defines the display order for grouped scanner results.
var CategoryOrder = []Category{
	CategoryAssociated,
	CategoryControlPane,
	CategoryStartupItem,
	CategoryQuickLook,
	CategoryScreenSaver,
	CategoryInputMethod,
	CategoryFont,
}

// Item is one file or directory found by a scanner.
type Item struct {
	Path        string
	Category    Category
	MatchReason string
}
