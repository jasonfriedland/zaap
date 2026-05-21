package cleanup

import "os"

type Result struct {
	Path   string
	DryRun bool
	Err    error
}

type Remover struct {
	DryRun bool
}

func (r Remover) Remove(path string) Result {
	result := Result{Path: path, DryRun: r.DryRun}
	if r.DryRun {
		return result
	}
	result.Err = os.RemoveAll(path)
	return result
}
