package main

import (
	"fmt"
	"strconv"

	"math"

	tvp "github.com/emicklei/tviewplus"
	"github.com/rivo/tview"
	"google.golang.org/genproto/googleapis/api/metric"
)

type MetricStats struct {
	MetricKind      *tvp.StringHolder
	MetricValueType *tvp.StringHolder
	Frequency       *tvp.StringHolder
	MinValue        *tvp.StringHolder
	MaxValue        *tvp.StringHolder
	Count           *tvp.StringHolder
}

func NewMetricStats() *MetricStats {
	return &MetricStats{
		MetricKind:      new(tvp.StringHolder),
		MetricValueType: new(tvp.StringHolder),
		Frequency:       new(tvp.StringHolder),
		MinValue:        new(tvp.StringHolder),
		MaxValue:        new(tvp.StringHolder),
		Count:           new(tvp.StringHolder),
	}
}

func (s *MetricStats) update(mon *Monitor) {
	desc, ts := mon.store.getTimeSeries(mon.ProjectList.Selection.Value, mon.MetricDescriptorList.Selection.Value)
	if len(ts) == 0 {
		s.Count.Set("")
		s.MinValue.Set("")
		s.MaxValue.Set("")
		return
	}
	count := 0
	var min float64 = math.MaxFloat64
	var max float64 = -min
	s.MetricKind.Set(metric.MetricDescriptor_MetricKind_name[int32(desc.MetricKind)])
	s.MetricValueType.Set(metric.MetricDescriptor_ValueType_name[int32(desc.GetValueType())])
	for _, each := range ts {
		for _, other := range each.Points {
			if dist := other.Value.GetDistributionValue(); dist != nil {
				//count = count + int(dist.Count)
				mean := dist.Mean
				if mean < min {
					min = mean
				}
				if mean > max {
					max = mean
				}
			} else {
				logger.Println("point type", each.ValueType)
			}
		}
	}
	s.Count.Set(strconv.Itoa(count))
	s.MinValue.Set(fmt.Sprintf("%.4f", min))
	s.MaxValue.Set(fmt.Sprintf("%.4f", max))
}

func (s *MetricStats) addUITo(a *tview.Application, c *tview.Flex) {
	{
		row := tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(NewStaticView(" Kind"), 12, 1, false).
			AddItem(tvp.NewReadOnlyTextView(a, s.MetricKind), 12, 1, false)
		c.AddItem(row, 1, 1, false)
	}
	{
		row := tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(NewStaticView(" Type"), 12, 1, false).
			AddItem(tvp.NewReadOnlyTextView(a, s.MetricValueType), 12, 1, false)
		c.AddItem(row, 1, 1, false)
	}
	{
		row := tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(NewStaticView(" Count"), 12, 1, false).
			AddItem(tvp.NewReadOnlyTextView(a, s.Count), 12, 1, false)
		c.AddItem(row, 1, 1, false)
	}
	{
		row := tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(NewStaticView(" Frequency"), 12, 1, false).
			AddItem(tvp.NewReadOnlyTextView(a, s.Frequency), 12, 1, false)
		c.AddItem(row, 1, 1, false)
	}
	{
		row := tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(NewStaticView(" Minimum"), 12, 1, false).
			AddItem(tvp.NewReadOnlyTextView(a, s.MinValue), 12, 1, false)
		c.AddItem(row, 1, 1, false)
	}
	{
		row := tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(NewStaticView(" Maximum"), 12, 1, false).
			AddItem(tvp.NewReadOnlyTextView(a, s.MaxValue), 12, 1, false)
		c.AddItem(row, 1, 1, false)
	}

}
