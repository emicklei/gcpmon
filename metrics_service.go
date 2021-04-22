package main

import (
	"context"
	"fmt"
	"sync"

	"google.golang.org/genproto/googleapis/api/metric"
	"google.golang.org/genproto/googleapis/api/monitoredres"
	. "google.golang.org/genproto/googleapis/monitoring/v3"
	"google.golang.org/protobuf/types/known/emptypb"
)

// implements monitoringpb.MetricServiceServer
type MetricsService struct {
	metricDescriptors            sync.Map
	monitoredResourceDescriptors sync.Map
	monitor                      *Monitor
}

// Lists monitored resource descriptors that match a filter. This method does not require a Workspace.
func (s *MetricsService) ListMonitoredResourceDescriptors(ctx context.Context, req *ListMonitoredResourceDescriptorsRequest) (*ListMonitoredResourceDescriptorsResponse, error) {
	s.monitor.Printf("ListMonitoredResourceDescriptors:%#v", req)
	resp := &ListMonitoredResourceDescriptorsResponse{}
	s.metricDescriptors.Range(func(k, v interface{}) bool {
		resp.ResourceDescriptors = append(resp.ResourceDescriptors, v.(*monitoredres.MonitoredResourceDescriptor))
		return true
	})
	return resp, nil
}

// Gets a single monitored resource descriptor. This method does not require a Workspace.
func (s *MetricsService) GetMonitoredResourceDescriptor(ctx context.Context, req *GetMonitoredResourceDescriptorRequest) (*monitoredres.MonitoredResourceDescriptor, error) {
	s.monitor.Printf("GetMonitoredResourceDescriptor:%#v", req)
	if v, ok := s.monitoredResourceDescriptors.Load(req.Name); ok {
		return v.(*monitoredres.MonitoredResourceDescriptor), nil
	}
	return nil, fmt.Errorf("no such monitored resource descriptor:%s", req.Name)
}

// Lists metric descriptors that match a filter. This method does not require a Workspace.
func (s *MetricsService) ListMetricDescriptors(ctx context.Context, req *ListMetricDescriptorsRequest) (*ListMetricDescriptorsResponse, error) {
	s.monitor.Printf("ListMetricDescriptors:%#v", req)
	resp := &ListMetricDescriptorsResponse{}
	s.metricDescriptors.Range(func(k, v interface{}) bool {
		resp.MetricDescriptors = append(resp.MetricDescriptors, v.(*metric.MetricDescriptor))
		return true
	})
	return resp, nil
}

// Gets a single metric descriptor. This method does not require a Workspace.
func (s *MetricsService) GetMetricDescriptor(ctx context.Context, req *GetMetricDescriptorRequest) (*metric.MetricDescriptor, error) {
	s.monitor.Printf("GetMetricDescriptor:%s", req.Name)
	if v, ok := s.metricDescriptors.Load(req.Name); ok {
		return v.(*metric.MetricDescriptor), nil
	}
	return nil, fmt.Errorf("no such metric descriptor:%s", req.Name)
}

// Creates a new metric descriptor.
// User-created metric descriptors define
// [custom metrics](https://cloud.google.com/monitoring/custom-metrics).
func (s *MetricsService) CreateMetricDescriptor(ctx context.Context, req *CreateMetricDescriptorRequest) (*metric.MetricDescriptor, error) {
	s.monitor.store.addMetricDescriptor(req.Name, req.MetricDescriptor)
	go func() {
		s.monitor.updateProjects()
		s.monitor.updateMetricDescriptors()
	}()
	return req.MetricDescriptor, nil
}

// Deletes a metric descriptor. Only user-created
// [custom metrics](https://cloud.google.com/monitoring/custom-metrics) can be
// deleted.
func (s *MetricsService) DeleteMetricDescriptor(ctx context.Context, req *DeleteMetricDescriptorRequest) (*emptypb.Empty, error) {
	s.monitor.Printf("DeleteMetricDescriptor:%s", req.Name)
	s.metricDescriptors.Delete(req.Name)
	return new(emptypb.Empty), nil
}

// Lists time series that match a filter. This method does not require a Workspace.
func (s *MetricsService) ListTimeSeries(ctx context.Context, req *ListTimeSeriesRequest) (*ListTimeSeriesResponse, error) {
	s.monitor.Printf("ListTimeSeries:%s", req.Name)
	resp := &ListTimeSeriesResponse{}
	return resp, nil
}

// Creates or adds data to one or more time series.
// The response is empty if all time series in the request were written.
// If any time series could not be written, a corresponding failure message is
// included in the error response.
func (s *MetricsService) CreateTimeSeries(ctx context.Context, req *CreateTimeSeriesRequest) (*emptypb.Empty, error) {
	//s.monitor.Printf("!")
	s.monitor.store.addTimeSeries(req.Name, req.TimeSeries)
	go func() {
		s.monitor.updateMetricStats()
	}()
	return new(emptypb.Empty), nil
}

func intervalDisplay(i *TimeInterval) string {
	return fmt.Sprintf("start:%d end:%d", i.StartTime.Seconds, i.EndTime.Seconds)
}

func pointDisplay(v *TypedValue) string {
	if d := v.GetDistributionValue(); d != nil {
		return fmt.Sprintf("mean:%v count:%d", d.Mean, d.Count)
	}
	return v.String()
}
