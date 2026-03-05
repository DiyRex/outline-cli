package changelog

import (
	"strings"
	"testing"
)

func TestFormatMarkdown_ConventionalCommits(t *testing.T) {
	commits := []Commit{
		{Hash: "abc12345", Author: "Alice", Message: "feat(auth): add login endpoint", Type: "feat", Scope: "auth", Subject: "add login endpoint"},
		{Hash: "def67890", Author: "Bob", Message: "fix: resolve null pointer", Type: "fix", Scope: "", Subject: "resolve null pointer"},
		{Hash: "ghi11111", Author: "Charlie", Message: "docs: update README", Type: "docs", Scope: "", Subject: "update README"},
		{Hash: "jkl22222", Author: "Dave", Message: "chore: bump dependencies", Type: "chore", Scope: "", Subject: "bump dependencies"},
	}

	opts := Options{From: "v1.0", To: "v1.1", IncludeAuthors: false}
	md := formatMarkdown(commits, opts)

	// Check structure
	if !strings.Contains(md, "# Changelog v1.0..v1.1") {
		t.Error("Should contain changelog header")
	}
	if !strings.Contains(md, "## Features") {
		t.Error("Should contain Features section")
	}
	if !strings.Contains(md, "## Bug Fixes") {
		t.Error("Should contain Bug Fixes section")
	}
	if !strings.Contains(md, "## Documentation") {
		t.Error("Should contain Documentation section")
	}
	if !strings.Contains(md, "## Chores") {
		t.Error("Should contain Chores section")
	}
	if !strings.Contains(md, "**auth:**") {
		t.Error("Should contain scope prefix for scoped commits")
	}
	if !strings.Contains(md, "`abc1234") {
		t.Error("Should contain commit hash")
	}
}

func TestFormatMarkdown_WithAuthors(t *testing.T) {
	commits := []Commit{
		{Hash: "abc12345", Author: "Alice", Message: "feat: new feature", Type: "feat", Subject: "new feature"},
	}

	opts := Options{From: "v1.0", To: "v1.1", IncludeAuthors: true}
	md := formatMarkdown(commits, opts)

	if !strings.Contains(md, "(Alice)") {
		t.Error("Should include author name when IncludeAuthors is true")
	}
}

func TestFormatMarkdown_WithoutAuthors(t *testing.T) {
	commits := []Commit{
		{Hash: "abc12345", Author: "Alice", Message: "feat: new feature", Type: "feat", Subject: "new feature"},
	}

	opts := Options{From: "v1.0", To: "v1.1", IncludeAuthors: false}
	md := formatMarkdown(commits, opts)

	if strings.Contains(md, "(Alice)") {
		t.Error("Should not include author name when IncludeAuthors is false")
	}
}

func TestFormatMarkdown_UnknownTypes(t *testing.T) {
	commits := []Commit{
		{Hash: "abc12345", Author: "Alice", Message: "custom: something special", Type: "custom", Subject: "something special"},
	}

	opts := Options{From: "v1.0", To: "v1.1"}
	md := formatMarkdown(commits, opts)

	if !strings.Contains(md, "## Custom") {
		t.Error("Should capitalize unknown commit types")
	}
}

func TestFormatMarkdown_EmptyCommits(t *testing.T) {
	opts := Options{From: "v1.0", To: "v1.1"}
	md := formatMarkdown(nil, opts)

	if !strings.Contains(md, "# Changelog v1.0..v1.1") {
		t.Error("Should still contain header even with no commits")
	}
}

func TestGenerate_InvalidRange(t *testing.T) {
	_, err := Generate(Options{
		From:     "nonexistent-tag-abc",
		To:       "nonexistent-tag-xyz",
		RepoPath: ".",
	})
	if err == nil {
		t.Error("Expected error for invalid git range")
	}
}

func TestConventionalRegex(t *testing.T) {
	tests := []struct {
		input string
		match bool
		ctype string
		scope string
		subj  string
	}{
		{"feat: add feature", true, "feat", "", "add feature"},
		{"fix(api): resolve bug", true, "fix", "api", "resolve bug"},
		{"docs: update README", true, "docs", "", "update README"},
		{"feat!: breaking change", true, "feat", "", "breaking change"},
		{"feat(scope)!: breaking", true, "feat", "scope", "breaking"},
		{"not a conventional commit", false, "", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			matches := conventionalRe.FindStringSubmatch(tt.input)
			if tt.match {
				if len(matches) == 0 {
					t.Errorf("Expected match for %q", tt.input)
					return
				}
				if matches[1] != tt.ctype {
					t.Errorf("Type = %q, want %q", matches[1], tt.ctype)
				}
				if matches[2] != tt.scope {
					t.Errorf("Scope = %q, want %q", matches[2], tt.scope)
				}
				if matches[3] != tt.subj {
					t.Errorf("Subject = %q, want %q", matches[3], tt.subj)
				}
			} else {
				if len(matches) > 0 {
					t.Errorf("Expected no match for %q", tt.input)
				}
			}
		})
	}
}
