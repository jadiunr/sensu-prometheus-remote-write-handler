package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/m3db/prometheus_remote_client_golang/promremote"
	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/sensu/sensu-plugin-sdk/sensu"
)

// Config represents the handler plugin config.
type Config struct {
	sensu.PluginConfig
	Endpoint string
	Timeout string
	Headers []string
	IncludeEventStatus bool
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-prometheus-remote-write-handler",
			Short:    "Prometheus remote write Handler for Sensu",
			Keyspace: "sensu.io/plugins/sensu-prometheus-remote-write-handler/config",
		},
	}

	options = []sensu.ConfigOption{
		&sensu.PluginConfigOption[string]{
			Path:      "endpoint",
			Argument:  "endpoint",
			Shorthand: "e",
			Default:   "",
			Usage:     "Remote write endpoint",
			Value:     &plugin.Endpoint,
		},
		&sensu.PluginConfigOption[string]{
			Path: "timeout",
			Argument: "timeout",
			Shorthand: "t",
			Default: "10s",
			Usage: "Remote write timeout",
			Value: &plugin.Timeout,
		},
		&sensu.SlicePluginConfigOption[string]{
			Path: "header",
			Argument: "header",
			Shorthand: "H",
			Default: []string{},
			Usage: "Additional header(s) to send in remote write request",
			Value: &plugin.Headers,
		},
		&sensu.PluginConfigOption[bool]{
			Path: "include-event-status",
			Argument: "include-event-status",
			Shorthand: "i",
			Default: false,
			Usage: "If true, the check status result will be captured as a metric",
			Value: &plugin.IncludeEventStatus,
		},
	}
)

func main() {
	handler := sensu.NewGoHandler(&plugin.PluginConfig, options, checkArgs, executeHandler)
	handler.Execute()
}

func checkArgs(event *corev2.Event) error {
	return nil
}

func executeHandler(event *corev2.Event) error {
	timeout, err := time.ParseDuration(plugin.Timeout)
	if err != nil {
		return err
	}

	cfg := promremote.NewConfig(
		promremote.WriteURLOption(plugin.Endpoint),
		promremote.HTTPClientTimeoutOption(timeout),
	)

	client, err := promremote.NewClient(cfg)
	if err != nil {
		return err
	}

	var timeSeriesList promremote.TSList

	for _, point := range event.Metrics.Points {
		var labels []promremote.Label
		labels = append(labels, promremote.Label{
			Name: "__name__",
			Value: strings.Split(point.Name, ".")[0],
		})
		labels = append(labels, promremote.Label{
			Name: "sensu_entity_name",
			Value: event.Entity.Name,
		})
		for _, tag := range point.Tags {
			labels = append(labels, promremote.Label{
				Name: tag.Name,
				Value: tag.Value,
			})
		}
		timestamp, err := convertInt64ToTime(point.Timestamp)
		if err != nil {
			return err
		}
		timeSeriesList = append(timeSeriesList, promremote.TimeSeries{
			Labels: labels,
			Datapoint: promremote.Datapoint{
				Timestamp: timestamp,
				Value: point.Value,
			},
		})
	}

	if plugin.IncludeEventStatus && event.HasCheck() {
		var labels []promremote.Label
		labels = append(labels, promremote.Label{
			Name: "entity",
			Value: event.Entity.Name,
		})
		labels = append(labels, promremote.Label{
			Name: "check",
			Value: event.Check.Name,
		})
		timestamp, err := convertInt64ToTime(event.Timestamp)
		if err != nil {
			return err
		}
		timeSeriesList = append(timeSeriesList, promremote.TimeSeries{
			Labels: append(labels, promremote.Label{
				Name: "__name__",
				Value: "sensu_event_status",
			}),
			Datapoint: promremote.Datapoint{
				Timestamp: timestamp,
				Value: float64(event.Check.Status),
			},
		})
		timeSeriesList = append(timeSeriesList, promremote.TimeSeries{
			Labels: append(labels, promremote.Label{
				Name: "__name__",
				Value: "sensu_event_occurrences",
			}),
			Datapoint: promremote.Datapoint{
				Timestamp: timestamp,
				Value: float64(event.Check.Occurrences),
			},
		})
		timeSeriesList = append(timeSeriesList, promremote.TimeSeries{
			Labels: append(labels, promremote.Label{
				Name: "__name__",
				Value: "sensu_event_silenced",
			}),
			Datapoint: promremote.Datapoint{
				Timestamp: timestamp,
				Value: func() float64 { if event.IsSilenced() { return 1 } else { return 0 } }(),
			},
		})
	}

	writeOptions := promremote.WriteOptions{
		Headers: map[string]string{},
	}
	if len(plugin.Headers) > 0 {
		for _, header := range plugin.Headers {
			headerSplit := strings.SplitN(header, ":", 2)
			writeOptions.Headers[headerSplit[0]] = headerSplit[1]
		}
	}

	ctx := context.Background()
	res, err := client.WriteTimeSeries(ctx, timeSeriesList, writeOptions)
	if err != nil {
		return fmt.Errorf("%v, %v", res, err)
	}

	return nil
}

func convertInt64ToTime(timestamp int64) (time.Time, error) {
	stringTimestamp := strconv.FormatInt(timestamp, 10)
	if len(stringTimestamp) > 10 {
		stringTimestamp = stringTimestamp[:10]
	}
	t, err := strconv.ParseInt(stringTimestamp, 10, 64)
	if err != nil {
		return time.Now(), err
	}

	return time.Unix(t, 0), nil
}
