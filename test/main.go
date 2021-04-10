package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"google.golang.org/api/option"
)

func main() {
	localOptions := []option.ClientOption{
		option.WithEndpoint("localhost:9443")}

	exporter, err := stackdriver.NewExporter(
		stackdriver.Options{
			ProjectID:               "google-project-id",
			BundleCountThreshold:    3,
			ReportingInterval:       1 * time.Second,
			MonitoringClientOptions: localOptions,
			TraceClientOptions:      localOptions})
	if err != nil {
		log.Fatal(err)
	}
	defer exporter.Flush()

	// Export to Stackdriver Monitoring.
	if err = exporter.StartMetricsExporter(); err != nil {
		log.Fatal(err)
	}

	// Subscribe views to see stats in Stackdriver Monitoring.
	if err := view.Register(
		ochttp.ClientRoundtripLatencyDistribution,
		ochttp.ClientReceivedBytesDistribution,
	); err != nil {
		log.Fatal(err)
	}

	// Export to Stackdriver Trace.
	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	// Automatically add a Stackdriver trace header to outgoing requests:
	client := &http.Client{
		Transport: &ochttp.Transport{
			Propagation: &propagation.HTTPFormat{},
		},
	}
	_ = client // use client

	// All outgoing requests from client will include a Stackdriver Trace header.
	// See the ochttp package for how to handle incoming requests.

	for i := 0; i < 100; i++ {
		time.Sleep(1 * time.Second)
		fmt.Print(".")
		resp, err := client.Get("https://melrōse.org")
		if err != nil {
			fmt.Println(err)
		} else {
			io.ReadAll(resp.Body)
			resp.Body.Close()
		}
	}
}
