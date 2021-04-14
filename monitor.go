package main

import (
	"fmt"

	tvp "github.com/emicklei/tviewplus"
)

type Monitor struct {
	store                *EventStore
	metricStats          *MetricStats
	ProjectList          *tvp.StringListSelectionHolder
	MetricDescriptorList *tvp.StringListSelectionHolder
	BatchWriteSpansList  *tvp.StringListSelectionHolder
	Labels               *tvp.StringHolder
	Console              *tvp.StringHolder
}

func NewMonitor(s *EventStore) *Monitor {
	return &Monitor{
		store:                s,
		metricStats:          NewMetricStats(),
		ProjectList:          new(tvp.StringListSelectionHolder),
		MetricDescriptorList: new(tvp.StringListSelectionHolder),
		Console:              new(tvp.StringHolder),
		Labels:               new(tvp.StringHolder),
	}
}

func (m *Monitor) Printf(format string, v ...interface{}) {
	m.Console.Append(fmt.Sprintf(format, v...))
}
func (m *Monitor) Println(v ...interface{}) {
	m.Console.Append(fmt.Sprintln(v...))
}

func (m *Monitor) updateProjects() {
	names := []string{}
	m.store.events.Range(func(k, _ interface{}) bool {
		names = append(names, k.(string))
		return true
	})
	m.ProjectList.Set(names)
	if len(names) == 1 {
		m.ProjectList.Select(0)
	}
}

func (m *Monitor) updateMetricDescriptors() {
	p := m.ProjectList.Selection.Value
	if p == noSelection {
		m.MetricDescriptorList.Set([]string{})
		return
	}
	v, ok := m.store.events.Load(p)
	if !ok {
		return
	}
	pe := v.(*ProjectEvents)
	names := []string{}
	pe.metricDescriptors.Range(func(k, _ interface{}) bool {
		names = append(names, k.(string))
		return true
	})
	m.MetricDescriptorList.Set(names)
}

func (m *Monitor) updateMetricStats() {
	m.metricStats.update(m)
}
