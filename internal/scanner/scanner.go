package scanner

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// PortResult holds the result of a single port scan.
type PortResult struct {
	Host string
	Port int
	Open bool
}

// ScanHost scans a range of TCP ports on the given host concurrently.
// workers controls the level of concurrency.
func ScanHost(host string, startPort, endPort, workers int, timeout time.Duration) []PortResult {
	ports := make(chan int, workers)
	results := make(chan PortResult, endPort-startPort+1)
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for port := range ports {
				address := fmt.Sprintf("%s:%d", host, port)
				conn, err := net.DialTimeout("tcp", address, timeout)
				open := err == nil
				if open {
					conn.Close()
				}
				results <- PortResult{Host: host, Port: port, Open: open}
			}
		}()
	}

	go func() {
		for p := startPort; p <= endPort; p++ {
			ports <- p
		}
		close(ports)
		wg.Wait()
		close(results)
	}()

	var out []PortResult
	for r := range results {
		if r.Open {
			out = append(out, r)
		}
	}
	return out
}

// OpenPorts returns just the list of open port numbers from results.
func OpenPorts(results []PortResult) []int {
	ports := make([]int, 0, len(results))
	for _, r := range results {
		ports = append(ports, r.Port)
	}
	return ports
}
