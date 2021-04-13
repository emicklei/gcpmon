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

func (s *EventStore) getTimeSeries(project string, metricType string) (list []*monitoring.TimeSeries) {
	v, ok := s.events.Load(project)
	if !ok {
		return list
	}
	pe := v.(*ProjectEvents)
	w, ok := pe.timeSeries.Load(metricType)
	if !ok {
		// keys := []string{}
		// pe.timeSeries.Range(func(k, v interface{}) bool {
		// 	keys = append(keys, k.(string))
		// 	return true
		// })
		// log.Println("metric", metricType, keys)
		return list
	}
	return w.([]*monitoring.TimeSeries)
}
