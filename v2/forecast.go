package forecast

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// URL example:  "https://api.forecast.io/forecast/APIKEY/LATITUDE,LONGITUDE,TIME?units=ca"
const (
	BASEURL = "https://api.forecast.io/forecast"
)

type Flags struct {
	DarkSkyUnavailable string
	DarkSkyStations    []string
	DataPointStations  []string
	ISDStations        []string
	LAMPStations       []string
	METARStations      []string
	METNOLicense       string
	Sources            []string
	Units              string
}

type DataPoint struct {
	Time                   float64
	Summary                string
	Icon                   string
	SunriseTime            float64
	SunsetTime             float64
	PrecipIntensity        float64
	PrecipIntensityMax     float64
	PrecipIntensityMaxTime float64
	PrecipProbability      float64
	PrecipType             string
	PrecipAccumulation     float64
	Temperature            float64
	TemperatureMin         float64
	TemperatureMinTime     float64
	TemperatureMax         float64
	TemperatureMaxTime     float64
	DewPoint               float64
	WindSpeed              float64
	WindBearing            float64
	CloudCover             float64
	Humidity               float64
	Pressure               float64
	Visibility             float64
	Ozone                  float64
}

type DataBlock struct {
	Summary string
	Icon    string
	Data    []DataPoint
}

type alert struct {
	Title   string
	Expires float64
	URI     string
}

type Forecast struct {
	Latitude  float64
	Longitude float64
	Timezone  string
	Offset    float64
	Currently DataPoint
	Minutely  DataBlock
	Hourly    DataBlock
	Daily     DataBlock
	Alerts    []alert
	Flags     Flags
	APICalls  int
}

type Units string

const (
	CA Units = "ca"
	SI Units = "si"
)

func Get(key string, lat string, long string, time string, units Units) (*Forecast, error) {
	coord := lat + "," + long

	var url string
	if time == "now" {
		url = BASEURL + "/" + key + "/" + coord + "?units=" + string(units)
	} else {
		url = BASEURL + "/" + key + "/" + coord + "," + time + "?units=" + string(units)
	}

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var f Forecast
	err = json.Unmarshal(body, &f)
	if err != nil {
		return nil, err
	}

	calls, _ := strconv.Atoi(res.Header.Get("X-Forecast-API-Calls"))
	f.APICalls = calls

	return &f, nil
}
