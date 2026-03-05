package changelog

import (
	"fmt"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"time"
)

type Options struct {
	From           string
	To             string
	IncludeAuthors bool
	RepoPath       string
}

type Commit struct {
	Hash    string
	Author  string
	Message string
	Type    string
	Scope   string
	Subject string
}

var conventionalRe = regexp.MustCompile(`^(\w+)(?:\(([^)]*)\))?!?:\s*(.+)$`)

func Generate(opts Options) (string, error) {
	if opts.RepoPath == "" {
		opts.RepoPath = "."
	}

	format := "%H|%an|%s"
	args := []string{"-C", opts.RepoPath, "log", "--pretty=format:" + format, fmt.Sprintf("%s..%s", opts.From, opts.To)}
	out, err := exec.Command("git", args...).Output()
	if err != nil {
		return "", fmt.Errorf("git log failed: %w (is %s..%s a valid range?)", err, opts.From, opts.To)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") {
		return "", fmt.Errorf("no commits found between %s and %s", opts.From, opts.To)
	}

	var commits []Commit
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 3)
		if len(parts) < 3 {
			continue
		}
		c := Commit{
			Hash:    parts[0][:8],
			Author:  parts[1],
			Message: parts[2],
		}
		matches := conventionalRe.FindStringSubmatch(c.Message)
		if len(matches) > 0 {
			c.Type = strings.ToLower(matches[1])
			c.Scope = matches[2]
			c.Subject = matches[3]
		} else {
			c.Type = "other"
			c.Subject = c.Message
		}
		commits = append(commits, c)
	}

	return formatMarkdown(commits, opts), nil
}

func formatMarkdown(commits []Commit, opts Options) string {
	categories := map[string][]Commit{}
	for _, c := range commits {
		categories[c.Type] = append(categories[c.Type], c)
	}

	categoryOrder := []struct {
		key   string
		title string
	}{
		{"feat", "Features"},
		{"fix", "Bug Fixes"},
		{"perf", "Performance"},
		{"refactor", "Refactoring"},
		{"docs", "Documentation"},
		{"test", "Tests"},
		{"ci", "CI/CD"},
		{"build", "Build"},
		{"chore", "Chores"},
		{"style", "Style"},
		{"other", "Other Changes"},
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Changelog %s..%s\n\n", opts.From, opts.To))
	sb.WriteString(fmt.Sprintf("*Generated on %s*\n\n", time.Now().Format("2006-01-02")))

	for _, cat := range categoryOrder {
		items, ok := categories[cat.key]
		if !ok || len(items) == 0 {
			continue
		}
		sb.WriteString(fmt.Sprintf("## %s\n\n", cat.title))

		sort.Slice(items, func(i, j int) bool {
			return items[i].Scope < items[j].Scope
		})

		for _, c := range items {
			line := "- "
			if c.Scope != "" {
				line += fmt.Sprintf("**%s:** ", c.Scope)
			}
			line += c.Subject
			if opts.IncludeAuthors {
				line += fmt.Sprintf(" (%s)", c.Author)
			}
			line += fmt.Sprintf(" (`%s`)", c.Hash)
			sb.WriteString(line + "\n")
		}
		sb.WriteString("\n")
		delete(categories, cat.key)
	}

	// Any remaining uncategorized types
	if len(categories) > 0 {
		var keys []string
		for k := range categories {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			sb.WriteString(fmt.Sprintf("## %s\n\n", titleCase(k)))
			for _, c := range categories[k] {
				line := fmt.Sprintf("- %s (`%s`)", c.Subject, c.Hash)
				if opts.IncludeAuthors {
					line = fmt.Sprintf("- %s (%s) (`%s`)", c.Subject, c.Author, c.Hash)
				}
				sb.WriteString(line + "\n")
			}
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func titleCase(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
