package practice

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"encoding/json"
	"time"
)

type WeatherResponse2 struct {
	Hourly []HourlyForecast2 `json:"hourly"`
}

type HourlyForecast2 struct {
	Dt int64 `json:"dt"`
	Temp float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	Humidity int `json:"humidity"`
	Pop float64 `json:"pop"`
	Weather []WeatherInfo2 `json:"weather"`
}

type WeatherInfo2 struct {
	Main string `json:"main"`
	Description string `json:"description"`
}

func GetWeatherData2() {
	apiKey := "fd1a9fe1168620101f0bcef74da06706"
	lat := 47.6038321
	lon := -122.330062
	basicUrl := fmt.Sprintf("https://api.openweathermap.org/data/3.0/onecall?lat=%f&lon=%f&units=metric&lang=en&exclude=current,minutely,daily,alerts&appid=%s", lat, lon, apiKey)
	
	response, err := http.Get(basicUrl)
	if err != nil {
		fmt.Println("Error getting from weather API", err)
		os.Exit(1)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error parsing the response: ", err)
		os.Exit(1)
	}
	var weatherData WeatherResponse2
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		fmt.Println("Error unmarshaling response body: ", err)
		os.Exit(1)
	}
	for i := 0; i < 5 && i < len(weatherData.Hourly); i++ {
		h := weatherData.Hourly[i]
		fmt.Printf("Time: %s | Temp: %.1f°C | Feels Like: %.1f°C | Humidity: %d%% | Condition: %s (%s) | Precip Chance: %.0f%%\n",
			unixToHour2(h.Dt),
			h.Temp,
			h.FeelsLike,
			h.Humidity,
			h.Weather[0].Main,
			h.Weather[0].Description,
			h.Pop*100,
		)
	}

	
}

func unixToHour2(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.Format("3:04 PM")
}