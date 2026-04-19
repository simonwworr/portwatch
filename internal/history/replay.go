package history

import "time"

// ReplayEvent represents a single port state at a point in time.
type ReplayEvent struct {
	Host      string
	Ports     []int
	ScannedAt time.Time
}

// ReplayOptions controls how a replay is filtered.
type ReplayOptions struct {
	Host  string
	Since time.Time
	Until time.Time
}

// Replay returns ordered scan entries for a host within the given time window,
// allowing callers to "replay" port history step by step.
func Replay(store Store, opts ReplayOptions) ([]ReplayEvent, error) {
	entries, err := store.ForHost(opts.Host)
	if err != nil {
		return nil, err
	}

	var events []ReplayEvent
	for _, e := range entries {
		if !opts.Since.IsZero() && e.ScannedAt.Before(opts.Since) {
			continue
		}
		if !opts.Until.IsZero() && e.ScannedAt.After(opts.Until) {
			continue
		}
		events = append(events, ReplayEvent{
			Host:      e.Host,
			Ports:     e.Ports,
			ScannedAt: e.ScannedAt,
		})
	}
	return events, nil
}
