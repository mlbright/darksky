package forecast

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

// URL example:  "https://api.darksky.net/forecast/APIKEY/LATITUDE,LONGITUDE,TIME?units=ca"
const (
	BASEURL = "https://api.darksky.net/forecast"
)

type Flags struct {
	DarkSkyUnavailable string   `json:"darksky-unavailable"`
	DarkSkyStations    []string `json:"darksky-stations"`
	DataPointStations  []string `json:"datapoint-stations"`
	ISDStations        []string `json:"isds-stations"`
	LAMPStations       []string `json:"lamp-stations"`
	METARStations      []string `json:"metars-stations"`
	METNOLicense       string   `json:"metnol-license"`
	Sources            []string `json:"sources"`
	Units              string   `json:"units"`
}

type DataPoint struct {
	Time                   float64 `json:"time"`
	Summary                string  `json:"summary"`
	Icon                   string  `json:"icon"`
	SunriseTime            float64 `json:"sunriseTime"`
	SunsetTime             float64 `json:"sunsetTime"`
	PrecipIntensity        float64 `json:"precipIntensity"`
	PrecipIntensityMax     float64 `json:"precipIntensityMax"`
	PrecipIntensityMaxTime float64 `json:"precipIntensityMaxTime"`
	PrecipProbability      float64 `json:"precipProbability"`
	PrecipType             string  `json:"precipType"`
	PrecipAccumulation     float64 `json:"precipAccumulation"`
	Temperature            float64 `json:"temperature"`
	TemperatureMin         float64 `json:"temperatureMin"`
	TemperatureMinTime     float64 `json:"temperatureMinTime"`
	TemperatureMax         float64 `json:"temperatureMax"`
	TemperatureMaxTime     float64 `json:"temperatureMaxTime"`
	ApparentTemperature    float64 `json:"apparentTemperature"`
	DewPoint               float64 `json:"dewPoint"`
	WindSpeed              float64 `json:"windSpeed"`
	WindBearing            float64 `json:"windBearing"`
	CloudCover             float64 `json:"cloudCover"`
	Humidity               float64 `json:"humidity"`
	Pressure               float64 `json:"pressure"`
	Visibility             float64 `json:"visibility"`
	Ozone                  float64 `json:"ozone"`
	MoonPhase              float64 `json:"moonPhase"`
}

type DataBlock struct {
	Summary string      `json:"summary"`
	Icon    string      `json:"icon"`
	Data    []DataPoint `json:"data"`
}

type alert struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Time        float64 `json:"time"`
	Expires     float64 `json:"expires"`
	URI         string  `json:"uri"`
}

type Forecast struct {
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Timezone  string    `json:"timezone"`
	Offset    float64   `json:"offset"`
	Currently DataPoint `json:"currently"`
	Minutely  DataBlock `json:"minutely"`
	Hourly    DataBlock `json:"hourly"`
	Daily     DataBlock `json:"daily"`
	Alerts    []alert   `json:"alerts"`
	Flags     Flags     `json:"flags"`
	APICalls  int       `json:"apicalls"`
	Code      int       `json:"code"`
}

type Units string

const (
	CA   Units = "ca"
	SI   Units = "si"
	US   Units = "us"
	UK   Units = "uk"
	AUTO Units = "auto"
)

func Get(key string, lat string, long string, time string, units Units) (*Forecast, error) {
	res, err := GetResponse(key, lat, long, time, units)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	f, err := FromJSON(body)
	if err != nil {
		return nil, err
	}

	calls, _ := strconv.Atoi(res.Header.Get("X-Forecast-API-Calls"))
	f.APICalls = calls

	return f, nil
}

func FromJSON(json_blob []byte) (*Forecast, error) {
	var f Forecast
	err := json.Unmarshal(json_blob, &f)
	if err != nil {
		return nil, err
	}

	return &f, nil
}

func GetResponse(key string, lat string, long string, time string, units Units) (*http.Response, error) {
	coord := lat + "," + long

	var url string
	if time == "now" {
		url = BASEURL + "/" + key + "/" + coord + "?units=" + string(units)
	} else {
		url = BASEURL + "/" + key + "/" + coord + "," + time + "?units=" + string(units)
	}

	res, err := http.Get(url)
	if err != nil {
		return res, err
	}

	return res, nil
}
