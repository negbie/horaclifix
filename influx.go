package main

import (
	"errors"
	"log"
	"time"

	influx "github.com/influxdata/influxdb/client/v2"
)

func NewMetric(measurement string, tags map[string]string, fields map[string]interface{}) *InfluxMetric {
	return &InfluxMetric{
		measurement: measurement,
		tags:        tags,
		fields:      fields,
		time:        time.Now(),
	}
}

func (influxDB *InfluxClient) Send(i *IPFIX, s string) {
	tags := map[string]string{
		"host":          *name,
		"metricType":    s,
		"incomingRealm": string(i.QOS.IncRealm),
		"outgoingRealm": string(i.QOS.OutRealm),
	}

	fields := i.mapMetricQOS()

	if err := influxDB.send(NewMetric("horaclifix", tags, fields)); err != nil {
		log.Printf("[WARN] Could not send metric to influxDB: %s\n", err.Error())
		if influxDB.errorFunc != nil {
			influxDB.errorFunc(err)
		}
	}
}

func (influxDB *InfluxClient) send(im *InfluxMetric) error {
	if influxDB == nil {
		return errors.New("Failed to create influxDB client")
	}

	pt, err := influx.NewPoint(im.measurement, im.tags, im.fields, im.time)
	if err != nil {
		return err
	}

	influxDB.pointsChannel <- pt
	return nil
}

func (influxDB *InfluxClient) bulk(points chan *influx.Point, ticker *time.Ticker) {
	pointsBuffer := make([]*influx.Point, influxDB.batchSize)
	currentBatchSize := 0
	for {
		select {
		case <-ticker.C:
			err := influxDB.flush(pointsBuffer, currentBatchSize)
			if err != nil && influxDB.errorFunc != nil {
				influxDB.errorFunc(err)
			}
			currentBatchSize = 0
		case point := <-points:
			pointsBuffer[currentBatchSize] = point
			currentBatchSize++
			if influxDB.batchSize == currentBatchSize {
				err := influxDB.flush(pointsBuffer, currentBatchSize)
				if err != nil && influxDB.errorFunc != nil {
					influxDB.errorFunc(err)
				}
				currentBatchSize = 0
			}
		}
	}
}

func (influxDB *InfluxClient) flush(points []*influx.Point, size int) error {
	if size > 0 {
		newBatch, err := influx.NewBatchPoints(influxDB.batchConfig)
		if err != nil {
			return err
		}
		newBatch.AddPoints(points[0:size])
		err = influxDB.client.Write(newBatch)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewInfluxClient(config *InfluxClientConfig) (*InfluxClient, error) {
	// This is just a connection to see if we have no errors
	httpConfig := influx.HTTPConfig{Addr: config.Endpoint, Timeout: 5 * time.Second}
	influxDBClient, err := influx.NewHTTPClient(httpConfig)
	if err != nil {
		return nil, err
	}
	// Create database if it's not already created
	_, err = influxDBClient.Query(influx.Query{
		Command: "CREATE DATABASE horaclifix",
	})
	if err != nil {
		return nil, err
	}

	iClient := &InfluxClient{
		client:        influxDBClient,
		database:      config.Database,
		batchSize:     config.BatchSize,
		pointsChannel: make(chan *influx.Point, config.BatchSize*50),
		batchConfig:   influx.BatchPointsConfig{Database: config.Database},
		errorFunc:     config.ErrorFunc,
	}
	go iClient.bulk(iClient.pointsChannel, time.NewTicker(config.FlushTimeout))
	return iClient, nil
}
