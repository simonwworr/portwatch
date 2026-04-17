package history

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"
)

// ExportFormat defines the output format for history export.
type ExportFormat string

const (
	FormatCSV  ExportFormat = "csv"
	FormatJSON ExportFormat = "json"
)

// ExportOptions controls what gets exported.
type ExportOptions struct {
	Host   string
	Since  time.Time
	Format ExportFormat
}

// Export writes history entries to w in the requested format.
func Export(store *Store, opts ExportOptions, w io.Writer) error {
	entries := store.ForHost(opts.Host)
	if !opts.Since.IsZero() {
		filtered := entries[:0]
		for _, e := range entries {
			if !e.ScannedAt.Before(opts.Since) {
				filtered = append(filtered, e)
			}
		}
		entries = filtered
	}

	switch opts.Format {
	case FormatCSV:
		return exportCSV(entries, w)
	case FormatJSON:
		return exportJSON(entries, w)
	default:
		return fmt.Errorf("unsupported format: %s", opts.Format)
	}
}

func exportCSV(entries []Entry, w io.Writer) error {
	cw := csv.NewWriter(w)
	_ = cw.Write([]string{"host", "scanned_at", "open_ports"})
	for _, e := range entries {
		ports := ""
		for i, p := range e.OpenPorts {
			if i > 0 {
				ports += ";"
			}
			ports += strconv.Itoa(p)
		}
		_ = cw.Write([]string{e.Host, e.ScannedAt.Format(time.RFC3339), ports})
	}
	cw.Flush()
	return cw.Error()
}

func exportJSON(entries []Entry, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}
