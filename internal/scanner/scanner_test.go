package scanner_test

import (
	"net"
	"testing"
	"time"

	"github.com/portwatch/internal/scanner"
)

// startTestServer opens a TCP listener on a random port and returns the port + a close func.
func startTestServer(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	return port, func() { ln.Close() }
}

func TestScanHost_DetectsOpenPort(t *testing.T) {
	port, stop := startTestServer(t)
	defer stop()

	results := scanner.ScanHost("127.0.0.1", port, port, 5, time.Second)
	if len(results) != 1 {
		t.Fatalf("expected 1 open port, got %d", len(results))
	}
	if results[0].Port != port {
		t.Errorf("expected port %d, got %d", port, results[0].Port)
	}
	if !results[0].Open {
		t.Error("expected port to be open")
	}
}

func TestScanHost_ClosedPort(t *testing.T) {
	// Port 1 is almost certainly closed in a test environment.
	results := scanner.ScanHost("127.0.0.1", 1, 1, 2, 200*time.Millisecond)
	if len(results) != 0 {
		t.Errorf("expected no open ports, got %d", len(results))
	}
}

func TestOpenPorts(t *testing.T) {
	input := []scanner.PortResult{
		{Host: "localhost", Port: 80, Open: true},
		{Host: "localhost", Port: 443, Open: true},
	}
	ports := scanner.OpenPorts(input)
	if len(ports) != 2 || ports[0] != 80 || ports[1] != 443 {
		t.Errorf("unexpected ports: %v", ports)
	}
}
