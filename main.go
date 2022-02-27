package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/barrysunderland/influx/drivers"
)

func main() {

	client, err := drivers.ConnectToInfluxDB()
	if err != nil {
		fmt.Println("err")
		os.Exit(1)
	}

	drivers.CreateBucketIfNotExists(context.Background(), client)

	var newT drivers.ThermostatSetting
	numSteps := 100.0
	for i := 0.0; i < numSteps; i++ {
		max := 30.4 * math.Sin(i/10.0)
		newT = drivers.ThermostatSetting{
			User: "tempCorp",
			Max:  max,
			Avg:  28.1,
		}
		drivers.Write_event_with_params_constror(client, newT)

		fmt.Printf("\rwriting record %v of %v", i+1, numSteps)
		time.Sleep(time.Second / 10)

	}
	fmt.Println(": done. view in UI at http://localhost:8086")

}
