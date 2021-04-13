package main

import (
	"fmt"
	"strconv"

	"math"

	tvp "github.com/emicklei/tviewplus"
	"github.com/rivo/tview"
)

type MetricStats struct {
	Frequency *tvp.StringHolder
	MinValue  *tvp.StringHolder
	MaxValue  *tvp.StringHolder
	Count     *tvp.StringHolder
}

func NewMetricStats() *MetricStats {
	return &MetricStats{
		Frequency: new(tvp.StringHolder),
		MinValue:  new(tvp.StringHolder),
		MaxValue:  new(tvp.StringHolder),
		Count:     new(tvp.StringHolder),
	}
}

func (s *MetricStats) update(mon *Monitor) {
	ts := mon.store.getTimeSeries(mon.ProjectList.Selection.Value, mon.MetricDescriptorList.Selection.Value)
	if len(ts) == 0 {
		s.Count.Set("")
		s.MinValue.Set("")
		s.MaxValue.Set("")
		return
	}
	s.Count.Set(strconv.Itoa(len(ts)))
	var min int64 = math.MaxInt64
	var max int64 = -min
	for _, each := range ts {
		for _, other := range each.Points {
			if v := other.Value.GetInt64Value(); v > max {
				max = v
			} else {
				if v < min {
					min = v
				}
			}
		}
	}
	s.MinValue.Set(fmt.Sprintf("%d", min))
	s.MaxValue.Set(fmt.Sprintf("%d", max))
}

func (s *MetricStats) addUITo(a *tview.Application, c *tview.Flex) {
	{
		row := tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(NewStaticView(" Count"), 20, 1, false).
			AddItem(tvp.NewReadOnlyTextView(a, s.Count), 6, 1, false)
		c.AddItem(row, 1, 1, false)
	}
	{
		row := tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(NewStaticView(" Frequency"), 20, 1, false).
			AddItem(tvp.NewReadOnlyTextView(a, s.Frequency), 6, 1, false)
		c.AddItem(row, 1, 1, false)
	}
	{
		row := tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(NewStaticView(" Minimum"), 20, 1, false).
			AddItem(tvp.NewReadOnlyTextView(a, s.MinValue), 6, 1, false)
		c.AddItem(row, 1, 1, false)
	}
	{
		row := tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(NewStaticView(" Maximum"), 20, 1, false).
			AddItem(tvp.NewReadOnlyTextView(a, s.MaxValue), 6, 1, false)
		c.AddItem(row, 1, 1, false)
	}

}
