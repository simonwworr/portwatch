package history

import (
	"os"
	"path/filepath"
	"time"
)

// CleanupOptions controls how old history files are removed.
type CleanupOptions struct {
	// Dir is the directory containing history files.
	Dir string
	// MaxAge removes files older than this duration. Zero means no age limit.
	MaxAge time.Duration
	// DryRun reports what would be deleted without deleting.
	DryRun bool
}

// CleanupResult summarises what was (or would be) removed.
type CleanupResult struct {
	Removed []string
	Errors  []error
}

// Cleanup removes history files that exceed the age threshold.
// Only files matching the pattern "*.json" inside Dir are considered.
func Cleanup(opts CleanupOptions) (CleanupResult, error) {
	if opts.Dir == "" {
		return CleanupResult{}, nil
	}

	entries, err := os.ReadDir(opts.Dir)
	if err != nil {
		if os.IsNotExist(err) {
			return CleanupResult{}, nil
		}
		return CleanupResult{}, err
	}

	var result CleanupResult
	now := time.Now()

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if filepath.Ext(e.Name()) != ".json" {
			continue
		}

		if opts.MaxAge > 0 {
			info, err := e.Info()
			if err != nil {
				result.Errors = append(result.Errors, err)
				continue
			}
			if now.Sub(info.ModTime()) < opts.MaxAge {
				continue
			}
		}

		full := filepath.Join(opts.Dir, e.Name())
		result.Removed = append(result.Removed, full)

		if !opts.DryRun {
			if err := os.Remove(full); err != nil {
				result.Errors = append(result.Errors, err)
			}
		}
	}

	return result, nil
}
