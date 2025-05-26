/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
	"net/http"
	"net/url"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vibes",
	Short: "vibes tells you weather",
	Long: `vibes tells you about the vibes`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) { 
		replyGeneralWeather()
	},
}

type TempUnit string

const (
	Celsius TempUnit = "celsius"
	Fahrenheit TempUnit = "fahrenheit"
)

func getThresholds(unit TempUnit) (freezing, cold, cool, mild, tempChange float64) {
	if unit == Fahrenheit {
		return celsiusToFahrenheit(FreezingThresholdCelsius),
				celsiusToFahrenheit(ColdThresholdCelsius),
				celsiusToFahrenheit(CoolThresholdCelsius),
				celsiusToFahrenheit(MildThresholdCelsius),
				tempChangeCelsiusToFahrenheit(TempChangeThresholdCelsius)
	}
	return FreezingThresholdCelsius, ColdThresholdCelsius, CoolThresholdCelsius, MildThresholdCelsius, TempChangeThresholdCelsius
}

func celsiusToFahrenheit(c float64) float64 {
	return c * 9 / 5 + 32
}

func tempChangeCelsiusToFahrenheit(c float64) float64 {
	return c * 9 / 5
}

const (
    // Temperature categories (°C)
    FreezingThresholdCelsius     = 5.0
    ColdThresholdCelsius        = 10.0  
    CoolThresholdCelsius        = 15.0
    MildThresholdCelsius        = 20.0
	
	TempChangeThresholdCelsius   = 5.0

	// Rain chance %
	HighRainChance       = 60
    ModerateRainChance   = 30

	// API params
	ForecastHours = 4
	LocationTimeout = 10 * time.Second

	// API endpoints
	IPInfoURL = "https://ipinfo.io/json"
	OpenMeteoBaseURL = "https://api.open-meteo.com/v1/forecast"
	GeocodingBaseURL = "https://geocoding-api.open-meteo.com/v1/search"
)

type WeatherResponse struct {
	Hourly HourlyData `json:"hourly"`
}

type HourlyData struct {
	Time []string `json:"time"`
	Temperature []float64 `json:"temperature_2m"`
	RainChance []int `json:"precipitation_probability"`
	WindSpeed []float64 `json:"wind_speed_10m"`
}

type IPToCoordinates struct {
	City string 		`json:"city"`
	Region string		`json:"region"`
	Coordinates string  `json:"loc"`
}

type Location struct {
	City string
	Region string
	Lat float64
	Lon float64
}

var (
	httpClient = &http.Client {
		Timeout: LocationTimeout,
	}
	unitFlag string
)

func getLocation() (Location, error) {

	ipResponse, err := httpClient.Get(IPInfoURL)
	if err != nil {
		return Location{}, fmt.Errorf("failed to get IP info: %w", err)
	}
	defer ipResponse.Body.Close()

	if ipResponse.StatusCode != http.StatusOK {
		return Location{}, fmt.Errorf("location service returned status %d", ipResponse.StatusCode)
	}

	ipBody, err := io.ReadAll(ipResponse.Body)
	if err != nil {
		return Location{}, fmt.Errorf("failed to read JSON response: %w", err)
	}
	var IPData IPToCoordinates
	err = json.Unmarshal(ipBody, &IPData)
	if err != nil {
		return Location{}, fmt.Errorf("failed to unmarshal the IP coordinates: %w", err)
	}

	city := IPData.City
	region := IPData.Region
	coordinates := strings.Split(IPData.Coordinates, ",")
	if len(coordinates) != 2 {
		return Location{}, fmt.Errorf("error splitting coordinates, length != 2 (lat, lon), %s", IPData.Coordinates)
	}
	lat, err := strconv.ParseFloat(strings.TrimSpace(coordinates[0]), 64)
	if err != nil {
		return Location{}, fmt.Errorf("error converting latitude string to float: %w", err)
	}
	lon, err := strconv.ParseFloat(strings.TrimSpace(coordinates[1]), 64)
	if err != nil {
		return Location{}, fmt.Errorf("error converting longitude string to float: %w", err)
	}

	var location Location = Location {
		City: city,
		Region: region,
		Lat: lat,
		Lon: lon,
	}

	return location, nil
	
}

func getWeatherData(lat float64, lon float64, unit TempUnit) (WeatherResponse, error) {
	meteoApi := fmt.Sprintf("%s?latitude=%f&longitude=%f&hourly=temperature_2m,precipitation_probability,precipitation,rain,wind_speed_10m&forecast_hours=4&timezone=auto&temperature_unit=%s", OpenMeteoBaseURL, lat, lon, string(unit))
	response, err := httpClient.Get(meteoApi)
	if err != nil {
		return WeatherResponse{}, fmt.Errorf("error getting data from open-meteo: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return WeatherResponse{}, fmt.Errorf("error parsing the response: %w", err)
	}

	var weatherData WeatherResponse
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		return WeatherResponse{}, fmt.Errorf("error unmarshalling the response: %w", err)
	}

	return weatherData, nil
}

func analyzeWeather(location Location, weatherData WeatherResponse, unit TempUnit) {
	freezing, cold, cool, mild, tempChange := getThresholds(unit)
	city := location.City
	region := location.Region
	fmt.Printf("Here's the current weather condition report for %s, %s:\n", city, region)
	significantRise := false
	significantDrop := false
	currTemp := weatherData.Hourly.Temperature[0]
	
	// general temperature
	if currTemp <= freezing {
		fmt.Println("Freezing! Bundle up")
	} else if currTemp <= cold {
		fmt.Println("Proper jacket weather, maybe gloves")
	} else if currTemp <= cool {
		fmt.Println("Classic hoodie/light jacket zone")
	} else if currTemp <= mild {
		fmt.Println("Good weather, maybe just a light layer")
	} else {
		fmt.Println("T-shirt weather!")
	}
	
	minTemp := currTemp
	maxTemp := currTemp
	
	for _, temp := range weatherData.Hourly.Temperature {
		if temp < minTemp {
			minTemp = temp
		}
		if temp > maxTemp {
			maxTemp = temp
		}
	}
	
	// temp change: rise/drop 5°C
	if maxTemp - currTemp >= tempChange {
		significantRise = true
	}
	
	if currTemp - minTemp >= tempChange {
		significantDrop = true
	}

	unitSymbol := "°C"
	if unit == Fahrenheit {
		unitSymbol = "°F"
	}
	
	if significantDrop && significantRise {
		fmt.Println("temp will change significantly in next four hours")
	} else if significantDrop {
		fmt.Printf("temp will drop %.0f%s in the next four hours\n", tempChange, unitSymbol)
	} else if significantRise {
		fmt.Printf("it'll get %.0f%s hotter in the next four hours\n", tempChange, unitSymbol)
	} else {
		fmt.Println("temp will be around the same in the next four hours")
	}
	
	fmt.Printf("Current temp: %.1f%s\n", currTemp, unitSymbol)
	
	// rain chance
	maxRainChance := 0
	peakHour := 0
	for i, chance := range weatherData.Hourly.RainChance {
		if chance > maxRainChance {
			maxRainChance = chance
			peakHour = i
		}
	}
	if maxRainChance > HighRainChance {
		if peakHour == 0 {
			fmt.Printf("Definitely bring an umbrella! Very likely to rain right now\n")
		} else {
			fmt.Printf("Definitely bring an umbrella! Very likely to rain in %d hours \n", peakHour)
		}
	} else if maxRainChance > ModerateRainChance {
		if peakHour == 0 {
			fmt.Println("Might rain now - maybe keep a jacket handy.")
		} else {
		fmt.Printf("Might rain in %d hours - maybe keep a jacket handy. \n", peakHour)
		}
	} else {
		fmt.Println("No rain expected.")
	}
}

type LocationResult struct {
	Results []ResultData `json:"results"`
}

type ResultData struct {
	Lat float64	`json:"latitude"`
	Lon float64	`json:"longitude"`
	Region string `json:"admin1"`
}

func askUserForLocation() (Location, error) {

	fmt.Println("Please enter your city name (e.g., 'Seattle', 'Boston', 'Tokyo'):")
	// fmt.Scanln(&cityName)
	reader := bufio.NewReader(os.Stdin)
	cityName, _ := reader.ReadString('\n')
	cityName = strings.TrimSpace(cityName)

	queryURL := fmt.Sprintf("%s?name=%s", GeocodingBaseURL, url.QueryEscape(cityName))
	response, err := httpClient.Get(queryURL)
	if err != nil {
		return Location{}, fmt.Errorf("error getting the coordinates JSON: %w", err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return Location{}, fmt.Errorf("error reading the JSON body: %w", err)
	}
	var locationData LocationResult
	err = json.Unmarshal(body, &locationData)
	if err != nil {
		return Location{}, fmt.Errorf("error unmarshalling the JSON body: %w", err)
	}

	if len(locationData.Results) == 0 {
		return Location{}, fmt.Errorf("no results found for city: %s", cityName)
	}
	var location Location = Location {
		City: cityName,
		Region: locationData.Results[0].Region,
		Lat: locationData.Results[0].Lat,
		Lon: locationData.Results[0].Lon,
	}
	return location, nil

}

func replyGeneralWeather() {
	unit := Fahrenheit
	normalizedUnit := strings.ToLower(unitFlag)
	if normalizedUnit == "celsius" || normalizedUnit == "c" {
		unit = Celsius
	}

	location, err := getLocation()
	if err != nil {
		fmt.Println("error getting location automatically: ", err)
		fmt.Println("let's try this manually...")
		location, err = askUserForLocation()
		if err != nil {
			fmt.Println("error getting location manually: ", err)
			fmt.Println("Unable to determine location. Please try again later.")
			return
		}
	}
	weatherData, err := getWeatherData(location.Lat, location.Lon, unit)
	if err != nil {
		fmt.Println("error getting weatherData:", err)
		return
	}
	analyzeWeather(location, weatherData, unit)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.Flags().StringVarP(&unitFlag, "unit", "u", "fahrenheit", "Temperature unit: 'celsius'/'c' or 'fahrenheit'/'f'")
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.weather.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


