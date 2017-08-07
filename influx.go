package main

import (
	"errors"
	"log"
	"time"

	influx "github.com/influxdata/influxdb/client/v2"
)

func NewMetric(measurement string, tags map[string]string, fields map[string]interface{}) *Metric {
	return &Metric{
		measurement: measurement,
		tags:        tags,
		fields:      fields,
		time:        time.Now(),
	}
}

func (influxDB *InfluxClient) Send(i *IPFIX, s string) {
	tags := map[string]string{
		"Host":       *name,
		"MetricType": s,
	}
	fields := map[string]interface{}{
		"inc.rtp.mos":          float32(i.Data.QOS.IncMos) / 100,
		"out.rtp.mos":          float32(i.Data.QOS.OutMos) / 100,
		"inc.rtp.rval":         float32(i.Data.QOS.IncrVal) / 100,
		"out.rtp.rval":         float32(i.Data.QOS.OutrVal) / 100,
		"inc.rtp.packets":      float32(i.Data.QOS.IncRtpPackets),
		"out.rtp.packets":      float32(i.Data.QOS.OutRtpPackets),
		"inc.rtcp.packets":     float32(i.Data.QOS.IncRtcpPackets),
		"out.rtcp.packets":     float32(i.Data.QOS.OutRtcpPackets),
		"inc.rtp.lostPackets":  float32(i.Data.QOS.IncRtpLostPackets),
		"out.rtp.lostPackets":  float32(i.Data.QOS.OutRtpLostPackets),
		"inc.rtcp.lostPackets": float32(i.Data.QOS.IncRtcpLostPackets),
		"out.rtcp.lostPackets": float32(i.Data.QOS.OutRtcpLostPackets),
		"inc.rtp.avgJitter":    float32(i.Data.QOS.IncRtpAvgJitter),
		"out.rtp.avgJitter":    float32(i.Data.QOS.OutRtpAvgJitter),
		"inc.rtp.maxJitter":    float32(i.Data.QOS.IncRtpMaxJitter),
		"out.rtp.maxJitter":    float32(i.Data.QOS.OutRtpMaxJitter),
		"inc.rtcp.avgJitter":   float32(i.Data.QOS.IncRtcpAvgJitter),
		"out.rtcp.avgJitter":   float32(i.Data.QOS.OutRtcpAvgJitter),
		"inc.rtcp.maxJitter":   float32(i.Data.QOS.IncRtcpMaxJitter),
		"out.rtcp.maxJitter":   float32(i.Data.QOS.OutRtcpMaxJitter),
		"inc.rtcp.avgLat":      float32(i.Data.QOS.IncRtcpAvgLat),
		"out.rtcp.avgLat":      float32(i.Data.QOS.OutRtcpAvgLat),
		"inc.rtcp.maxLat":      float32(i.Data.QOS.IncRtcpMaxLat),
		"out.rtcp.maxLat":      float32(i.Data.QOS.OutRtcpMaxLat),
	}

	if err := influxDB.send(NewMetric("horaclifix", tags, fields)); err != nil {
		log.Printf("Could not send metric to influxDB: %s\n", err.Error())
		if influxDB.errorFunc != nil {
			influxDB.errorFunc(err)
		}
	}
}

func (influxDB *InfluxClient) send(metric *Metric) error {
	if influxDB == nil {
		return errors.New("Failed to create influxDB client")
	}

	pt, err := influx.NewPoint(metric.measurement, metric.tags, metric.fields, metric.time)
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
