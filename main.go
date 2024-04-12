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
		ForecastDay []struct {
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {

	region := "kef,tunisia"

	if len(os.Args) >= 2 {
		region = os.Args[1]
	}

	fmt.Println("Starting ")

	res, err := http.Get("http://api.weatherapi.com/v1/forecast.json?key=135ea5d37ffa43c391c155108241204&q=" + region + "&days=1&aqi=no&alerts=no")

	if err != nil { // check if there was an error
		panic(err)
	}
	defer res.Body.Close() // close the body request when done

	if res.StatusCode != 200 {
		panic("WEATHER API NOT AVAILABLE")
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}

	// fmt.Println(weather)

	Location, current, hours := weather.Location, weather.Current, weather.Forecast.ForecastDay[0].Hour
	fmt.Println("********************************")
	fmt.Printf("%s , %s : %0.0fC ,%s \n", Location.Name, Location.Country, current.TempC, current.Condition.Text)
	fmt.Println("********************************")

	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0)
		if date.Before(time.Now()) {
			continue
		}
		message := fmt.Sprintf(
			"%s - %.0fC, %.0f%%, %s\n",
			date.Format("15:04"),
			hour.TempC,
			hour.ChanceOfRain,
			hour.Condition.Text,
		)
		if hour.ChanceOfRain < 40 {
			fmt.Print(message)
		} else {
			color.Red(message)
		}
	}
}
