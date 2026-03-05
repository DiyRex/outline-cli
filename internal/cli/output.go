package cli

import (
	"fmt"
	"os"
	"strings"
)

type OutputManager struct {
	Quiet   bool
	Verbose bool
	NoColor bool
	YesMode bool
}

var Output = &OutputManager{}

func (o *OutputManager) Info(format string, args ...interface{}) {
	if o.Quiet {
		return
	}
	fmt.Fprintf(os.Stdout, format+"\n", args...)
}

func (o *OutputManager) Success(format string, args ...interface{}) {
	if o.Quiet {
		return
	}
	msg := fmt.Sprintf(format, args...)
	if o.NoColor {
		fmt.Fprintf(os.Stdout, "%s\n", msg)
	} else {
		fmt.Fprintf(os.Stdout, "\033[32m%s\033[0m\n", msg)
	}
}

func (o *OutputManager) Error(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if o.NoColor {
		fmt.Fprintf(os.Stderr, "Error: %s\n", msg)
	} else {
		fmt.Fprintf(os.Stderr, "\033[31mError: %s\033[0m\n", msg)
	}
}

func (o *OutputManager) Warn(format string, args ...interface{}) {
	if o.Quiet {
		return
	}
	msg := fmt.Sprintf(format, args...)
	if o.NoColor {
		fmt.Fprintf(os.Stderr, "Warning: %s\n", msg)
	} else {
		fmt.Fprintf(os.Stderr, "\033[33mWarning: %s\033[0m\n", msg)
	}
}

func (o *OutputManager) Debug(format string, args ...interface{}) {
	if !o.Verbose {
		return
	}
	msg := fmt.Sprintf(format, args...)
	if o.NoColor {
		fmt.Fprintf(os.Stderr, "[debug] %s\n", msg)
	} else {
		fmt.Fprintf(os.Stderr, "\033[90m[debug] %s\033[0m\n", msg)
	}
}

func (o *OutputManager) Progress(current, total int, label string) {
	if o.Quiet {
		return
	}
	pct := 0
	if total > 0 {
		pct = (current * 100) / total
	}
	barWidth := 30
	filled := (pct * barWidth) / 100
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
	fmt.Fprintf(os.Stderr, "\r  [%s] %d/%d %s", bar, current, total, label)
	if current >= total {
		fmt.Fprintln(os.Stderr)
	}
}

func (o *OutputManager) Confirm(prompt string) bool {
	if o.YesMode {
		return true
	}
	fmt.Fprintf(os.Stderr, "%s [y/N]: ", prompt)
	var response string
	fmt.Scanln(&response)
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}
