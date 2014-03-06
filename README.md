forecast
========

Forecast.io v2 API wrapper in Go (Golang)

Documentation: https://developer.forecast.io/docs/v2

Example usage:

```
package main

import (
    "fmt"
    forecast "github.com/mlbright/forecast/v2"
    "io/ioutil"
    "log"
    "strings"
)

func main() {

    keybytes, err := ioutil.ReadFile("api_key.txt")
    if err != nil {
        log.Fatal(err)
    }
    key := string(keybytes)
    key = strings.TrimSpace(key)

    lat := "43.6595"
    long := "-79.3433"

    f, err := forecast.Get(key, lat, long, "now", forecast.CA)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%s: %s\n", f.Timezone, f.Currently.Summary)
    fmt.Printf("humidity: %.2f\n", f.Currently.Humidity)
    fmt.Printf("temperature: %.2f Celsius\n", f.Currently.Temperature)
    fmt.Printf("wind speed: %.2f\n", f.Currently.WindSpeed)

}
```
