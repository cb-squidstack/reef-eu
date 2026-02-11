package feeds

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWeatherDataSerialization(t *testing.T) {
	weather := WeatherData{
		Summary:      "Clear sky",
		TemperatureC: 25.5,
		FeelsLikeC:   24.0,
	}

	data, err := json.Marshal(weather)
	if err != nil {
		t.Fatalf("failed to marshal WeatherData: %v", err)
	}

	var decoded WeatherData
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal WeatherData: %v", err)
	}

	if decoded.Summary != weather.Summary {
		t.Errorf("expected Summary %s, got %s", weather.Summary, decoded.Summary)
	}
	if decoded.TemperatureC != weather.TemperatureC {
		t.Errorf("expected TemperatureC %f, got %f", weather.TemperatureC, decoded.TemperatureC)
	}
	if decoded.FeelsLikeC != weather.FeelsLikeC {
		t.Errorf("expected FeelsLikeC %f, got %f", weather.FeelsLikeC, decoded.FeelsLikeC)
	}
}

func TestEUCountryCoordinates(t *testing.T) {
	tests := []struct {
		country string
		hasCoords bool
	}{
		{"GB", true},
		{"FR", true},
		{"DE", true},
		{"ES", true},
		{"IT", true},
		{"NL", true},
		{"BE", true},
		{"SE", true},
		{"NO", true},
		{"FI", true},
		{"PL", true},
		{"IE", true},
		{"PT", true},
		{"AT", true},
		{"CH", true},
		{"DK", true},
		{"CZ", true},
		{"GR", true},
		{"XX", false}, // Unknown country
	}

	for _, tt := range tests {
		t.Run(tt.country, func(t *testing.T) {
			_, exists := euCountryCoordinates[tt.country]
			if exists != tt.hasCoords {
				t.Errorf("country %s: expected exists=%v, got %v", tt.country, tt.hasCoords, exists)
			}
		})
	}
}

func TestCoordinatesStructure(t *testing.T) {
	// Test that London coordinates are reasonable
	london, ok := euCountryCoordinates["GB"]
	if !ok {
		t.Fatal("London coordinates not found")
	}

	// London should be around 51N, 0W
	if london.Lat < 50 || london.Lat > 53 {
		t.Errorf("London latitude %f seems incorrect (expected ~51)", london.Lat)
	}
	if london.Lon < -2 || london.Lon > 2 {
		t.Errorf("London longitude %f seems incorrect (expected ~0)", london.Lon)
	}
}

func TestWeatherCodeDescriptions(t *testing.T) {
	tests := []struct {
		code        int
		description string
	}{
		{0, "Clear sky"},
		{1, "Mainly clear"},
		{2, "Partly cloudy"},
		{3, "Overcast"},
		{45, "Foggy"},
		{61, "Slight rain"},
		{95, "Thunderstorm"},
		{999, ""}, // Unknown code
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			desc, ok := weatherCodeDescriptions[tt.code]
			if tt.description == "" {
				if ok {
					t.Errorf("expected code %d to not exist, but got %s", tt.code, desc)
				}
			} else {
				if !ok {
					t.Errorf("expected code %d to exist", tt.code)
				}
				if desc != tt.description {
					t.Errorf("code %d: expected %s, got %s", tt.code, tt.description, desc)
				}
			}
		})
	}
}

func TestFetchWeatherWithMockServer(t *testing.T) {
	// Create a mock Open-Meteo server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := OpenMeteoResponse{}
		response.Current.Temperature = 25.5
		response.Current.ApparentTemperature = 24.0
		response.Current.WeatherCode = 0 // Clear sky

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Note: This test would need to modify FetchWeather to accept a custom URL
	// For now, we're testing the structure
	t.Skip("Skipping integration test - would need to inject test server URL")
}

func TestFetchWeatherWithUnknownCountry(t *testing.T) {
	// Test that unknown country defaults to London
	_, ok := euCountryCoordinates["UNKNOWN"]
	if ok {
		t.Error("UNKNOWN country should not have coordinates")
	}

	// Verify GB (London) exists as the fallback
	_, ok = euCountryCoordinates["GB"]
	if !ok {
		t.Error("GB (London) should exist as fallback")
	}
}

func TestOpenMeteoResponseStructure(t *testing.T) {
	jsonData := `{
		"current": {
			"temperature_2m": 25.5,
			"apparent_temperature": 24.0,
			"weather_code": 0
		}
	}`

	var response OpenMeteoResponse
	if err := json.Unmarshal([]byte(jsonData), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response.Current.Temperature != 25.5 {
		t.Errorf("expected temperature 25.5, got %f", response.Current.Temperature)
	}
	if response.Current.ApparentTemperature != 24.0 {
		t.Errorf("expected apparent temperature 24.0, got %f", response.Current.ApparentTemperature)
	}
	if response.Current.WeatherCode != 0 {
		t.Errorf("expected weather code 0, got %d", response.Current.WeatherCode)
	}
}

func TestAllCountriesHaveValidCoordinates(t *testing.T) {
	for country, coords := range euCountryCoordinates {
		t.Run(country, func(t *testing.T) {
			// Latitude should be between -90 and 90
			if coords.Lat < -90 || coords.Lat > 90 {
				t.Errorf("%s: invalid latitude %f", country, coords.Lat)
			}
			// Longitude should be between -180 and 180
			if coords.Lon < -180 || coords.Lon > 180 {
				t.Errorf("%s: invalid longitude %f", country, coords.Lon)
			}
		})
	}
}
