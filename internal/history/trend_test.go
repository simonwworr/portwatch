package history

import (
	"testing"
	"time"
)

func buildTrendStore() *Store {
	now := time.Now().UTC()
	s := &Store{}
	s.Append(Entry{Host: "h1", ScannedAt: now.Add(-2 * time.Hour), OpenPorts: []int{80, 443}})
	s.Append(Entry{Host: "h1", ScannedAt: now.Add(-1 * time.Hour), OpenPorts: []int{80, 8080}})
	s.Append(Entry{Host: "h1", ScannedAt: now, OpenPorts: []int{80}})
	return s
}

func TestTrend_Basic(t *testing.T) {
	s := buildTrendStore()
	ht := Trend(s, "h1")
	if ht.Host != "h1" {
		t.Fatalf("expected host h1, got %s", ht.Host)
	}
	if len(ht.Trends) != 3 {
		t.Fatalf("expected 3 port trends, got %d", len(ht.Trends))
	}
}

func TestTrend_Frequency(t *testing.T) {
	s := buildTrendStore()
	ht := Trend(s, "h1")
	for _, tr := range ht.Trends {
		if tr.Port == 80 {
			if tr.SeenCount != 3 {
				t.Errorf("port 80 seen count: want 3, got %d", tr.SeenCount)
			}
			if tr.Frequency != 1.0 {
				t.Errorf("port 80 frequency: want 1.0, got %f", tr.Frequency)
			}
		}
		if tr.Port == 443 {
			if tr.SeenCount != 1 {
				t.Errorf("port 443 seen count: want 1, got %d", tr.SeenCount)
			}
		}
	}
}

func TestTrend_NoEntries(t *testing.T) {
	s := &Store{}
	ht := Trend(s, "missing")
	if len(ht.Trends) != 0 {
		t.Errorf("expected empty trends, got %d", len(ht.Trends))
	}
}

func TestTrend_Sorted(t *testing.T) {
	s := buildTrendStore()
	ht := Trend(s, "h1")
	for i := 1; i < len(ht.Trends); i++ {
		if ht.Trends[i].Port < ht.Trends[i-1].Port {
			t.Errorf("trends not sorted at index %d", i)
		}
	}
}
