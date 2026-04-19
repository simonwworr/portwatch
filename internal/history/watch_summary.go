package history

import "sort"

// WatchSummary holds aggregated event counts for a host.
type WatchSummary struct {
	Host        string
	TotalEvents int
	Opened      int
	Closed      int
	TopPorts    []int
}

// SummarizeWatchLog returns per-host summaries from a WatchLog.
func SummarizeWatchLog(log *WatchLog) []WatchSummary {
	type portCount struct {
		port  int
		count int
	}

	type hostData struct {
		opened int
		closed int
		ports  map[int]int
	}

	hosts := map[string]*hostData{}

	for _, ev := range log.Events {
		hd, ok := hosts[ev.Host]
		if !ok {
			hd = &hostData{ports: map[int]int{}}
			hosts[ev.Host] = hd
		}
		switch ev.Type {
		case "opened":
			hd.opened++
		case "closed":
			hd.closed++
		}
		for _, p := range ev.Ports {
			hd.ports[p]++
		}
	}

	summaries := make([]WatchSummary, 0, len(hosts))
	for host, hd := range hosts {
		pcs := make([]portCount, 0, len(hd.ports))
		for p, c := range hd.ports {
			pcs = append(pcs, portCount{p, c})
		}
		sort.Slice(pcs, func(i, j int) bool {
			if pcs[i].count != pcs[j].count {
				return pcs[i].count > pcs[j].count
			}
			return pcs[i].port < pcs[j].port
		})
		top := []int{}
		for i, pc := range pcs {
			if i >= 5 {
				break
			}
			top = append(top, pc.port)
		}
		summaries = append(summaries, WatchSummary{
			Host:        host,
			TotalEvents: hd.opened + hd.closed,
			Opened:      hd.opened,
			Closed:      hd.closed,
			TopPorts:    top,
		})
	}
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].Host < summaries[j].Host
	})
	return summaries
}
