package main

import tvp "github.com/emicklei/tviewplus"

func (m *Monitor) setup() {
	m.MetricDescriptorList.AddSelectionChangeDependent(m.changedMetricDescriptor)
}

func (m *Monitor) changedMetricDescriptor(old, new tvp.SelectionWithIndex) {
	m.metricStats.update(m)
}
