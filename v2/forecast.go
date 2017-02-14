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
	Time                   int64   `json:"time"`
	Summary                string  `json:"summary"`
	Icon                   string  `json:"icon"`
	SunriseTime            int64   `json:"sunriseTime"`
	SunsetTime             int64   `json:"sunsetTime"`
	PrecipIntensity        float64 `json:"precipIntensity"`
	PrecipIntensityMax     float64 `json:"precipIntensityMax"`
	PrecipIntensityMaxTime int64   `json:"precipIntensityMaxTime"`
	PrecipProbability      float64 `json:"precipProbability"`
	PrecipType             string  `json:"precipType"`
	PrecipAccumulation     float64 `json:"precipAccumulation"`
	Temperature            float64 `json:"temperature"`
	TemperatureMin         float64 `json:"temperatureMin"`
	TemperatureMinTime     int64   `json:"temperatureMinTime"`
	TemperatureMax         float64 `json:"temperatureMax"`
	TemperatureMaxTime     int64   `json:"temperatureMaxTime"`
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
	Time        int64   `json:"time"`
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

type Lang string

const (
	Arabic 				Lang = "ar"
	Azerbaijani 			Lang = "az"
	Belarusian 			Lang = "be"
	Bosnian 			Lang = "bs"
	Catalan 			Lang = "ca"
	Czech 				Lang = "cs"
	German 				Lang = "de"
	Greek 				Lang = "el"
	English 			Lang = "en"
	Spanish 			Lang = "es"
	Estonian 			Lang = "et"
	French 				Lang = "fr"
	Croatian 			Lang = "hr"
	Hungarian 			Lang = "hu"
	Indonesian 			Lang = "id"
	Italian 			Lang = "it"
	Icelandic 			Lang = "is"
	Cornish 			Lang = "kw"
	Indonesia 			Lang = "nb"
	Dutch 				Lang = "nl"
	Polish 				Lang = "pl"
	Portuguese 			Lang = "pt"
	Russian 			Lang = "ru"
	Slovak 				Lang = "sk"
	Slovenian 			Lang = "sl"
	Serbian 			Lang = "sr"
	Swedish 			Lang = "sv"
	Tetum 				Lang = "te"
	Turkish 			Lang = "tr"
	Ukrainian 			Lang = "uk"
	IgpayAtinlay			Lang = "x-pig-latin"
	SimplifiedChinese		Lang = "zh"
	TraditionalChinese		Lang = "zh-tw"
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
