package config

import (
	"testing"
	"time"
)

func validConfig() *Config {
	return &Config{
		Hosts:    []string{"127.0.0.1"},
		Interval: 30 * time.Second,
		Timeout:  5 * time.Second,
	}
}

func TestValidate_ValidConfig(t *testing.T) {
	if err := Validate(validConfig()); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidate_NoHosts(t *testing.T) {
	cfg := validConfig()
	cfg.Hosts = nil
	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for empty hosts")
	}
}

func TestValidate_DuplicateHost(t *testing.T) {
	cfg := validConfig()
	cfg.Hosts = []string{"127.0.0.1", "127.0.0.1"}
	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for duplicate host")
	}
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Errors) == 0 {
		t.Fatal("expected at least one validation error")
	}
}

func TestValidate_ZeroInterval(t *testing.T) {
	cfg := validConfig()
	cfg.Interval = 0
	if err := Validate(cfg); err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestValidate_TimeoutExceedsInterval(t *testing.T) {
	cfg := validConfig()
	cfg.Timeout = 60 * time.Second
	cfg.Interval = 30 * time.Second
	if err := Validate(cfg); err == nil {
		t.Fatal("expected error when timeout >= interval")
	}
}

func TestValidate_EmptyHostEntry(t *testing.T) {
	cfg := validConfig()
	cfg.Hosts = []string{""}
	if err := Validate(cfg); err == nil {
		t.Fatal("expected error for empty host string")
	}
}

func TestValidationError_Message(t *testing.T) {
	ve := &ValidationError{Errors: []string{"bad host", "bad interval"}}
	msg := ve.Error()
	if msg == "" {
		t.Fatal("expected non-empty error message")
	}
}
