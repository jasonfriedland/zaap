package scanner

type Category string

const (
	CategoryAssociated  Category = "Associated files"
	CategoryControlPane Category = "Control Panels"
	CategoryStartupItem Category = "Startup Items"
	CategoryQuickLook   Category = "QuickLook Plugins"
	CategoryScreenSaver Category = "Screen Savers"
	CategoryInputMethod Category = "Input Methods"
	CategoryFont        Category = "Fonts"
)

var CategoryOrder = []Category{
	CategoryAssociated,
	CategoryControlPane,
	CategoryStartupItem,
	CategoryQuickLook,
	CategoryScreenSaver,
	CategoryInputMethod,
	CategoryFont,
}

type Item struct {
	Path        string
	Category    Category
	MatchReason string
}
