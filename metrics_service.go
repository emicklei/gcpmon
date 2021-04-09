package main

import (
	"context"
	"fmt"
	"log"
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
}

// Lists monitored resource descriptors that match a filter. This method does not require a Workspace.
func (s *MetricsService) ListMonitoredResourceDescriptors(ctx context.Context, req *ListMonitoredResourceDescriptorsRequest) (*ListMonitoredResourceDescriptorsResponse, error) {
	log.Printf("ListMonitoredResourceDescriptors:%#v", req)
	resp := &ListMonitoredResourceDescriptorsResponse{}
	s.metricDescriptors.Range(func(k, v interface{}) bool {
		resp.ResourceDescriptors = append(resp.ResourceDescriptors, v.(*monitoredres.MonitoredResourceDescriptor))
		return true
	})
	return resp, nil
}

// Gets a single monitored resource descriptor. This method does not require a Workspace.
func (s *MetricsService) GetMonitoredResourceDescriptor(ctx context.Context, req *GetMonitoredResourceDescriptorRequest) (*monitoredres.MonitoredResourceDescriptor, error) {
	log.Printf("GetMonitoredResourceDescriptor:%#v", req)
	if v, ok := s.monitoredResourceDescriptors.Load(req.Name); ok {
		return v.(*monitoredres.MonitoredResourceDescriptor), nil
	}
	return nil, fmt.Errorf("no such monitored resource descriptor:%s", req.Name)
}

// Lists metric descriptors that match a filter. This method does not require a Workspace.
func (s *MetricsService) ListMetricDescriptors(ctx context.Context, req *ListMetricDescriptorsRequest) (*ListMetricDescriptorsResponse, error) {
	log.Printf("ListMetricDescriptors:%#v", req)
	resp := &ListMetricDescriptorsResponse{}
	s.metricDescriptors.Range(func(k, v interface{}) bool {
		resp.MetricDescriptors = append(resp.MetricDescriptors, v.(*metric.MetricDescriptor))
		return true
	})
	return resp, nil
}

// Gets a single metric descriptor. This method does not require a Workspace.
func (s *MetricsService) GetMetricDescriptor(ctx context.Context, req *GetMetricDescriptorRequest) (*metric.MetricDescriptor, error) {
	log.Printf("GetMetricDescriptor:%s", req.Name)
	if v, ok := s.metricDescriptors.Load(req.Name); ok {
		return v.(*metric.MetricDescriptor), nil
	}
	return nil, fmt.Errorf("no such metric descriptor:%s", req.Name)
}

// Creates a new metric descriptor.
// User-created metric descriptors define
// [custom metrics](https://cloud.google.com/monitoring/custom-metrics).
func (s *MetricsService) CreateMetricDescriptor(ctx context.Context, req *CreateMetricDescriptorRequest) (*metric.MetricDescriptor, error) {
	log.Printf("CreateMetricDescriptor:%s desc:%s", req.Name, req.MetricDescriptor.Name)
	s.metricDescriptors.Store(req.Name, req.MetricDescriptor)
	return req.MetricDescriptor, nil
}

// Deletes a metric descriptor. Only user-created
// [custom metrics](https://cloud.google.com/monitoring/custom-metrics) can be
// deleted.
func (s *MetricsService) DeleteMetricDescriptor(ctx context.Context, req *DeleteMetricDescriptorRequest) (*emptypb.Empty, error) {
	log.Printf("DeleteMetricDescriptor:%s", req.Name)
	s.metricDescriptors.Delete(req.Name)
	return new(emptypb.Empty), nil
}

// Lists time series that match a filter. This method does not require a Workspace.
func (s *MetricsService) ListTimeSeries(ctx context.Context, req *ListTimeSeriesRequest) (*ListTimeSeriesResponse, error) {
	log.Printf("ListTimeSeries:%s", req.Name)
	resp := &ListTimeSeriesResponse{}
	return resp, nil
}

// Creates or adds data to one or more time series.
// The response is empty if all time series in the request were written.
// If any time series could not be written, a corresponding failure message is
// included in the error response.
func (s *MetricsService) CreateTimeSeries(ctx context.Context, req *CreateTimeSeriesRequest) (*emptypb.Empty, error) {
	log.Printf("CreateTimeSeries:%s len:%d", req.Name, len(req.TimeSeries))
	for _, each := range req.TimeSeries {
		log.Printf("\ttype:%s\n", each.Metric.Type)
		log.Printf("\tresource:%v labels:%v unit:%v points:%d\n", each.Resource.Type, each.Metric.Labels, each.Unit, len(each.Points))
		for _, other := range each.Points {
			log.Printf("\t\tinterval:%s value:%s\n", intervalDisplay(other.Interval), pointDisplay(other.GetValue()))
		}
	}
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
