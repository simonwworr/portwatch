package history

// PinDeviation describes a port that was expected (pinned) but is now missing,
// or a port that is open but not pinned for a host.
type PinDeviation struct {
	Host    string
	Missing []int // pinned but not currently open
	Extra   []int // open but not pinned
}

// CheckPins compares current open ports against pinned ports for a host.
// It returns a PinDeviation if any discrepancy is found.
func CheckPins(store *PinStore, host string, openPorts []int) *PinDeviation {
	pins := store.ForHost(host)
	if len(pins) == 0 {
		return nil
	}

	pinned := make(map[int]bool)
	for _, p := range pins {
		for _, port := range p.Ports {
			pinned[port] = true
		}
	}

	open := make(map[int]bool)
	for _, p := range openPorts {
		open[p] = true
	}

	var missing, extra []int
	for port := range pinned {
		if !open[port] {
			missing = append(missing, port)
		}
	}
	for port := range open {
		if !pinned[port] {
			extra = append(extra, port)
		}
	}

	if len(missing) == 0 && len(extra) == 0 {
		return nil
	}
	return &PinDeviation{Host: host, Missing: missing, Extra: extra}
}
