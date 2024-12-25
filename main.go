package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
)

// Weather represents the structure of weather data obtained from the API
type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceOfRain int64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {
	// Default location is Ghaziabad; accepts input from command line arguments
	q := "Ghaziabad"
	if len(os.Args) >= 2 {
		q = os.Args[1]
	}

	// Fetch weather data from the API
	res, err := http.Get("http://api.weatherapi.com/v1/forecast.json?key=ca2218a7adc747e399d101230242703&q=" + q + "&days=1&aqi=no&alerts=no")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	// Check if API response is successful
	if res.StatusCode != 200 {
		panic("weather API is not working")
	}

	// Read API response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	// Unmarshal JSON response into Weather struct
	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}

	// Extract relevant data from weather struct
	location, current, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour

	// Print location and current weather condition
	fmt.Printf("%s, %s: %.0f°C, %s\n", location.Name, location.Country, current.TempC, current.Condition.Text)
	fmt.Print("\n")

	// Print table header for hourly forecast
	fmt.Println("Time - Temp - Rain% - Condition")
	fmt.Print("\n")

	// Print hourly forecast with color coding for rain percentage
	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0)

		// Skip past hours
		if date.Before(time.Now()) {
			continue
		}

		// Construct message for the hourly forecast
		message := fmt.Sprintf(
			"%s - %.0f°C - %d%% - %s\n",
			date.Format("15:00"),
			hour.TempC,
			hour.ChanceOfRain,
			hour.Condition.Text,
		)

		// Print message with color coding for rain percentage
		if hour.ChanceOfRain < 40 {
			color.Red(message)
		} else {
			fmt.Print(message)
		}
	}
}
