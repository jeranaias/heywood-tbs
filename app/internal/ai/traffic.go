package ai

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"heywood-tbs/internal/models"
)

// TBS main campus coordinates (The Basic School, MCB Quantico)
const (
	tbsLat = 38.5230
	tbsLon = -77.3186
)

// TrafficService calculates real route data for off-base appointments.
type TrafficService struct{}

// RouteInfo holds calculated route data for an off-base appointment.
type RouteInfo struct {
	Appointment    models.XOScheduleItem
	DistanceMiles  float64
	BaseTravelMins int    // from OSRM (no traffic)
	EstTravelMins  int    // adjusted for conditions
	DepartBy       string // recommended departure time (HHMM)
	Conditions     string // traffic condition summary
}

// CalculateRoutes computes travel info for off-base appointments.
// Uses OSRM (free, no API key) for real route calculations.
// Adjusts estimates based on weather conditions and time of day.
func (t *TrafficService) CalculateRoutes(appointments []models.XOScheduleItem, weather *WeatherData) []RouteInfo {
	var routes []RouteInfo

	for _, appt := range appointments {
		if appt.OnBase || appt.Latitude == 0 {
			continue
		}

		ri := RouteInfo{Appointment: appt}

		// Fetch real route from OSRM
		dist, durSec, err := fetchOSRMRoute(tbsLat, tbsLon, appt.Latitude, appt.Longitude)
		if err != nil {
			// Fallback: estimate from straight-line distance
			dist = haversine(tbsLat, tbsLon, appt.Latitude, appt.Longitude)
			durSec = dist / 40 * 3600 // assume 40 mph avg
		}

		ri.DistanceMiles = dist
		baseMins := int(math.Ceil(durSec / 60))
		ri.BaseTravelMins = baseMins

		// Adjust for conditions
		multiplier := 1.0
		var conditions []string

		// Time-of-day adjustment based on appointment time
		hour := parseHour(appt.StartTime)
		if hour >= 7 && hour <= 9 {
			multiplier += 0.3
			conditions = append(conditions, "morning rush hour on I-95/Route 1")
		} else if hour >= 16 && hour <= 18 {
			multiplier += 0.35
			conditions = append(conditions, "evening rush hour on I-95/Route 1")
		}

		// Weather adjustment
		if weather != nil {
			switch {
			case strings.Contains(strings.ToLower(weather.Condition), "snow") ||
				strings.Contains(strings.ToLower(weather.Condition), "ice"):
				multiplier += 0.5
				conditions = append(conditions, "snow/ice — reduced speeds likely")
			case strings.Contains(strings.ToLower(weather.Condition), "rain") ||
				strings.Contains(strings.ToLower(weather.Condition), "drizzle") ||
				strings.Contains(strings.ToLower(weather.Condition), "shower"):
				multiplier += 0.2
				conditions = append(conditions, "wet roads — allow extra stopping distance")
			case strings.Contains(strings.ToLower(weather.Condition), "fog"):
				multiplier += 0.25
				conditions = append(conditions, "fog — reduced visibility")
			}
		}

		// I-95 corridor buffer for Pentagon/Arlington trips (>25 miles)
		if dist > 25 {
			multiplier += 0.15
			conditions = append(conditions, "I-95 corridor — variable congestion")
		}

		// Gate traffic buffer
		conditions = append(conditions, "allow extra time for base gate security")

		ri.EstTravelMins = int(math.Ceil(float64(baseMins) * multiplier))

		// Add 10-min buffer for parking/walking
		totalMins := ri.EstTravelMins + 10

		if len(conditions) == 0 {
			ri.Conditions = "Normal conditions"
		} else {
			ri.Conditions = strings.Join(conditions, "; ")
		}

		// Calculate NLT departure time (round down to nearest 5 min for conservatism)
		ri.DepartBy = subtractMinutesNLT(appt.StartTime, totalMins)

		routes = append(routes, ri)
	}

	return routes
}

// fetchOSRMRoute gets real route data from the free OSRM demo server.
func fetchOSRMRoute(fromLat, fromLon, toLat, toLon float64) (distMiles float64, durSeconds float64, err error) {
	url := fmt.Sprintf(
		"https://router.project-osrm.org/route/v1/driving/%.6f,%.6f;%.6f,%.6f?overview=false",
		fromLon, fromLat, toLon, toLat,
	)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return 0, 0, fmt.Errorf("OSRM fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, 0, fmt.Errorf("OSRM: status %d", resp.StatusCode)
	}

	var result struct {
		Code   string `json:"code"`
		Routes []struct {
			Distance float64 `json:"distance"` // meters
			Duration float64 `json:"duration"` // seconds
		} `json:"routes"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, 0, fmt.Errorf("OSRM parse: %w", err)
	}

	if result.Code != "Ok" || len(result.Routes) == 0 {
		return 0, 0, fmt.Errorf("OSRM: no route found")
	}

	route := result.Routes[0]
	return route.Distance * 0.000621371, route.Duration, nil // meters to miles
}

// haversine calculates straight-line distance in miles between two coordinates.
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 3959.0 // Earth radius in miles
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

func parseHour(hhmm string) int {
	if len(hhmm) < 2 {
		return 0
	}
	h := 0
	for _, c := range hhmm[:2] {
		h = h*10 + int(c-'0')
	}
	return h
}

// subtractMinutesNLT subtracts minutes and rounds DOWN to the nearest 5-min mark.
func subtractMinutesNLT(hhmm string, mins int) string {
	if len(hhmm) < 4 {
		return hhmm
	}
	h := parseHour(hhmm)
	m := int(hhmm[2]-'0')*10 + int(hhmm[3]-'0')

	totalMins := h*60 + m - mins
	if totalMins < 0 {
		totalMins += 24 * 60
	}

	// Round down to nearest 5 minutes (more conservative)
	totalMins = (totalMins / 5) * 5

	return fmt.Sprintf("%02d%02d", totalMins/60, totalMins%60)
}

// FormatTrafficForPrompt formats route info for injection into the system prompt.
func FormatTrafficForPrompt(routes []RouteInfo) string {
	if len(routes) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString("Travel Advisory for Today's Off-Base Appointments:\n")
	for _, r := range routes {
		fmt.Fprintf(&b, "- %s–%s: %s\n", r.Appointment.StartTime, r.Appointment.EndTime, r.Appointment.Title)
		fmt.Fprintf(&b, "  Location: %s\n", r.Appointment.Location)
		fmt.Fprintf(&b, "  Distance: %.1f miles | Base travel: %d min | Est. with conditions: %d min\n",
			r.DistanceMiles, r.BaseTravelMins, r.EstTravelMins)
		fmt.Fprintf(&b, "  Depart NLT: %s\n", r.DepartBy)
		fmt.Fprintf(&b, "  Conditions: %s\n", r.Conditions)
	}
	return b.String()
}
