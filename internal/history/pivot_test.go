package history

import (
	"testing"
	"time"
)

func buildPivotStore() Store {
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	s := NewMemoryStore()
	// host-a: day 1 ports 80,443 | day 2 ports 80
	s.Append(Entry{Host: "host-a", Time: base, Ports: []int{80, 443}})
	s.Append(Entry{Host: "host-a", Time: base.Add(24 * time.Hour), Ports: []int{80}})
	// host-b: day 1 ports 22 | day 2 ports 22,8080
	s.Append(Entry{Host: "host-b", Time: base, Ports: []int{22}})
	s.Append(Entry{Host: "host-b", Time: base.Add(24 * time.Hour), Ports: []int{22, 8080}})
	return s
}

func TestPivot_BucketLabels(t *testing.T) {
	result := Pivot(buildPivotStore(), "day")
	if len(result.BucketLabels) != 2 {
		t.Fatalf("expected 2 bucket labels, got %d", len(result.BucketLabels))
	}
	if result.BucketLabels[0] >= result.BucketLabels[1] {
		t.Error("bucket labels should be sorted ascending")
	}
}

func TestPivot_RowsPerHost(t *testing.T) {
	result := Pivot(buildPivotStore(), "day")
	if len(result.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(result.Rows))
	}
	if result.Rows[0].Host != "host-a" || result.Rows[1].Host != "host-b" {
		t.Error("rows should be sorted by host name")
	}
}

func TestPivot_PortsInBucket(t *testing.T) {
	result := Pivot(buildPivotStore(), "day")
	day1 := result.BucketLabels[0]
	portA := result.Rows[0].Buckets[day1]
	if len(portA) != 2 || portA[0] != 80 || portA[1] != 443 {
		t.Errorf("host-a day1 ports: expected [80 443], got %v", portA)
	}
}

func TestPivot_EmptyBucketReturnsNilOrEmpty(t *testing.T) {
	s := NewMemoryStore()
	s.Append(Entry{Host: "host-a", Time: time.Now(), Ports: []int{80}})
	result := Pivot(s, "day")
	if len(result.Rows) != 1 {
		t.Fatalf("expected 1 row")
	}
	// Only one bucket should exist.
	if len(result.BucketLabels) != 1 {
		t.Errorf("expected 1 bucket label, got %d", len(result.BucketLabels))
	}
}

func TestPivot_EmptyStore(t *testing.T) {
	result := Pivot(NewMemoryStore(), "day")
	if len(result.BucketLabels) != 0 || len(result.Rows) != 0 {
		t.Error("expected empty result for empty store")
	}
}

func TestPivot_MonthGranularity(t *testing.T) {
	base := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	s := NewMemoryStore()
	s.Append(Entry{Host: "host-a", Time: base, Ports: []int{80}})
	s.Append(Entry{Host: "host-a", Time: base.Add(10 * 24 * time.Hour), Ports: []int{443}})
	result := Pivot(s, "month")
	// Both entries fall in January 2024 → single bucket.
	if len(result.BucketLabels) != 1 {
		t.Errorf("expected 1 month bucket, got %d", len(result.BucketLabels))
	}
}
