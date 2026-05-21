package terminal

import (
	"fmt"
	"io"

	"github.com/jasonfriedland/zaap/internal/scanner"
)

func PrintItems(w io.Writer, items []scanner.Item) {
	grouped := GroupItems(items)
	for _, category := range scanner.CategoryOrder {
		paths := grouped[category]
		if len(paths) == 0 {
			continue
		}
		fmt.Fprintf(w, "\n%s:\n", category)
		for _, path := range paths {
			fmt.Fprintf(w, "  - %s\n", path)
		}
	}
}

func GroupItems(items []scanner.Item) map[scanner.Category][]string {
	grouped := make(map[scanner.Category][]string)
	for _, item := range items {
		grouped[item.Category] = append(grouped[item.Category], item.Path)
	}
	return grouped
}
