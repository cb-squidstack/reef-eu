package feeds

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WeatherData represents weather information
type WeatherData struct {
	Summary      string  `json:"summary"`
	TemperatureC float64 `json:"temperatureC"`
	FeelsLikeC   float64 `json:"feelsLikeC"`
}

// Coordinates represents latitude and longitude
type Coordinates struct {
	Lat float64
	Lon float64
}

// OpenMeteoResponse represents the API response from Open-Meteo
type OpenMeteoResponse struct {
	Current struct {
		Temperature         float64 `json:"temperature_2m"`
		ApparentTemperature float64 `json:"apparent_temperature"`
		WeatherCode         int     `json:"weather_code"`
	} `json:"current"`
}

// EU country coordinates (major cities)
var euCountryCoordinates = map[string]Coordinates{
	"GB": {Lat: 51.5074, Lon: -0.1278},  // London
	"FR": {Lat: 48.8566, Lon: 2.3522},   // Paris
	"DE": {Lat: 52.5200, Lon: 13.4050},  // Berlin
	"ES": {Lat: 40.4168, Lon: -3.7038},  // Madrid
	"IT": {Lat: 41.9028, Lon: 12.4964},  // Rome
	"NL": {Lat: 52.3676, Lon: 4.9041},   // Amsterdam
	"BE": {Lat: 50.8503, Lon: 4.3517},   // Brussels
	"SE": {Lat: 59.3293, Lon: 18.0686},  // Stockholm
	"NO": {Lat: 59.9139, Lon: 10.7522},  // Oslo
	"FI": {Lat: 60.1699, Lon: 24.9384},  // Helsinki
	"PL": {Lat: 52.2297, Lon: 21.0122},  // Warsaw
	"IE": {Lat: 53.3498, Lon: -6.2603},  // Dublin
	"PT": {Lat: 38.7223, Lon: -9.1393},  // Lisbon
	"AT": {Lat: 48.2082, Lon: 16.3738},  // Vienna
	"CH": {Lat: 46.9481, Lon: 7.4474},   // Bern
	"DK": {Lat: 55.6761, Lon: 12.5683},  // Copenhagen
	"CZ": {Lat: 50.0755, Lon: 14.4378},  // Prague
	"GR": {Lat: 37.9838, Lon: 23.7275},  // Athens
}

// Weather code to description mapping (WMO Weather interpretation codes)
var weatherCodeDescriptions = map[int]string{
	0:  "Clear sky",
	1:  "Mainly clear",
	2:  "Partly cloudy",
	3:  "Overcast",
	45: "Foggy",
	48: "Depositing rime fog",
	51: "Light drizzle",
	53: "Moderate drizzle",
	55: "Dense drizzle",
	61: "Slight rain",
	63: "Moderate rain",
	65: "Heavy rain",
	71: "Slight snow",
	73: "Moderate snow",
	75: "Heavy snow",
	77: "Snow grains",
	80: "Slight rain showers",
	81: "Moderate rain showers",
	82: "Violent rain showers",
	85: "Slight snow showers",
	86: "Heavy snow showers",
	95: "Thunderstorm",
	96: "Thunderstorm with slight hail",
	99: "Thunderstorm with heavy hail",
}

// FetchWeather fetches weather data for a given country using Open-Meteo API
func FetchWeather(country string) (*WeatherData, error) {
	// Get coordinates for the country
	coords, ok := euCountryCoordinates[country]
	if !ok {
		// Default to London if country not found
		coords = euCountryCoordinates["GB"]
	}

	// Build Open-Meteo API URL
	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f&current=temperature_2m,apparent_temperature,weather_code",
		coords.Lat, coords.Lon,
	)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Make API request
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("weather API call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather API returned status %d", resp.StatusCode)
	}

	// Parse response
	var apiResp OpenMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse weather response: %w", err)
	}

	// Convert weather code to description
	description, ok := weatherCodeDescriptions[apiResp.Current.WeatherCode]
	if !ok {
		description = "Unknown"
	}

	return &WeatherData{
		Summary:      description,
		TemperatureC: apiResp.Current.Temperature,
		FeelsLikeC:   apiResp.Current.ApparentTemperature,
	}, nil
}
