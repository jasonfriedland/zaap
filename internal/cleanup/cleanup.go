package cleanup

import "os"

// Result reports the outcome of a remove operation.
type Result struct {
	Path   string
	DryRun bool
	Err    error
}

// Remover deletes paths or reports dry-run actions.
type Remover struct {
	DryRun bool
}

// Remove deletes path unless dry-run mode is enabled.
func (r Remover) Remove(path string) Result {
	result := Result{Path: path, DryRun: r.DryRun}
	if r.DryRun {
		return result
	}
	result.Err = os.RemoveAll(path)
	return result
}
