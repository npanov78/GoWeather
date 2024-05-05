package main

import (
	"os"
	api2 "weather/api"
)

func startApp(args []string) {
	api := api2.WeatherApi{}
	api.InitApi()
	api.ParseArguments(args)
	json := api.MakeRequest()
	api.FormatOutput(json)
}

func main() {
	/*
		CLI util for forecast news, based on Weather API service
	*/
	args := os.Args[1:]
	startApp(args)
}
