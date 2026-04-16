package history

import "time"

// PruneOptions controls how old entries are removed.
type PruneOptions struct {
	// MaxAge removes entries older than this duration. Zero means no age limit.
	MaxAge time.Duration
	// MaxEntries keeps at most this many entries per host (most recent). Zero means no limit.
	MaxEntries int
}

// Prune removes entries from the log according to opts and returns the number removed.
func Prune(l *Log, opts PruneOptions) int {
	before := len(l.Entries)
	if opts.MaxAge > 0 {
		cutoff := time.Now().UTC().Add(-opts.MaxAge)
		filtered := l.Entries[:0]
		for _, e := range l.Entries {
			if e.Timestamp.After(cutoff) {
				filtered = append(filtered, e)
			}
		}
		l.Entries = filtered
	}
	if opts.MaxEntries > 0 {
		byHost := map[string][]Entry{}
		for _, e := range l.Entries {
			byHost[e.Host] = append(byHost[e.Host], e)
		}
		l.Entries = l.Entries[:0]
		for _, entries := range byHost {
			if len(entries) > opts.MaxEntries {
				entries = entries[len(entries)-opts.MaxEntries:]
			}
			l.Entries = append(l.Entries, entries...)
		}
	}
	return before - len(l.Entries)
}
