package api

import (
	"encoding/json"
	"fmt"
	"os"
	"io"
	"net/http"
	"time"
)

type WeatherResponse struct {
	Hourly []HourlyForecast `json:"hourly"`
}

type HourlyForecast struct {
	Dt int64 `json:"dt"`
	Temp float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	Humidity int `json:"humidity"`
	Pop float64 `json:"pop"`
	Weather []WeatherInfo `json:"weather"`
}

type WeatherInfo struct {
	Main string `json:"main"`
	Description string `json:"description"`
}

func GetWeatherData() {
	apiKey := "fd1a9fe1168620101f0bcef74da06706"
	lat := 47.6038321
	lon := -122.330062
	basicURL := fmt.Sprintf("https://api.openweathermap.org/data/3.0/onecall?lat=%f&lon=%f&units=metric&lang=en&exclude=current,minutely,daily,alerts&appid=%s", lat, lon, apiKey)
	
	response, err := http.Get(basicURL)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading the response:", err)
		os.Exit(1)
	}
	
	var weatherData WeatherResponse
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		fmt.Println("Error unmarshalling the json: ", err)
		os.Exit(1)
	}

	fmt.Println("Next 5 hours forecast:")
	for i := 0; i < 5 && i < len(weatherData.Hourly); i++ {
		h := weatherData.Hourly[i]
		fmt.Printf("Time: %s | Temp: %.1f°C | Feels Like: %.1f°C | Humidity: %d%% | Condition: %s (%s) | Precip Chance: %.0f%%\n", 
		unixToHour(h.Dt),
		h.Temp,
		h.FeelsLike,
		h.Humidity,
		h.Weather[0].Main,
		h.Weather[0].Description,
		h.Pop*100,
	)
	}
	// fmt.Println(string(body))

	// var prettyJSON map[string]interface{}
	// json.Unmarshal(body, &prettyJSON)
	// formatted, _ := json.MarshalIndent(prettyJSON, "", "  ")
	// fmt.Println(string(formatted))
}


func unixToHour(timeStamp int64) string {
	t := time.Unix(timeStamp, 0)
	return t.Format("3:04 PM")
}