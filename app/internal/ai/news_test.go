package ai

import (
	"strings"
	"testing"
	"time"
)

func TestFormatAge(t *testing.T) {
	tests := []struct {
		name     string
		offset   time.Duration
		wantSub  string // substring expected in result
	}{
		{"30 minutes ago", 30 * time.Minute, "30 min ago"},
		{"2 hours ago", 2 * time.Hour, "2h ago"},
		{"23 hours ago", 23 * time.Hour, "23h ago"},
		{"1 day ago", 25 * time.Hour, "1d ago"},
		{"3 days ago", 72 * time.Hour, "3d ago"},
		{"zero time", 0, "recent"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var input time.Time
			if tc.name == "zero time" {
				input = time.Time{} // zero value
			} else {
				input = time.Now().Add(-tc.offset)
			}
			got := formatAge(input)
			if !strings.Contains(got, tc.wantSub) {
				t.Errorf("formatAge() = %q, want substring %q", got, tc.wantSub)
			}
		})
	}
}

func TestFormatNewsForPrompt_Empty(t *testing.T) {
	got := FormatNewsForPrompt(nil)
	if !strings.Contains(got, "No recent") {
		t.Errorf("FormatNewsForPrompt(nil) = %q, expected mention of no news", got)
	}

	got2 := FormatNewsForPrompt([]NewsItem{})
	if !strings.Contains(got2, "No recent") {
		t.Errorf("FormatNewsForPrompt([]) = %q, expected mention of no news", got2)
	}
}

func TestFormatNewsForPrompt_WithItems(t *testing.T) {
	items := []NewsItem{
		{
			Title:     "Marines Deploy New Training System",
			Source:    "Military Times",
			Published: time.Now().Add(-2 * time.Hour),
			Link:      "https://example.com/1",
		},
		{
			Title:     "Quantico Base Expansion Approved",
			Source:    "Stars and Stripes",
			Published: time.Now().Add(-5 * time.Hour),
			Link:      "https://example.com/2",
		},
	}

	got := FormatNewsForPrompt(items)

	if !strings.Contains(got, "Headlines") {
		t.Error("expected result to contain 'Headlines' header")
	}
	if !strings.Contains(got, "Marines Deploy New Training System") {
		t.Error("expected result to contain first headline title")
	}
	if !strings.Contains(got, "Military Times") {
		t.Error("expected result to contain first headline source")
	}
	if !strings.Contains(got, "Quantico Base Expansion Approved") {
		t.Error("expected result to contain second headline title")
	}
	if !strings.Contains(got, "1.") || !strings.Contains(got, "2.") {
		t.Error("expected numbered list format")
	}
}
