package main

import (
	"bytes"
	"fmt"

	tvp "github.com/emicklei/tviewplus"
	"google.golang.org/genproto/googleapis/api/label"
)

func (m *Monitor) setup() {
	m.MetricDescriptorList.AddSelectionChangeDependent(m.changedMetricDescriptor)
	m.ProjectList.AddSelectionChangeDependent(m.changedProject)
}

func (m *Monitor) changedProject(old, next tvp.SelectionWithIndex) {
	m.updateMetricDescriptors()
}

func (m *Monitor) changedMetricDescriptor(old, next tvp.SelectionWithIndex) {
	d := m.store.getMetricDescriptor(m.ProjectList.Selection.Value, next.Value)
	if d != nil {
		b := new(bytes.Buffer)
		for _, each := range d.Labels {
			fmt.Fprintf(b, "%s:%s (%s)\n", each.Key, each.GetDescription(), label.LabelDescriptor_ValueType_name[int32(each.ValueType)])
		}
		m.Labels.Set(b.String())
	}
	m.metricStats.update(m)
}
