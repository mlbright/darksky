package forecast

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

// URL example:  "https://api.darksky.net/forecast/APIKEY/LATITUDE,LONGITUDE,TIME?units=ca&lang=en"
const (
	BASEURL = "https://api.darksky.net/forecast"
)

type Flags struct {
	DarkSkyUnavailable string   `json:"darksky-unavailable,omitempty"`
	DarkSkyStations    []string `json:"darksky-stations,omitempty"`
	DataPointStations  []string `json:"datapoint-stations,omitempty"`
	ISDStations        []string `json:"isds-stations,omitempty"`
	LAMPStations       []string `json:"lamp-stations,omitempty"`
	MADISStations      []string `json:"madis-stations,omitempty"`
	METARStations      []string `json:"metars-stations,omitempty"`
	METNOLicense       string   `json:"metnol-license,omitempty"`
	Sources            []string `json:"sources,omitempty"`
	Units              string   `json:"units,omitempty"`
}

type DataPoint struct {
	Time                       int64   `json:"time,omitempty"`
	Summary                    string  `json:"summary,omitempty"`
	Icon                       string  `json:"icon,omitempty"`
	SunriseTime                int64   `json:"sunriseTime,omitempty"`
	SunsetTime                 int64   `json:"sunsetTime,omitempty"`
	PrecipIntensity            float64 `json:"precipIntensity,omitempty"`
	PrecipIntensityMax         float64 `json:"precipIntensityMax,omitempty"`
	PrecipIntensityMaxTime     int64   `json:"precipIntensityMaxTime,omitempty"`
	PrecipProbability          float64 `json:"precipProbability,omitempty"`
	PrecipType                 string  `json:"precipType,omitempty"`
	PrecipAccumulation         float64 `json:"precipAccumulation,omitempty"`
	Temperature                float64 `json:"temperature,omitempty"`
	TemperatureMin             float64 `json:"temperatureMin,omitempty"`
	TemperatureMinTime         int64   `json:"temperatureMinTime,omitempty"`
	TemperatureMax             float64 `json:"temperatureMax,omitempty"`
	TemperatureMaxTime         int64   `json:"temperatureMaxTime,omitempty"`
	ApparentTemperature        float64 `json:"apparentTemperature,omitempty"`
	ApparentTemperatureMin     float64 `json:"apparentTemperatureMin,omitempty"`
	ApparentTemperatureMinTime int64   `json:"apparentTemperatureMinTime,omitempty"`
	ApparentTemperatureMax     float64 `json:"apparentTemperatureMax,omitempty"`
	ApparentTemperatureMaxTime int64   `json:"apparentTemperatureMaxTime,omitempty"`
	NearestStormBearing        float64 `json:"nearestStormBearing,omitempty"`
	NearestStormDistance       float64 `json:"nearestStormDistance,omitempty"`
	DewPoint                   float64 `json:"dewPoint,omitempty"`
	WindSpeed                  float64 `json:"windSpeed,omitempty"`
	WindBearing                float64 `json:"windBearing,omitempty"`
	CloudCover                 float64 `json:"cloudCover,omitempty"`
	Humidity                   float64 `json:"humidity,omitempty"`
	Pressure                   float64 `json:"pressure,omitempty"`
	Visibility                 float64 `json:"visibility,omitempty"`
	Ozone                      float64 `json:"ozone,omitempty"`
	MoonPhase                  float64 `json:"moonPhase,omitempty"`
}

type DataBlock struct {
	Summary string      `json:"summary,omitempty"`
	Icon    string      `json:"icon,omitempty"`
	Data    []DataPoint `json:"data,omitempty"`
}

type alert struct {
	Title       string   `json:"title,omitempty"`
	Regions     []string `json:"regions,omitempty"`
	Severity    string   `json:"severity,omitempty"`
	Description string   `json:"description,omitempty"`
	Time        int64    `json:"time,omitempty"`
	Expires     float64  `json:"expires,omitempty"`
	URI         string   `json:"uri,omitempty"`
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

type Lang string

const (
	Arabic             Lang = "ar"
	Azerbaijani        Lang = "az"
	Belarusian         Lang = "be"
	Bosnian            Lang = "bs"
	Catalan            Lang = "ca"
	Czech              Lang = "cs"
	German             Lang = "de"
	Greek              Lang = "el"
	English            Lang = "en"
	Spanish            Lang = "es"
	Estonian           Lang = "et"
	French             Lang = "fr"
	Croatian           Lang = "hr"
	Hungarian          Lang = "hu"
	Indonesian         Lang = "id"
	Italian            Lang = "it"
	Icelandic          Lang = "is"
	Cornish            Lang = "kw"
	Indonesia          Lang = "nb"
	Dutch              Lang = "nl"
	Polish             Lang = "pl"
	Portuguese         Lang = "pt"
	Russian            Lang = "ru"
	Slovak             Lang = "sk"
	Slovenian          Lang = "sl"
	Serbian            Lang = "sr"
	Swedish            Lang = "sv"
	Tetum              Lang = "te"
	Turkish            Lang = "tr"
	Ukrainian          Lang = "uk"
	IgpayAtinlay       Lang = "x-pig-latin"
	SimplifiedChinese  Lang = "zh"
	TraditionalChinese Lang = "zh-tw"
)

func Get(key string, lat string, long string, time string, units Units, lang Lang) (*Forecast, error) {
	res, err := GetResponse(key, lat, long, time, units, lang)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	f, err := FromJSON(res.Body)
	if err != nil {
		return nil, err
	}

	calls, _ := strconv.Atoi(res.Header.Get("X-Forecast-API-Calls"))
	f.APICalls = calls

	return f, nil
}

func FromJSON(reader io.Reader) (*Forecast, error) {
	var f Forecast
	if err := json.NewDecoder(reader).Decode(&f); err != nil {
		return nil, err
	}

	return &f, nil
}

func GetResponse(key string, lat string, long string, time string, units Units, lang Lang) (*http.Response, error) {
	coord := lat + "," + long

	var url string
	if time == "now" {
		url = BASEURL + "/" + key + "/" + coord + "?units=" + string(units) + "&lang=" + string(lang)
	} else {
		url = BASEURL + "/" + key + "/" + coord + "," + time + "?units=" + string(units) + "&lang=" + string(lang)
	}

	res, err := http.Get(url)
	if err != nil {
		return res, err
	}

	return res, nil
}
