package drivers

import (
	"context"
	"os"
	"testing"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/domain"
	"github.com/joho/godotenv"
)

func Test_ConnectToInfluxDB(t *testing.T) {

	//load environment variable from a file for test purposes
	godotenv.Load("./.env.test")

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Successful connection to InfluxDB",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			influxClient, err := ConnectToInfluxDB()
			if (err != nil) != tt.wantErr {
				t.Errorf("ConnectToInfluxDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			health, err := influxClient.Health(context.Background())
			if (err != nil) && health.Status == domain.HealthCheckStatusPass {
				t.Errorf("connectToInfluxDB() error. database not healthy")
				return
			}
			influxClient.Close()
		})
	}
}

func init_testDB(t *testing.T) influxdb2.Client {
	t.Helper()                   // Tells `go test` that this is an helper
	godotenv.Load("./.env.test") //load environement variable
	orgName := os.Getenv("DOCKER_INFLUXDB_INIT_ORG")
	bucketName := os.Getenv("DOCKER_INFLUXDB_INIT_BUCKET")

	client, err := ConnectToInfluxDB() // create the client

	if err != nil {
		t.Errorf("impossible to connect to DB")
	}

	// Clean the database by deleting the bucket
	ctx := context.Background()
	bucketsAPI := client.BucketsAPI()
	dBucket, err := bucketsAPI.FindBucketByName(ctx, bucketName)
	if err == nil {
		client.BucketsAPI().DeleteBucketWithID(context.Background(), *dBucket.Id)
	}

	// create new empty bucket
	dOrg, _ := client.OrganizationsAPI().FindOrganizationByName(ctx, orgName)
	_, err = client.BucketsAPI().CreateBucketWithNameWithID(ctx, *dOrg.Id, bucketName)

	if err != nil {
		t.Errorf("impossible to new create bucket")
	}

	return client
}

func Test_write_event_with_line_protocol(t *testing.T) {

	tests := []struct {
		name  string
		f     func(influxdb2.Client, []ThermostatSetting)
		datas []ThermostatSetting
	}{
		{
			name: "Write new record with line protocol",
			// Your data Points
			datas: []ThermostatSetting{{User: "foo", Avg: 35.5, Max: 42}},
			f: func(c influxdb2.Client, datas []ThermostatSetting) {
				// Send all the data to the DB
				for _, data := range datas {
					Write_event_with_line_protocol(c, data)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// helper to initialise and clean the database
			client := init_testDB(t)
			// call function under test
			tt.f(client, tt.datas)
			// TODO Validate the data
		})
	}
}
