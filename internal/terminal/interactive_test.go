package terminal

import (
	"bytes"
	"strings"
	"testing"

	"github.com/jasonfriedland/zaap/internal/catalog"
)

func TestSelectAppReturnsChosenApp(t *testing.T) {
	apps := []catalog.App{
		{Name: "First", Path: "/Applications/First.app"},
		{Name: "Second", Path: "/Applications/Second.app"},
	}
	var out bytes.Buffer
	prompter := NewPrompter(strings.NewReader("2\n"), &out)

	app, ok, err := prompter.SelectApp(apps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected selection")
	}
	if app.Name != "Second" {
		t.Fatalf("expected Second, got %s", app.Name)
	}
	if !strings.Contains(out.String(), "Enter number:") {
		t.Fatalf("expected prompt output, got %q", out.String())
	}
}

func TestSelectAppExitAndInvalidInput(t *testing.T) {
	apps := []catalog.App{{Name: "First", Path: "/Applications/First.app"}}

	var exitOut bytes.Buffer
	_, ok, err := NewPrompter(strings.NewReader("0\n"), &exitOut).SelectApp(apps)
	if err != nil {
		t.Fatalf("unexpected exit error: %v", err)
	}
	if ok {
		t.Fatal("expected no selection for exit")
	}

	var invalidOut bytes.Buffer
	_, _, err = NewPrompter(strings.NewReader("nope\n"), &invalidOut).SelectApp(apps)
	if err == nil {
		t.Fatal("expected invalid input error")
	}
}

func TestConfirmAndAssociatedMode(t *testing.T) {
	var out bytes.Buffer
	prompter := NewPrompter(strings.NewReader("Y\nall\n"), &out)

	confirmed, err := prompter.Confirm("Delete? ")
	if err != nil {
		t.Fatalf("unexpected confirm error: %v", err)
	}
	if !confirmed {
		t.Fatal("expected confirmation")
	}

	mode, err := prompter.AssociatedMode()
	if err != nil {
		t.Fatalf("unexpected mode error: %v", err)
	}
	if mode != "all" {
		t.Fatalf("expected all, got %s", mode)
	}
}

func TestConfirmPathUsesBaseName(t *testing.T) {
	var out bytes.Buffer
	prompter := NewPrompter(strings.NewReader("n\n"), &out)

	confirmed, err := prompter.ConfirmPath("/tmp/example/file.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if confirmed {
		t.Fatal("expected negative confirmation")
	}
	if !strings.Contains(out.String(), "Delete file.txt?") {
		t.Fatalf("expected basename prompt, got %q", out.String())
	}
}
