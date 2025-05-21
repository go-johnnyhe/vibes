package api

import (
	"encoding/json"
	"fmt"
	"os"
	"io/ioutil"
	"net/http"
)

func GetWeatherData() {
	apiKey := "fd1a9fe1168620101f0bcef74da06706"
	lat := 47.6038321
	lon := -122.330062
	basicURL := fmt.Sprintf("https://api.openweathermap.org/data/3.0/onecall?lat=%f&lon=%f&appid=%s", lat, lon, apiKey)
	
	response, err := http.Get(basicURL)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading the response:", err)
		os.Exit(1)
	}
	// fmt.Println(string(body))
	var prettyJSON map[string]interface{}
	json.Unmarshal(body, &prettyJSON)
	formatted, _ := json.MarshalIndent(prettyJSON, "", "  ")
	fmt.Println(string(formatted))
}
