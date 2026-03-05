package cli

import (
	"fmt"
	"os"
	"strings"
)

const (
	ExitSuccess       = 0
	ExitGeneralError  = 1
	ExitConfigError   = 2
	ExitAuthError     = 3
	ExitNotFound      = 4
	ExitRateLimited   = 5
	ExitValidation    = 6
)

func ExitWithError(code int, err error) {
	Output.Error("%s", err.Error())
	os.Exit(code)
}

func ExitCodeFromAPIError(err error) int {
	if err == nil {
		return ExitSuccess
	}
	msg := err.Error()
	if strings.Contains(msg, "401") || strings.Contains(msg, "403") || strings.Contains(msg, "authentication") {
		return ExitAuthError
	}
	if strings.Contains(msg, "404") || strings.Contains(msg, "not_found") {
		return ExitNotFound
	}
	if strings.Contains(msg, "429") || strings.Contains(msg, "rate") {
		return ExitRateLimited
	}
	if strings.Contains(msg, "not configured") {
		return ExitConfigError
	}
	return ExitGeneralError
}

func ExitCodeDescription(code int) string {
	switch code {
	case ExitSuccess:
		return "Success"
	case ExitGeneralError:
		return "General error"
	case ExitConfigError:
		return "Configuration error"
	case ExitAuthError:
		return "Authentication error"
	case ExitNotFound:
		return "Resource not found"
	case ExitRateLimited:
		return "Rate limited"
	case ExitValidation:
		return "Validation error"
	default:
		return fmt.Sprintf("Unknown exit code %d", code)
	}
}
