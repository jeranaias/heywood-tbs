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
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("https://wttr.in/Quantico,VA?format=j1")
	if err != nil {
		return nil, fmt.Errorf("weather fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("weather fetch: status %d", resp.StatusCode)
	}

	var result struct {
		CurrentCondition []struct {
			TempF      string `json:"temp_F"`
			FeelsLikeF string `json:"FeelsLikeF"`
			WeatherDesc []struct {
				Value string `json:"value"`
			} `json:"weatherDesc"`
			WindspeedMiles string `json:"windspeedMiles"`
			Winddir16Point string `json:"winddir16Point"`
			Humidity       string `json:"humidity"`
		} `json:"current_condition"`
		Weather []struct {
			MaxTempF string `json:"maxtempF"`
			MinTempF string `json:"mintempF"`
		} `json:"weather"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("weather parse: %w", err)
	}

	if len(result.CurrentCondition) == 0 {
		return nil, fmt.Errorf("weather: no current conditions")
	}

	cc := result.CurrentCondition[0]
	wd := &WeatherData{
		TempF:      atoi(cc.TempF),
		FeelsLikeF: atoi(cc.FeelsLikeF),
		WindMPH:    atoi(cc.WindspeedMiles),
		WindDir:    cc.Winddir16Point,
		Humidity:   atoi(cc.Humidity),
		FetchedAt:  time.Now(),
	}
	if len(cc.WeatherDesc) > 0 {
		wd.Condition = cc.WeatherDesc[0].Value
	}
	if len(result.Weather) > 0 {
		wd.ForecastHi = atoi(result.Weather[0].MaxTempF)
		wd.ForecastLo = atoi(result.Weather[0].MinTempF)
	}

	return wd, nil
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
