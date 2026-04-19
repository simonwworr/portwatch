package history

import (
	"testing"
	"time"
)

func buildFingerprintStore() Store {
	now := time.Now()
	s := NewMemoryStore()
	s.Append("host-a", Entry{Host: "host-a", Ports: []int{80, 443}, Timestamp: now})
	s.Append("host-b", Entry{Host: "host-b", Ports: []int{443, 80}, Timestamp: now})
	s.Append("host-c", Entry{Host: "host-c", Ports: []int{22, 8080}, Timestamp: now})
	s.Append("host-d", Entry{Host: "host-d", Ports: []int{}, Timestamp: now})
	return s
}

func TestFingerprints_Signature(t *testing.T) {
	s := buildFingerprintStore()
	fps := Fingerprints(s)

	sigMap := make(map[string]string)
	for _, f := range fps {
		sigMap[f.Host] = f.Signature
	}

	if sigMap["host-a"] != sigMap["host-b"] {
		t.Errorf("host-a and host-b should share signature, got %q vs %q", sigMap["host-a"], sigMap["host-b"])
	}
	if sigMap["host-a"] == sigMap["host-c"] {
		t.Errorf("host-a and host-c should differ")
	}
	if sigMap["host-d"] != "<empty>" {
		t.Errorf("expected <empty> for host-d, got %q", sigMap["host-d"])
	}
}

func TestFingerprints_SortedByHost(t *testing.T) {
	s := buildFingerprintStore()
	fps := Fingerprints(s)
	for i := 1; i < len(fps); i++ {
		if fps[i].Host < fps[i-1].Host {
			t.Errorf("results not sorted by host at index %d", i)
		}
	}
}

func TestGroupByFingerprint_SameSignature(t *testing.T) {
	s := buildFingerprintStore()
	groups := GroupByFingerprint(s)

	for sig, hosts := range groups {
		if sig == "80,443" || sig == "443,80" {
			if len(hosts) != 2 {
				t.Errorf("expected 2 hosts for signature %q, got %d", sig, len(hosts))
			}
		}
	}
}

func TestGroupByFingerprint_EmptyStore(t *testing.T) {
	s := NewMemoryStore()
	groups := GroupByFingerprint(s)
	if len(groups) != 0 {
		t.Errorf("expected empty groups, got %d", len(groups))
	}
}
