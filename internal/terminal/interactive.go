package terminal

import (
	"bufio"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jasonfriedland/zaap/internal/catalog"
)

type Prompter struct {
	In  *bufio.Reader
	Out io.Writer
}

func NewPrompter(in io.Reader, out io.Writer) Prompter {
	return Prompter{In: bufio.NewReader(in), Out: out}
}

func (p Prompter) SelectApp(apps []catalog.App) (catalog.App, bool, error) {
	fmt.Fprintln(p.Out, "Select an application to delete:")
	fmt.Fprintln(p.Out, "---------------------------------")
	for i, app := range apps {
		fmt.Fprintf(p.Out, "%d. %s\n", i+1, app.Name)
	}
	fmt.Fprintln(p.Out, "0. Exit")
	fmt.Fprint(p.Out, "\nEnter number: ")

	line, err := p.readLine()
	if err != nil {
		return catalog.App{}, false, err
	}
	selection, err := strconv.Atoi(line)
	if err != nil {
		return catalog.App{}, false, fmt.Errorf("invalid input")
	}
	if selection == 0 || selection < 0 || selection > len(apps) {
		return catalog.App{}, false, nil
	}
	return apps[selection-1], true, nil
}

func (p Prompter) Confirm(prompt string) (bool, error) {
	fmt.Fprint(p.Out, prompt)
	line, err := p.readLine()
	if err != nil {
		return false, err
	}
	return strings.ToLower(line) == "y", nil
}

func (p Prompter) AssociatedMode() (string, error) {
	fmt.Fprintln(p.Out, "\nDelete associated items? (y/n/all): ")
	return p.readLine()
}

func (p Prompter) ConfirmPath(path string) (bool, error) {
	return p.Confirm(fmt.Sprintf("Delete %s? (y/n): ", filepath.Base(path)))
}

func (p Prompter) readLine() (string, error) {
	line, err := p.In.ReadString('\n')
	if err != nil && len(line) == 0 {
		return "", err
	}
	return strings.TrimSpace(line), nil
}
