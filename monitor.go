package main

import (
	"fmt"

	tvp "github.com/emicklei/tviewplus"
)

type Monitor struct {
	ProjectList          *tvp.StringListSelectionHolder
	MetricDescriptorList *tvp.StringListSelectionHolder
	BatchWriteSpansList  *tvp.StringListSelectionHolder
	Console              *tvp.StringHolder
}

func NewMonitor() *Monitor {
	return &Monitor{
		ProjectList:          new(tvp.StringListSelectionHolder),
		MetricDescriptorList: new(tvp.StringListSelectionHolder),
		Console:              new(tvp.StringHolder),
	}
}

func (m *Monitor) Printf(format string, v ...interface{}) {
	m.Console.Append(fmt.Sprintf(format, v...))
	// log.Printf(format, v...)
}
