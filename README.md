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
)

func main() {

	keybytes, err := ioutil.ReadFile("api_key.txt")
	if err != nil {
		log.Fatal(err)
	}
	key := string(keybytes)

	lat := "43.6595"
	long := "-79.3433"

	f := forecast.Get(key, lat, long, "now")
	fmt.Println(f.Timezone)
	fmt.Println(f.Currently.Summary)
	fmt.Println(f.Currently.Humidity)
	fmt.Println(f.Currently.Temperature)
	fmt.Println(f.Flags.Units)
	fmt.Println(f.Currently.WindSpeed)
}
```