package ai

import (
	"strings"
	"testing"
)

func TestWmoCodeToCondition(t *testing.T) {
	tests := []struct {
		code int
		want string
	}{
		{0, "Clear sky"},
		{1, "Partly cloudy"},
		{3, "Partly cloudy"},
		{45, "Fog"},
		{51, "Drizzle"},
		{61, "Rain"},
		{71, "Snow"},
		{80, "Rain showers"},
		{85, "Snow showers"},
		{95, "Thunderstorm"},
		{99, "Thunderstorm"},
		{100, "Unknown"},
	}

	for _, tc := range tests {
		got := wmoCodeToCondition(tc.code)
		if got != tc.want {
			t.Errorf("wmoCodeToCondition(%d) = %q, want %q", tc.code, got, tc.want)
		}
	}
}

func TestDegreesToCardinal(t *testing.T) {
	tests := []struct {
		deg  float64
		want string
	}{
		{0, "N"},
		{90, "E"},
		{180, "S"},
		{270, "W"},
		{45, "NE"},
		{135, "SE"},
		{225, "SW"},
		{315, "NW"},
		{22, "NNE"},
		{350, "N"},
	}

	for _, tc := range tests {
		got := degreesToCardinal(tc.deg)
		if got != tc.want {
			t.Errorf("degreesToCardinal(%.0f) = %q, want %q", tc.deg, got, tc.want)
		}
	}
}

func TestUniformRecommendation(t *testing.T) {
	tests := []struct {
		name    string
		tempF   int
		wantAny []string // result must contain at least one of these substrings
	}{
		{
			name:    "freezing cold",
			tempF:   20,
			wantAny: []string{"Cold weather", "cold weather", "warming"},
		},
		{
			name:    "cold",
			tempF:   30,
			wantAny: []string{"Cold weather", "cold weather", "warming"},
		},
		{
			name:    "cool",
			tempF:   45,
			wantAny: []string{"warming layer", "Gloves"},
		},
		{
			name:    "comfortable",
			tempF:   65,
			wantAny: []string{"Standard", "No special"},
		},
		{
			name:    "warm",
			tempF:   80,
			wantAny: []string{"Heat", "water", "heat"},
		},
		{
			name:    "very hot",
			tempF:   95,
			wantAny: []string{"Heat Category III", "hydration", "WBGT"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := uniformRecommendation(tc.tempF, "Clear sky")
			found := false
			for _, substr := range tc.wantAny {
				if strings.Contains(got, substr) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("uniformRecommendation(%d) = %q, expected to contain one of %v",
					tc.tempF, got, tc.wantAny)
			}
		})
	}
}
