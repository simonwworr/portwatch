package history

import (
	"testing"
	"time"
)

func buildClusterStore(hosts map[string][][]int) Store {
	s := Store{}
	now := time.Now()
	for host, snapshots := range hosts {
		for i, ports := range snapshots {
			s[host] = append(s[host], Entry{
				Timestamp: now.Add(time.Duration(i) * time.Minute),
				Ports:     ports,
			})
		}
	}
	return s
}

func TestCluster_GroupsIdenticalFingerprints(t *testing.T) {
	store := buildClusterStore(map[string][][]int{
		"host-a": {{80, 443}},
		"host-b": {{80, 443}},
		"host-c": {{22}},
	})
	clusters := Cluster(store)
	if len(clusters) != 2 {
		t.Fatalf("expected 2 clusters, got %d", len(clusters))
	}
}

func TestCluster_HostsInSameCluster(t *testing.T) {
	store := buildClusterStore(map[string][][]int{
		"alpha": {{80, 443}},
		"beta":  {{80, 443}},
	})
	clusters := Cluster(store)
	if len(clusters) != 1 {
		t.Fatalf("expected 1 cluster, got %d", len(clusters))
	}
	if len(clusters[0].Hosts) != 2 {
		t.Errorf("expected 2 hosts in cluster, got %d", len(clusters[0].Hosts))
	}
}

func TestCluster_EmptyStore(t *testing.T) {
	clusters := Cluster(Store{})
	if len(clusters) != 0 {
		t.Errorf("expected empty result, got %d", len(clusters))
	}
}

func TestCluster_UsesLatestSnapshot(t *testing.T) {
	now := time.Now()
	store := Store{
		"host-x": {
			{Timestamp: now.Add(-2 * time.Minute), Ports: []int{22}},
			{Timestamp: now, Ports: []int{80, 443}},
		},
		"host-y": {
			{Timestamp: now, Ports: []int{80, 443}},
		},
	}
	clusters := Cluster(store)
	if len(clusters) != 1 {
		t.Fatalf("expected 1 cluster (both use latest 80,443), got %d", len(clusters))
	}
}
