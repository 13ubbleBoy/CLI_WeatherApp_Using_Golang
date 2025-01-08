package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// converting the json data into go struct
type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"` // this is a tag

	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text`
		} `json:"condition"`
	} `json:"current"`

	Forecast struct {
		Forecastday []struct { // slice of struct
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text`
				} `json:"condition"`
				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {
	q := "Bengaluru"

	if len(os.Args) >= 2 {
		q = os.Args[1]
	}
	/*
		--> provide a city name at the time of execution, if we dont provide anything it will take Bengaluru by default.
			go run main.go Mumbai
			go run main.go London
	*/

	response, error := http.Get("http://api.weatherapi.com/v1/forecast.json?key=ec180872243c4f57a4f153631230105&q=" + q + "&days=1&aqi=no&alerts=no")
	if error != nil {
		panic(error)
	}

	// close the response body when main() function finishes its execution.
	defer response.Body.Close()

	if response.StatusCode != 200 {
		panic("Weather API not available.")
	}

	body, error := io.ReadAll(response.Body) // body has slice of bytes, to show it we have to convert it to a string.
	if error != nil {
		panic(error)
	}

	var weather Weather                    // weather is a variable of type Weather
	error = json.Unmarshal(body, &weather) // this will take 'body' and convert it to Weather.
	if error != nil {
		panic(error)
	}

	location, current, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour
	fmt.Printf("%s, %s: %.0f°C, %s\n", location.Name, location.Country, current.TempC, current.Condition.Text)

	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0)

		// show future weather not past.
		if date.Before(time.Now()) {
			continue
		}

		blue := "\033[34m"
		resetBlue := "\033[0m"
		yellow := "\033[33m"
		resetYellow := "\033[0m"

		// if temperature is less than 16°C then print it in blue
		if hour.TempC < 16 {
			fmt.Printf("%s - %s%.0f°C%s, %s\n", date.Format("15:04"), blue, hour.TempC, resetBlue, hour.Condition.Text)
		} else if hour.TempC > 24 {
			fmt.Printf("%s - %s%.0f°C%s, %s\n", date.Format("15:04"), yellow, hour.TempC, resetYellow, hour.Condition.Text)
		} else {
			fmt.Printf("%s - %.0f°C, %s\n", date.Format("15:04"), hour.TempC, hour.Condition.Text)
		}

	}
}
