forecast
========

Forecast.io v2 API wrapper in Go (Golang)

Documentation: https://developer.forecast.io/docs/v2

Example usage:

```
package main

import (
	"fmt"
	"io/ioutil"
	"log"
    "github.com/mlbright/forecast/v2"
)

func main() {

	keybytes, err := ioutil.ReadFile("api_key.txt")
	if err != nil {
		log.Fatal(err)
	}
	key := string(keybytes)

	lat := "43.6595"
	long := "-79.3433"

	forecast := forecast.Get(key, lat, long, "now")
	fmt.Println(forecast.Timezone)
	fmt.Println(forecast.Currently.Summary)
    fmt.Println(forecast.Currently.Humidity)
    fmt.Println(forecast.Currently.Temperature)
    fmt.Println(forecast.Flags.Units)
}
```