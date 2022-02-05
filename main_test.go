package main

import (
	"context"
	"testing"

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
