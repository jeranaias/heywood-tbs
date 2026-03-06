package ai

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// WeatherService fetches and caches weather data for Quantico, VA.
type WeatherService struct {
	mu      sync.Mutex
	cached  *WeatherData
	expires time.Time
}

// WeatherData holds parsed weather information.
type WeatherData struct {
	TempF      int
	FeelsLikeF int
	Condition  string
	WindMPH    int
	WindDir    string
	Humidity   int
	ForecastHi int
	ForecastLo int
	FetchedAt  time.Time
}

// Get returns cached weather data, fetching from wttr.in if expired.
func (w *WeatherService) Get() (*WeatherData, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.cached != nil && time.Now().Before(w.expires) {
		return w.cached, nil
	}

	wd, err := fetchWeather()
	if err != nil {
		return nil, err
	}

	w.cached = wd
	w.expires = time.Now().Add(30 * time.Minute)
	return wd, nil
}

func fetchWeather() (*WeatherData, error) {
	// Open-Meteo: free, no API key, works from any datacenter
	// Quantico, VA: 38.52°N, 77.29°W
	const url = "https://api.open-meteo.com/v1/forecast?" +
		"latitude=38.52&longitude=-77.29" +
		"&current=temperature_2m,relative_humidity_2m,apparent_temperature,weather_code,wind_speed_10m,wind_direction_10m" +
		"&daily=temperature_2m_max,temperature_2m_min" +
		"&temperature_unit=fahrenheit&wind_speed_unit=mph&timezone=America%2FNew_York" +
		"&forecast_days=1"

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("weather fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("weather fetch: status %d", resp.StatusCode)
	}

	var result struct {
		Current struct {
			Temp       float64 `json:"temperature_2m"`
			FeelsLike  float64 `json:"apparent_temperature"`
			Humidity   int     `json:"relative_humidity_2m"`
			WeatherCode int    `json:"weather_code"`
			WindSpeed  float64 `json:"wind_speed_10m"`
			WindDir    float64 `json:"wind_direction_10m"`
		} `json:"current"`
		Daily struct {
			TempMax []float64 `json:"temperature_2m_max"`
			TempMin []float64 `json:"temperature_2m_min"`
		} `json:"daily"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("weather parse: %w", err)
	}

	wd := &WeatherData{
		TempF:      int(result.Current.Temp),
		FeelsLikeF: int(result.Current.FeelsLike),
		Condition:  wmoCodeToCondition(result.Current.WeatherCode),
		WindMPH:    int(result.Current.WindSpeed),
		WindDir:    degreesToCardinal(result.Current.WindDir),
		Humidity:   result.Current.Humidity,
		FetchedAt:  time.Now(),
	}
	if len(result.Daily.TempMax) > 0 {
		wd.ForecastHi = int(result.Daily.TempMax[0])
	}
	if len(result.Daily.TempMin) > 0 {
		wd.ForecastLo = int(result.Daily.TempMin[0])
	}

	return wd, nil
}

// wmoCodeToCondition converts WMO weather codes to human-readable conditions.
func wmoCodeToCondition(code int) string {
	switch {
	case code == 0:
		return "Clear sky"
	case code <= 3:
		return "Partly cloudy"
	case code <= 49:
		return "Fog"
	case code <= 59:
		return "Drizzle"
	case code <= 69:
		return "Rain"
	case code <= 79:
		return "Snow"
	case code <= 82:
		return "Rain showers"
	case code <= 86:
		return "Snow showers"
	case code <= 99:
		return "Thunderstorm"
	default:
		return "Unknown"
	}
}

// degreesToCardinal converts wind direction degrees to 16-point compass.
func degreesToCardinal(deg float64) string {
	dirs := []string{"N", "NNE", "NE", "ENE", "E", "ESE", "SE", "SSE",
		"S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW"}
	idx := int((deg + 11.25) / 22.5) % 16
	return dirs[idx]
}

// FormatWeatherForPrompt formats weather data for injection into a system prompt.
func FormatWeatherForPrompt(wd *WeatherData) string {
	uniform := uniformRecommendation(wd.TempF, wd.Condition)
	return fmt.Sprintf(
		"Current Weather at Quantico, VA (as of %s):\n"+
			"Temperature: %d°F (feels like %d°F)\n"+
			"Conditions: %s\n"+
			"Wind: %s %d mph\n"+
			"Humidity: %d%%\n"+
			"Today's High/Low: %d°F / %d°F\n"+
			"Uniform Recommendation: %s",
		wd.FetchedAt.Format("1504"),
		wd.TempF, wd.FeelsLikeF,
		wd.Condition,
		wd.WindDir, wd.WindMPH,
		wd.Humidity,
		wd.ForecastHi, wd.ForecastLo,
		uniform,
	)
}

func uniformRecommendation(tempF int, condition string) string {
	switch {
	case tempF < 32:
		return "Cold weather gear mandatory. Warming stations required for extended outdoor training. Gloves and warming layers for all hands."
	case tempF < 50:
		return "Cammies with warming layer. Gloves recommended for extended outdoor periods. Cold weather PT gear authorized."
	case tempF < 70:
		return "Standard cammies. No special weather considerations."
	case tempF < 85:
		return "Heat Category I-II likely. Ensure water availability. Monitor for heat casualties during strenuous activity."
	default:
		return "Heat Category III+ likely. Modified training schedule recommended. Mandatory hydration plan. WBGT monitoring required."
	}
}

func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
