package api

import (
	json2 "encoding/json"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	aqi    = "yes"
	alerts = "no"
)

type Api interface {
	InitApi()
	ParseArguments()
	MakeRequest()
	FormatOutput()
}

type WeatherApi struct {
	key    string // Weather API service key
	days   string // Days of forecast, which API would return
	q      string // The City name
	aqi    string // Air Quality
	alerts string
}

func (w *WeatherApi) InitApi() {
	/*
		Init interface for WeatherApi
	*/
	w.aqi = aqi
	w.alerts = alerts

	if key, exists := os.LookupEnv("API_KEY"); exists {
		w.key = key
		return
	}

	w.key = "1111111111111111111111111111111"
}

func (w *WeatherApi) ParseArguments(args []string) {
	/*
		Setter for WeatherApi, init from system args
	*/
	var (
		days bool
		city bool
	)

	for pos, arg := range args {
		if arg[0] != '-' {
			continue
		}
		if strings.Compare(arg, "-d") == 0 {
			w.days = args[pos+1]
			days = true
			continue
		}
		if strings.Compare(arg, "-c") == 0 {
			w.q = args[pos+1]
			city = true
			continue
		}
		if strings.Compare(arg, "-help") == 0 {
			fmt.Println("CLI Forecast, based on Go")
			fmt.Println("Use 'forecast -d <days> -c <city>")
			fmt.Println("	or -help for this message")
			fmt.Println("If you have an API_KEY of https://www.weatherapi.com/ use 'export API_KEY=<key>")
			os.Exit(1)
		}

		log.Fatal("Unexpected argument, try with -help")
	}

	if !days {
		w.days = "1"
	}
	if !city {
		log.Fatal("The city was not specified, try with -help")
	}
}

func (w *WeatherApi) MakeRequest() map[string]interface{} {
	/*
		Make Request on Weather API, parse and return JSON format
	*/
	var url string = fmt.Sprintf("http://api.weatherapi.com/v1/forecast.json?key=%s&q=%s&days=%s&aqi=%s&alerts=%s", w.key, w.q, w.days, w.aqi, w.alerts)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Error while init Http Request")
	}
	req.Header.Add("accept", "application/json")

	res, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		log.Fatal("Error while making GET request")
	}

	body, _ := io.ReadAll(res.Body)
	jsonStr := string(body)

	var message map[string]interface{}

	_ = json2.Unmarshal([]byte(jsonStr), &message)
	if err != nil {
		log.Fatal("Error while making json from Request Body")
	}

	return message
}

func (w *WeatherApi) FormatOutput(json map[string]interface{}) {
	/*
		Formatted output in CLI table
	*/
	location := json["location"].(map[string]interface{})
	caption := fmt.Sprintf("Forecast for %s (%s)\n", location["name"], location["country"])

	t := table.NewWriter()
	t.SetCaption(caption)
	t.AppendHeader(table.Row{"Date", "Weather", "MaxT", "MinT", "AvgT", "ChanceOfRain", "ChanceOfSnow", "MaxWindSpeed", "AirQuality", "Sunrise", "Sunset", "Moonrise", "Moonset"})

	forecast := json["forecast"].(map[string]interface{})
	days := forecast["forecastday"].([]interface{})

	for i := 0; i < len(days); i++ {
		dayInfo := days[i].(map[string]interface{})
		day := dayInfo["day"].(map[string]interface{})
		condition := day["condition"].(map[string]interface{})
		air_quality := day["air_quality"].(map[string]interface{})
		astro := dayInfo["astro"].(map[string]interface{})

		date := dayInfo["date"]
		maxt := day["maxtemp_c"]
		mint := day["mintemp_c"]
		avgt := day["avgtemp_c"]
		chanceofrain := day["daily_chance_of_rain"]
		chanceofsnow := day["daily_chance_of_snow"]
		maxwinds := day["maxwind_kph"]
		weather := condition["text"]
		co := air_quality["co"]
		sunrise := astro["sunrise"]
		sunset := astro["sunset"]
		moonrise := astro["moonrise"]
		moonset := astro["moonset"]

		t.AppendRow(table.Row{date, weather, maxt, mint, avgt, chanceofrain, chanceofsnow, maxwinds, co, sunrise, sunset, moonrise, moonset})
	}

	fmt.Println(t.Render())

}
