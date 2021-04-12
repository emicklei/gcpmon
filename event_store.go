package main

import (
	"sync"

	"google.golang.org/genproto/googleapis/api/metric"
	"google.golang.org/genproto/googleapis/monitoring/v3"
)

type EventStore struct {
	events sync.Map // project-id -> *ProjectEvents
}

type ProjectEvents struct {
	metricDescriptors            sync.Map
	monitoredResourceDescriptors sync.Map
	timeSeries                   sync.Map // type -> []TimeSeries
	traceSpans                   sync.Map
}

func (s *EventStore) addMetricDescriptor(p string, d *metric.MetricDescriptor) {
	v, _ := s.events.LoadOrStore(p, new(ProjectEvents))
	pe := v.(*ProjectEvents)
	pe.metricDescriptors.Store(d.Name, d)
}

func (s *EventStore) addTimeSeries(p string, t []*monitoring.TimeSeries) {
	v, _ := s.events.LoadOrStore(p, new(ProjectEvents))
	pe := v.(*ProjectEvents)
	for _, each := range t {
		w, _ := pe.timeSeries.LoadOrStore(each.Metric.Type, []*monitoring.TimeSeries{})
		ts := w.([]*monitoring.TimeSeries)
		ts = append(ts, each)
	}
}
