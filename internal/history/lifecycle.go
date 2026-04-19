package history

import "time"

// LifecycleEvent represents a port's open/close lifecycle span.
type LifecycleEvent struct {
	Host      string
	Port      int
	OpenedAt  time.Time
	ClosedAt  *time.Time // nil if still open
	Duration  time.Duration
}

// Lifecycle computes open/close lifecycle spans for each port on a host.
func Lifecycle(store Store, host string) []LifecycleEvent {
	entries := store.ForHost(host)
	if len(entries) == 0 {
		return nil
	}

	type span struct {
		openedAt time.Time
	}

	open := map[int]span{}
	var events []LifecycleEvent

	for i, e := range entries {
		current := toSet(e.Ports)

		if i == 0 {
			for p := range current {
				open[p] = span{openedAt: e.Time}
			}
			continue
		}

		prev := toSet(entries[i-1].Ports)

		// newly opened
		for p := range current {
			if !prev[p] {
				open[p] = span{openedAt: e.Time}
			}
		}

		// newly closed
		for p := range prev {
			if !current[p] {
				if s, ok := open[p]; ok {
					closed := e.Time
					events = append(events, LifecycleEvent{
						Host:     host,
						Port:     p,
						OpenedAt: s.openedAt,
						ClosedAt: &closed,
						Duration: closed.Sub(s.openedAt),
					})
					delete(open, p)
				}
			}
		}
	}

	// still-open ports
	for p, s := range open {
		events = append(events, LifecycleEvent{
			Host:     host,
			Port:     p,
			OpenedAt: s.openedAt,
		})
	}

	return events
}
