package main

import (
	"context"
	"log"

	"google.golang.org/genproto/googleapis/api/metric"
	"google.golang.org/genproto/googleapis/api/monitoredres"
	. "google.golang.org/genproto/googleapis/monitoring/v3"
	"google.golang.org/protobuf/types/known/emptypb"
)

// implements monitoringpb.MetricServiceServer
type MetricsService struct {
}

// Lists monitored resource descriptors that match a filter. This method does not require a Workspace.
func (s *MetricsService) ListMonitoredResourceDescriptors(ctx context.Context, req *ListMonitoredResourceDescriptorsRequest) (*ListMonitoredResourceDescriptorsResponse, error) {
	log.Printf("ListMonitoredResourceDescriptors:%#v", req)
	return nil, nil
}

// Gets a single monitored resource descriptor. This method does not require a Workspace.
func (s *MetricsService) GetMonitoredResourceDescriptor(ctx context.Context, req *GetMonitoredResourceDescriptorRequest) (*monitoredres.MonitoredResourceDescriptor, error) {
	log.Printf("GetMonitoredResourceDescriptor:%#v", req)
	resp := new(monitoredres.MonitoredResourceDescriptor)
	return resp, nil
}

// Lists metric descriptors that match a filter. This method does not require a Workspace.
func (s *MetricsService) ListMetricDescriptors(ctx context.Context, req *ListMetricDescriptorsRequest) (*ListMetricDescriptorsResponse, error) {
	log.Printf("ListMetricDescriptors:%#v", req)
	resp := &ListMetricDescriptorsResponse{}
	return resp, nil
}

// Gets a single metric descriptor. This method does not require a Workspace.
func (s *MetricsService) GetMetricDescriptor(ctx context.Context, req *GetMetricDescriptorRequest) (*metric.MetricDescriptor, error) {
	log.Printf("GetMetricDescriptor:%#v", req)
	resp := new(metric.MetricDescriptor)
	return resp, nil
}

// Creates a new metric descriptor.
// User-created metric descriptors define
// [custom metrics](https://cloud.google.com/monitoring/custom-metrics).
func (s *MetricsService) CreateMetricDescriptor(ctx context.Context, req *CreateMetricDescriptorRequest) (*metric.MetricDescriptor, error) {
	log.Printf("CreateMetricDescriptor:%#v", req)
	return nil, nil
}

// Deletes a metric descriptor. Only user-created
// [custom metrics](https://cloud.google.com/monitoring/custom-metrics) can be
// deleted.
func (s *MetricsService) DeleteMetricDescriptor(ctx context.Context, req *DeleteMetricDescriptorRequest) (*emptypb.Empty, error) {
	log.Printf("DeleteMetricDescriptor:%#v", req)
	return nil, nil
}

// Lists time series that match a filter. This method does not require a Workspace.
func (s *MetricsService) ListTimeSeries(ctx context.Context, req *ListTimeSeriesRequest) (*ListTimeSeriesResponse, error) {
	log.Printf("ListTimeSeries:%#v", req)
	return nil, nil
}

// Creates or adds data to one or more time series.
// The response is empty if all time series in the request were written.
// If any time series could not be written, a corresponding failure message is
// included in the error response.
func (s *MetricsService) CreateTimeSeries(ctx context.Context, req *CreateTimeSeriesRequest) (*emptypb.Empty, error) {
	log.Printf("CreateTimeSeries:%#v", req)
	return nil, nil
}
