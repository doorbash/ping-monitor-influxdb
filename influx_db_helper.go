package main

import (
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type InfluxDBHelper struct {
	Client    influxdb2.Client
	AppConfig *AppConfig
}

func (i *InfluxDBHelper) Close() {
	i.Client.Close()
}

func (i *InfluxDBHelper) SendData(data PingResult) {
	p := influxdb2.NewPointWithMeasurement("Ping").
		AddTag("host", data.Address).
		AddField("response_time", data.Time).
		AddField("ttl", data.TTL).
		AddField("error", data.Error).
		SetTime(time.Now())

	i.Client.WriteAPI(i.AppConfig.Influx.Org, i.AppConfig.Influx.Bucket).WritePoint(p)
}

func NewInfluxDBHelper(appConfig *AppConfig) *InfluxDBHelper {
	client := influxdb2.NewClient(appConfig.Influx.URL, appConfig.Influx.Secret)

	client.Options().SetFlushInterval(10000)

	return &InfluxDBHelper{
		Client:    client,
		AppConfig: appConfig,
	}
}
