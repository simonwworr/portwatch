package history

import "sort"

// MemoryStore is an in-memory Store used for testing.
type MemoryStore struct {
	entries []Entry
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (m *MemoryStore) Append(e Entry) {
	m.entries = append(m.entries, e)
}

func (m *MemoryStore) ForHost(host string) []Entry {
	var out []Entry
	for _, e := range m.entries {
		if e.Host == host {
			out = append(out, e)
		}
	}
	return out
}

func (m *MemoryStore) Hosts() []string {
	seen := map[string]struct{}{}
	for _, e := range m.entries {
		seen[e.Host] = struct{}{}
	}
	hosts := make([]string, 0, len(seen))
	for h := range seen {
		hosts = append(hosts, h)
	}
	sort.Strings(hosts)
	return hosts
}

func (m *MemoryStore) All() []Entry {
	return m.entries
}
