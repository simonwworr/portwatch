package history

import "time"

// FilterOptions defines criteria for filtering history entries.
type FilterOptions struct {
	Host  string
	Port  int
	Since time.Time
	Until time.Time
}

// Filter returns entries from the store matching all non-zero criteria.
func Filter(store *Store, opts FilterOptions) []Entry {
	var out []Entry

	hosts := store.Hosts()
	for _, host := range hosts {
		if opts.Host != "" && host != opts.Host {
			continue
		}
		entries := store.ForHost(host)
		for _, e := range entries {
			if !opts.Since.IsZero() && e.Time.Before(opts.Since) {
				continue
			}
			if !opts.Until.IsZero() && e.Time.After(opts.Until) {
				continue
			}
			if opts.Port != 0 && !containsPort(e.Ports, opts.Port) {
				continue
			}
			out = append(out, e)
		}
	}
	return out
}
