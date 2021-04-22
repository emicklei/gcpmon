package main

import (
	"sync"
	"time"

	"google.golang.org/genproto/googleapis/api/metric"
	"google.golang.org/genproto/googleapis/devtools/cloudtrace/v2"
	"google.golang.org/genproto/googleapis/monitoring/v3"
)

type EventStore struct {
	events *sync.Map // project-id -> *ProjectEvents
}

func newEventStore() *EventStore { return &EventStore{new(sync.Map)} }

type ProjectEvents struct {
	metricDescriptors            *sync.Map
	monitoredResourceDescriptors *sync.Map
	timeSeriesStartedAt          time.Time
	timeSeries                   *sync.Map // type -> []TimeSeries
	traceSpans                   *sync.Map
}

func newProjectEvents() *ProjectEvents {
	return &ProjectEvents{
		metricDescriptors:            new(sync.Map),
		monitoredResourceDescriptors: new(sync.Map),
		timeSeriesStartedAt:          time.Now(),
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
		if len(ts) == 0 {
			pe.timeSeriesStartedAt = time.Now()
		}
		ts = append(ts, each)
		pe.timeSeries.Store(each.Metric.Type, ts)
	}
}

func (s *EventStore) getTracespans(project string, displayName string) (list []*cloudtrace.Span) {
	v, ok := s.events.Load(project)
	if !ok {
		return list
	}
	pe := v.(*ProjectEvents)
	w, ok := pe.traceSpans.Load(displayName)
	if !ok {
		return list
	}
	return w.([]*cloudtrace.Span)
}

func (s *EventStore) getTimeSeries(project string, metricType string) (desc *metric.MetricDescriptor, list []*monitoring.TimeSeries, t time.Time) {
	v, ok := s.events.Load(project)
	if !ok {
		return nil, list, time.Time{}
	}
	pe := v.(*ProjectEvents)
	w, ok := pe.timeSeries.Load(metricType)
	if !ok {
		return nil, list, pe.timeSeriesStartedAt
	}
	d, ok := pe.metricDescriptors.Load(metricType)
	if !ok {
		return nil, list, pe.timeSeriesStartedAt
	}
	return d.(*metric.MetricDescriptor), w.([]*monitoring.TimeSeries), pe.timeSeriesStartedAt
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

func (s *EventStore) addTraceSpans(project string, spans []*cloudtrace.Span) {
	v, _ := s.events.LoadOrStore(project, newProjectEvents())
	pe := v.(*ProjectEvents)
	// use pe lock instead of syncmap?
	for _, each := range spans {
		key := each.GetDisplayName().GetValue()
		w, _ := pe.traceSpans.LoadOrStore(key, []*cloudtrace.Span{})
		list := w.([]*cloudtrace.Span)
		list = append(list, each)
		pe.traceSpans.Store(key, w)
	}
}
