# gmetrics

A Google Monitoring emulator for (Stackdriver) metrics.


## self-signed cert

Because GCP monitoring clients use TLS to connect to the Cloud Monitoring service, the emulator needs to operate using a self-signed certificate.