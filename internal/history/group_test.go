package history

import (
	"testing"
	"time"
)

func buildGroupStore(t *testing.T) *Store {
	t.Helper()
	dir := t.TempDir()
	s, err := Load(dir)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	base := time.Date(2024, 3, 1, 10, 0, 0, 0, time.UTC)
	_ = s.Append("host1", base, []int{80, 443})
	_ = s.Append("host1", base.Add(24*time.Hour), []int{80, 443, 8080})
	_ = s.Append("host1", base.Add(8*24*time.Hour), []int{80})
	if err := Save(dir, s); err != nil {
		t.Fatalf("Save: %v", err)
	}
	return s
}

func TestGroupBy_Day(t *testing.T) {
	s := buildGroupStore(t)
	result, err := GroupBy(s, "host1", "day")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 3 {
		t.Fatalf("expected 3 day buckets, got %d", len(result))
	}
	if result[0].ScanCount != 1 {
		t.Errorf("expected scan count 1, got %d", result[0].ScanCount)
	}
}

func TestGroupBy_Month(t *testing.T) {
	s := buildGroupStore(t)
	result, err := GroupBy(s, "host1", "month")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 month bucket, got %d", len(result))
	}
	if result[0].Bucket != "2024-03" {
		t.Errorf("expected bucket 2024-03, got %s", result[0].Bucket)
	}
	if result[0].ScanCount != 3 {
		t.Errorf("expected 3 scans, got %d", result[0].ScanCount)
	}
}

func TestGroupBy_NoEntries(t *testing.T) {
	dir := t.TempDir()
	s, _ := Load(dir)
	result, err := GroupBy(s, "ghost", "day")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Errorf("expected nil result for missing host")
	}
}

func TestGroupBy_AvgPorts(t *testing.T) {
	dir := t.TempDir()
	s, _ := Load(dir)
	base := time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC)
	_ = s.Append("h", base, []int{80})
	_ = s.Append("h", base.Add(time.Hour), []int{80, 443})
	result, _ := GroupBy(s, "h", "day")
	if len(result) != 1 {
		t.Fatalf("expected 1 bucket")
	}
	if result[0].AvgPorts != 1.5 {
		t.Errorf("expected avg 1.5, got %f", result[0].AvgPorts)
	}
}
