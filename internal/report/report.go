package report

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/user/portwatch/internal/state"
)

// Format represents the output format for reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Report holds a summary of port changes for one or more hosts.
type Report struct {
	GeneratedAt time.Time
	Entries     []Entry
}

// Entry represents the diff result for a single host.
type Entry struct {
	Host   string
	Opened []int
	Closed []int
}

// New builds a Report from a map of host diffs.
func New(diffs map[string]state.Diff) *Report {
	r := &Report{GeneratedAt: time.Now()}
	for host, d := range diffs {
		r.Entries = append(r.Entries, Entry{
			Host:   host,
			Opened: d.Opened,
			Closed: d.Closed,
		})
	}
	return r
}

// Print writes the report to w in the given format.
func (r *Report) Print(w io.Writer, format Format) {
	if w == nil {
		w = os.Stdout
	}
	switch format {
	case FormatJSON:
		r.printJSON(w)
	default:
		r.printText(w)
	}
}

func (r *Report) printText(w io.Writer) {
	fmt.Fprintf(w, "Port Watch Report — %s\n", r.GeneratedAt.Format(time.RFC3339))
	fmt.Fprintln(w, strings.Repeat("-", 40))
	if len(r.Entries) == 0 {
		fmt.Fprintln(w, "No changes detected.")
		return
	}
	for _, e := range r.Entries {
		fmt.Fprintf(w, "Host: %s\n", e.Host)
		if len(e.Opened) > 0 {
			fmt.Fprintf(w, "  Opened: %v\n", e.Opened)
		}
		if len(e.Closed) > 0 {
			fmt.Fprintf(w, "  Closed: %v\n", e.Closed)
		}
	}
}

func (r *Report) printJSON(w io.Writer) {
	fmt.Fprintf(w, `{"generated_at":%q,"entries":[`, r.GeneratedAt.Format(time.RFC3339))
	for i, e := range r.Entries {
		if i > 0 {
			fmt.Fprint(w, ",")
		}
		fmt.Fprintf(w, `{"host":%q,"opened":%s,"closed":%s}`,
			e.Host, intSliceJSON(e.Opened), intSliceJSON(e.Closed))
	}
	fmt.Fprintln(w, "]}")
}

func intSliceJSON(s []int) string {
	if len(s) == 0 {
		return "[]"
	}
	parts := make([]string, len(s))
	for i, v := range s {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return "[" + strings.Join(parts, ",") + "]"
}
