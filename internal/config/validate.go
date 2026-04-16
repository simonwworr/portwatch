package config

import (
	"errors"
	"fmt"
	"net"
)

// ValidationError holds all validation issues found in a config.
type ValidationError struct {
	Errors []string
}

func (v *ValidationError) Error() string {
	if len(v.Errors) == 0 {
		return "validation failed"
	}
	msg := fmt.Sprintf("%d validation error(s):", len(v.Errors))
	for _, e := range v.Errors {
		msg += "\n  - " + e
	}
	return msg
}

// Validate checks a Config for common mistakes and returns a ValidationError
// if any issues are found, or nil if the config is valid.
func Validate(cfg *Config) error {
	var errs []string

	if len(cfg.Hosts) == 0 {
		errs = append(errs, "hosts list must not be empty")
	}

	seen := make(map[string]bool)
	for _, h := range cfg.Hosts {
		if h == "" {
			errs = append(errs, "host entry must not be empty")
			continue
		}
		if seen[h] {
			errs = append(errs, fmt.Sprintf("duplicate host: %s", h))
		}
		seen[h] = true
		if net.ParseIP(h) == nil {
			if _, err := net.LookupHost(h); err != nil {
				errs = append(errs, fmt.Sprintf("host %q does not resolve: %v", h, err))
			}
		}
	}

	if cfg.Interval <= 0 {
		errs = append(errs, "interval must be a positive duration")
	}

	if cfg.Timeout <= 0 {
		errs = append(errs, "timeout must be a positive duration")
	}

	if cfg.Timeout >= cfg.Interval {
		errs = append(errs, "timeout must be less than interval")
	}

	if len(errs) > 0 {
		return &ValidationError{Errors: errs}
	}
	return nil
}

// MustValidate is like Validate but panics on error — useful in tests.
func MustValidate(cfg *Config) {
	if err := Validate(cfg); err != nil {
		panic(errors.New("invalid config: " + err.Error()))
	}
}
