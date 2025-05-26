/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	// "math"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vibes",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) { 
		replyGeneralWeather()
	},
}

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

func replyGeneralWeather() {
	ip_response, err := http.Get("https://ipinfo.io/json")
	if err != nil {
		fmt.Println("couldn't extract IP", err)
		os.Exit(1)
	}
	defer ip_response.Body.Close()
	ip_body, err := io.ReadAll(ip_response.Body)
	if err != nil {
		fmt.Println("Error paring the IP json, ", err)
		os.Exit(1)
	}
	var IPData IPToCoordinates
	err = json.Unmarshal(ip_body, &IPData)
	if err != nil {
		fmt.Println("Error unmarshalling the ip coordinates", err)
		os.Exit(1)
	}

	city := IPData.City
	region := IPData.Region

	fmt.Printf("Here's the current weather condition report for %s, %s:\n", city, region)

	coordinates := strings.Split(IPData.Coordinates, ",")
	lat, err := strconv.ParseFloat(coordinates[0], 64)
	if err != nil {
		fmt.Println("Error converting string to float, ", err)
		os.Exit(1)
	}
	lon, err := strconv.ParseFloat(coordinates[1], 64)
		if err != nil {
		fmt.Println("Error converting string to float, ", err)
		os.Exit(1)
	}

	meteo_api := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&hourly=temperature_2m,precipitation_probability,precipitation,rain,wind_speed_10m&forecast_hours=4&timezone=auto", lat, lon)
	response, err := http.Get(meteo_api)
	if err != nil {
		fmt.Println("Error getting data from Open-Meteo, ", err)
		os.Exit(1)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error parsing the response, ", err)
		os.Exit(1)
	}

	var weatherData WeatherResponse
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		fmt.Println("Error unmarshalling body into weather data, ", err)
		os.Exit(1)
	}

	// temp rise/drop 5°C
	fiveDegreesRaise := false
	fiveDegreesDrop := false
	currTemp := weatherData.Hourly.Temperature[0]

	if currTemp <= 5.0 {
		fmt.Println("Freezing! Bundle up")
	} else if currTemp <= 10.0 {
		fmt.Println("Proper jacket weather, maybe gloves")
	} else if currTemp <= 15.0 {
		fmt.Println("Classic hoodie/light jacket zone")
	} else if currTemp <= 20.0 {
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

	if maxTemp - currTemp >= 5.0 {
		fiveDegreesRaise = true
	}

	if currTemp - minTemp >= 5.0 {
		fiveDegreesDrop = true
	}

	if fiveDegreesDrop && fiveDegreesRaise {
		fmt.Println("temp will change significantly in next four hours")
	} else if fiveDegreesDrop {
		fmt.Println("temp will drop 5°C in the next four hours")
	} else if fiveDegreesRaise {
		fmt.Println("it'll get 5°C hotter in the next four hours")
	} else {
		fmt.Println("temp will be around the same in the next four hours")
	}

	fmt.Printf("Current temp: %.1f°C\n", weatherData.Hourly.Temperature[0])
	
	// rain chance
	maxRainChance := 0
	peakHour := 0
	for i, chance := range weatherData.Hourly.RainChance {
		if chance > maxRainChance {
			maxRainChance = chance
			peakHour = i
		}
	}
	if maxRainChance > 60 {
		if peakHour == 0 {
    		fmt.Printf("Definitely bring an umbrella! Very likely to rain right now\n")
		} else {
			fmt.Printf("Definitely bring an umbrella! Very likely to rain in %d hours \n", peakHour)
		}
	} else if maxRainChance > 30 {
		if peakHour == 0 {
			fmt.Println("Might rain now - maybe keep a jacket handy.")
		} else {
		fmt.Printf("Might rain in %d hours - maybe keep a jacket handy. \n", peakHour)
		}
	} else {
		fmt.Println("No rain expected.")
	}


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

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.weather.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


