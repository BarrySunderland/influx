package drivers

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

var org string = os.Getenv("INFLUXDB_ORG")
var bucket string = os.Getenv("INFLUXDB_BUCKET")

// type influxConfig struct {
// 	org    string
// 	bucket string
// }

// Connect to an Influx Database reading the credentials from
// environment variables INFLUXDB_TOKEN, INFLUXDB_URL
// return influxdb Client or errors
func ConnectToInfluxDB() (influxdb2.Client, error) {

	dbToken := os.Getenv("INFLUXDB_TOKEN")
	if dbToken == "" {
		return nil, errors.New("INFLUXDB_TOKEN must be set")
	}

	dbURL := os.Getenv("INFLUXDB_URL")
	if dbURL == "" {
		return nil, errors.New("INFLUXDB_URL must be set")
	}

	client := influxdb2.NewClient(dbURL, dbToken)

	// validate client connection health
	_, err := client.Health(context.Background())

	return client, err
}

func CreateBucketIfNotExists(ctx context.Context, client influxdb2.Client) error {

	dBucket, err := client.BucketsAPI().FindBucketByName(ctx, bucket)
	if dBucket != nil {
		fmt.Println("found bucket: ", bucket)
		return nil
	}

	if dBucket == nil {
		// create new empty bucket
		dOrg, _ := client.OrganizationsAPI().FindOrganizationByName(ctx, org)
		_, err = client.BucketsAPI().CreateBucketWithName(ctx, dOrg, bucket)

		if err != nil {
			return errors.New("impossible to new create bucket")
		}

		fmt.Println("created new bucket: ", bucket)
		return nil
	}
	return nil

}

type ThermostatSetting struct {
	User string
	Max  float64 //temperature
	Avg  float64 //temperature
}

func Write_event_with_line_protocol(client influxdb2.Client, t ThermostatSetting) {
	// get non-blocking write client
	writeAPI := client.WriteAPI(org, bucket)
	// write Line Protocol
	record := fmt.Sprintf("thermostat,unit=temperature,user=%s avg=%f,max=%f", t.User, t.Avg, t.Max)
	writeAPI.WriteRecord(record)
	// Flush Writes
	writeAPI.Flush()
}

//The point data approach is lengthy to write,
//but also provides more structure.
//It’s convenient when data parameters are already
//in the desired format.
func Write_event_with_params_constror(client influxdb2.Client, t ThermostatSetting) {
	// Use blocking write client for writes to desired bucket
	writeAPI := client.WriteAPI(org, bucket)
	// Create point using full params constructor
	p := influxdb2.NewPoint("thermostat",
		map[string]string{"unit": "temperature", "user": t.User},
		map[string]interface{}{"avg": t.Avg, "max": t.Max},
		time.Now())
	writeAPI.WritePoint(p)
	// Flush writes
	writeAPI.Flush()
}

//Alternatively, you can use the builder NewPointWithMeasurement
// to construct the query step by step, which is easy to read.
func Write_event_with_fluent_Style(client influxdb2.Client, t ThermostatSetting) {
	// use blocking write client for writes to desired bucket
	writeAPI := client.WriteAPI(org, bucket)
	//create a point using fluent style
	p := influxdb2.NewPointWithMeasurement("thermostat").
		AddTag("unit", "temperature").
		AddTag("user", t.User).
		AddTag("avg", strconv.FormatFloat(t.Avg, 'f', 6, 64)).
		AddTag("max", strconv.FormatFloat(t.Max, 'f', 6, 64)).
		SetTime(time.Now())
	writeAPI.WritePoint(p)
	// Flush writes
	writeAPI.Flush()
}

// TO DO write tests for these functions
