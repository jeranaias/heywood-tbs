package data

import (
	"testing"
)

func TestAutoMapColumns_ExactMatch(t *testing.T) {
	headers := []string{"lastName", "firstName", "company", "rank", "phase"}
	mappings := AutoMapColumns(headers)

	want := map[string]string{
		"lastName":  "lastName",
		"firstName": "firstName",
		"company":   "company",
		"rank":      "rank",
		"phase":     "phase",
	}

	for i, m := range mappings {
		expected, ok := want[headers[i]]
		if !ok {
			continue
		}
		if m.FieldName != expected {
			t.Errorf("header %q: got FieldName=%q, want %q", headers[i], m.FieldName, expected)
		}
		if !m.AutoMatch {
			t.Errorf("header %q: AutoMatch should be true", headers[i])
		}
	}
}

func TestAutoMapColumns_CaseInsensitive(t *testing.T) {
	tests := []struct {
		header string
		field  string
	}{
		{"Last Name", "lastName"},
		{"FIRST NAME", "firstName"},
		{"Company", "company"},
		{"RANK", "rank"},
	}

	headers := make([]string, len(tests))
	for i, tc := range tests {
		headers[i] = tc.header
	}

	mappings := AutoMapColumns(headers)

	for i, tc := range tests {
		if mappings[i].FieldName != tc.field {
			t.Errorf("header %q: got FieldName=%q, want %q", tc.header, mappings[i].FieldName, tc.field)
		}
	}
}

func TestAutoMapColumns_FuzzyMatch(t *testing.T) {
	// "Student Last Name" contains "last name" (len >= 4), should match lastName
	headers := []string{"Student Last Name", "Student First Name"}
	mappings := AutoMapColumns(headers)

	if mappings[0].FieldName != "lastName" {
		t.Errorf("fuzzy: header %q got FieldName=%q, want lastName", headers[0], mappings[0].FieldName)
	}
	if mappings[1].FieldName != "firstName" {
		t.Errorf("fuzzy: header %q got FieldName=%q, want firstName", headers[1], mappings[1].FieldName)
	}
}

func TestAutoMapColumns_UnknownHeader(t *testing.T) {
	headers := []string{"FavoriteColor", "Shoe Size", "Zodiac Sign"}
	mappings := AutoMapColumns(headers)

	for _, m := range mappings {
		if m.FieldName != "" {
			t.Errorf("header %q should have empty FieldName, got %q", m.Header, m.FieldName)
		}
	}
}

func TestAutoMapColumns_NoDuplicates(t *testing.T) {
	// Both "Last Name" and "Surname" map to lastName — only the first should win
	headers := []string{"Last Name", "Surname"}
	mappings := AutoMapColumns(headers)

	fieldCount := make(map[string]int)
	for _, m := range mappings {
		if m.FieldName != "" {
			fieldCount[m.FieldName]++
		}
	}

	for field, count := range fieldCount {
		if count > 1 {
			t.Errorf("field %q was mapped %d times, expected at most once", field, count)
		}
	}

	// Verify at least one mapped to lastName
	mapped := false
	for _, m := range mappings {
		if m.FieldName == "lastName" {
			mapped = true
			break
		}
	}
	if !mapped {
		t.Error("expected at least one header to map to lastName")
	}
}
