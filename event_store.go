package main

import (
	"sync"

	"google.golang.org/genproto/googleapis/api/metric"
	"google.golang.org/genproto/googleapis/monitoring/v3"
)

type EventStore struct {
	events *sync.Map // project-id -> *ProjectEvents
}

func newEventStore() *EventStore { return &EventStore{new(sync.Map)} }

type ProjectEvents struct {
	metricDescriptors            *sync.Map
	monitoredResourceDescriptors *sync.Map
	timeSeries                   *sync.Map // type -> []TimeSeries
	traceSpans                   *sync.Map
}

func newProjectEvents() *ProjectEvents {
	return &ProjectEvents{
		metricDescriptors:            new(sync.Map),
		monitoredResourceDescriptors: new(sync.Map),
		timeSeries:                   new(sync.Map),
		traceSpans:                   new(sync.Map),
	}
}

func (s *EventStore) addMetricDescriptor(p string, d *metric.MetricDescriptor) {
	v, _ := s.events.LoadOrStore(p, newProjectEvents())
	pe := v.(*ProjectEvents)
	pe.metricDescriptors.Store(d.Type, d)
}

func (s *EventStore) addTimeSeries(p string, t []*monitoring.TimeSeries) {
	v, _ := s.events.LoadOrStore(p, newProjectEvents())
	pe := v.(*ProjectEvents)
	for _, each := range t {
		w, _ := pe.timeSeries.LoadOrStore(each.Metric.Type, []*monitoring.TimeSeries{})
		ts := w.([]*monitoring.TimeSeries)
		ts = append(ts, each)
		pe.timeSeries.Store(each.Metric.Type, ts)
	}
}

func (s *EventStore) getTimeSeries(project string, metricType string) (desc *metric.MetricDescriptor, list []*monitoring.TimeSeries) {
	v, ok := s.events.Load(project)
	if !ok {
		return nil, list
	}
	pe := v.(*ProjectEvents)
	w, ok := pe.timeSeries.Load(metricType)
	if !ok {
		return nil, list
	}
	d, ok := pe.metricDescriptors.Load(metricType)
	if !ok {
		return nil, list
	}
	return d.(*metric.MetricDescriptor), w.([]*monitoring.TimeSeries)
}

func (s *EventStore) getMetricDescriptor(project string, metricType string) *metric.MetricDescriptor {
	v, ok := s.events.Load(project)
	if !ok {
		return nil
	}
	pe := v.(*ProjectEvents)
	d, ok := pe.metricDescriptors.Load(metricType)
	if !ok {
		return nil
	}
	return d.(*metric.MetricDescriptor)
}
