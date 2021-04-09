# gcpmon

A Google Cloud Platform Monitoring emulator for (Stackdriver) metrics.


## Go

If you are using `contrib.go.opencensus.io/exporter/stackdriver` then change the endpoint by setting options:

    MonitoringClientOptions: []option.ClientOption{
	    option.WithEndpoint("localhost:9443")}},
## self-signed certificate

Because GCP monitoring clients use TLS to connect to the Cloud Monitoring service, the emulator needs to operate using a self-signed certificate.