package terminal

import (
	"bytes"
	"strings"
	"testing"

	"github.com/jasonfriedland/zaap/internal/scanner"
)

func TestGroupItems(t *testing.T) {
	items := []scanner.Item{
		{Path: "/tmp/pref.plist", Category: scanner.CategoryAssociated},
		{Path: "/tmp/agent.plist", Category: scanner.CategoryStartupItem},
		{Path: "/tmp/cache", Category: scanner.CategoryAssociated},
	}

	grouped := GroupItems(items)
	if len(grouped[scanner.CategoryAssociated]) != 2 {
		t.Fatalf("expected 2 associated items, got %d", len(grouped[scanner.CategoryAssociated]))
	}
	if grouped[scanner.CategoryStartupItem][0] != "/tmp/agent.plist" {
		t.Fatalf("unexpected startup item: %v", grouped[scanner.CategoryStartupItem])
	}
}

func TestPrintItemsGroupsByCategoryOrder(t *testing.T) {
	items := []scanner.Item{
		{Path: "/tmp/agent.plist", Category: scanner.CategoryStartupItem},
		{Path: "/tmp/pref.plist", Category: scanner.CategoryAssociated},
	}

	var out bytes.Buffer
	PrintItems(&out, items)
	got := out.String()

	for _, want := range []string{
		"Associated files:",
		"  - /tmp/pref.plist",
		"Startup Items:",
		"  - /tmp/agent.plist",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in output:\n%s", want, got)
		}
	}

	associatedIndex := strings.Index(got, "Associated files:")
	startupIndex := strings.Index(got, "Startup Items:")
	if associatedIndex < 0 || startupIndex < 0 || associatedIndex > startupIndex {
		t.Fatalf("expected associated files before startup items:\n%s", got)
	}
}
