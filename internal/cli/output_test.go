package cli

import (
	"fmt"
	"testing"
)

func TestExitCodeFromAPIError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected int
	}{
		{"nil error", nil, ExitSuccess},
		{"auth error 401", fmt.Errorf("API error 401: unauthorized"), ExitAuthError},
		{"auth error 403", fmt.Errorf("API error 403: forbidden"), ExitAuthError},
		{"not found 404", fmt.Errorf("API error 404: not_found"), ExitNotFound},
		{"rate limit 429", fmt.Errorf("API error 429: rate limited"), ExitRateLimited},
		{"config error", fmt.Errorf("URL not configured"), ExitConfigError},
		{"general error", fmt.Errorf("something went wrong"), ExitGeneralError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := ExitCodeFromAPIError(tt.err)
			if code != tt.expected {
				t.Errorf("ExitCodeFromAPIError() = %d, want %d", code, tt.expected)
			}
		})
	}
}

func TestExitCodeDescription(t *testing.T) {
	tests := []struct {
		code     int
		expected string
	}{
		{ExitSuccess, "Success"},
		{ExitGeneralError, "General error"},
		{ExitConfigError, "Configuration error"},
		{ExitAuthError, "Authentication error"},
		{ExitNotFound, "Resource not found"},
		{ExitRateLimited, "Rate limited"},
		{ExitValidation, "Validation error"},
		{99, "Unknown exit code 99"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			desc := ExitCodeDescription(tt.code)
			if desc != tt.expected {
				t.Errorf("ExitCodeDescription(%d) = %q, want %q", tt.code, desc, tt.expected)
			}
		})
	}
}

func TestOutputManager_Confirm_YesMode(t *testing.T) {
	o := &OutputManager{YesMode: true}
	if !o.Confirm("test?") {
		t.Error("Confirm() should return true when YesMode is set")
	}
}

func TestOutputManager_Quiet(t *testing.T) {
	o := &OutputManager{Quiet: true}
	// These should not panic in quiet mode
	o.Info("test %s", "info")
	o.Success("test %s", "success")
	o.Warn("test %s", "warn")
}

func TestOutputManager_Verbose(t *testing.T) {
	o := &OutputManager{Verbose: false}
	// Debug should not output when verbose is false
	o.Debug("test %s", "debug")
}
