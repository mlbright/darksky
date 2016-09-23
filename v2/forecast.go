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
	DarkSkyUnavailable string   `json:"darksky-unavailable,omitempty"`
	DarkSkyStations    []string `json:"darksky-stations,omitempty"`
	DataPointStations  []string `json:"datapoint-stations,omitempty"`
	ISDStations        []string `json:"isds-stations,omitempty"`
	LAMPStations       []string `json:"lamp-stations,omitempty"`
	METARStations      []string `json:"metars-stations,omitempty"`
	METNOLicense       string   `json:"metnol-license,omitempty"`
	Sources            []string `json:"sources,omitempty"`
	Units              string   `json:"units,omitempty"`
}

type DataPoint struct {
	Time                   float64 `json:"time,omitempty"`
	Summary                string  `json:"summary,omitempty"`
	Icon                   string  `json:"icon,omitempty"`
	SunriseTime            float64 `json:"sunriseTime,omitempty"`
	SunsetTime             float64 `json:"sunsetTime,omitempty"`
	PrecipIntensity        float64 `json:"precipIntensity,omitempty"`
	PrecipIntensityMax     float64 `json:"precipIntensityMax,omitempty"`
	PrecipIntensityMaxTime float64 `json:"precipIntensityMaxTime,omitempty"`
	PrecipProbability      float64 `json:"precipProbability,omitempty"`
	PrecipType             string  `json:"precipType,omitempty"`
	PrecipAccumulation     float64 `json:"precipAccumulation,omitempty"`
	Temperature            float64 `json:"temperature,omitempty"`
	TemperatureMin         float64 `json:"temperatureMin,omitempty"`
	TemperatureMinTime     float64 `json:"temperatureMinTime,omitempty"`
	TemperatureMax         float64 `json:"temperatureMax,omitempty"`
	TemperatureMaxTime     float64 `json:"temperatureMaxTime,omitempty"`
	ApparentTemperature    float64 `json:"apparentTemperature,omitempty"`
	DewPoint               float64 `json:"dewPoint,omitempty"`
	WindSpeed              float64 `json:"windSpeed,omitempty"`
	WindBearing            float64 `json:"windBearing,omitempty"`
	CloudCover             float64 `json:"cloudCover,omitempty"`
	Humidity               float64 `json:"humidity,omitempty"`
	Pressure               float64 `json:"pressure,omitempty"`
	Visibility             float64 `json:"visibility,omitempty"`
	Ozone                  float64 `json:"ozone,omitempty"`
	MoonPhase              float64 `json:"moonPhase,omitempty"`
}

type DataBlock struct {
	Summary string      `json:"summary,omitempty"`
	Icon    string      `json:"icon,omitempty"`
	Data    []DataPoint `json:"data,omitempty"`
}

type alert struct {
	Title       string  `json:"title,omitempty"`
	Description string  `json:"description,omitempty"`
	Time        float64 `json:"time,omitempty"`
	Expires     float64 `json:"expires,omitempty"`
	URI         string  `json:"uri,omitempty"`
}

type Forecast struct {
	Latitude  float64   `json:"latitude,omitempty"`
	Longitude float64   `json:"longitude,omitempty"`
	Timezone  string    `json:"timezone,omitempty"`
	Offset    float64   `json:"offset,omitempty"`
	Currently DataPoint `json:"currently,omitempty"`
	Minutely  DataBlock `json:"minutely,omitempty"`
	Hourly    DataBlock `json:"hourly,omitempty"`
	Daily     DataBlock `json:"daily,omitempty"`
	Alerts    []alert   `json:"alerts,omitempty"`
	Flags     Flags     `json:"flags,omitempty"`
	APICalls  int       `json:"apicalls,omitempty"`
	Code      int       `json:"code,omitempty"`
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
