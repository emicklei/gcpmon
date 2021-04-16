package main

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

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
	LatestLabels    *tvp.StringHolder
}

func NewMetricStats() *MetricStats {
	return &MetricStats{
		MetricKind:      new(tvp.StringHolder),
		MetricValueType: new(tvp.StringHolder),
		Frequency:       new(tvp.StringHolder),
		MinValue:        new(tvp.StringHolder),
		MaxValue:        new(tvp.StringHolder),
		Count:           new(tvp.StringHolder),
		LatestLabels:    new(tvp.StringHolder),
	}
}

func (s *MetricStats) update(mon *Monitor) {
	desc, ts, when := mon.store.getTimeSeries(mon.ProjectList.Selection.Value, mon.MetricDescriptorList.Selection.Value)
	if len(ts) == 0 {
		s.Count.Set("")
		s.MinValue.Set("")
		s.MaxValue.Set("")
		return
	}
	//count := 0
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
				if desc.GetValueType() == metric.MetricDescriptor_INT64 {
					f := float64(other.Value.GetInt64Value())
					if f < min {
						min = f
					}
					if f > max {
						max = f
					}
				} else {
					logger.Println("unhandled point type", each.ValueType)
				}
			}
		}
	}
	buf := new(bytes.Buffer)
	latest := ts[len(ts)-1]
	if len(latest.GetResource().Labels) > 0 {
		fmt.Fprintf(buf, "resource.labels values:\n")
		for k, v := range latest.GetResource().Labels {
			fmt.Fprintf(buf, "\t%s:%v\n", k, v)
		}
	}
	if latest.Metadata != nil {
		fmt.Fprintf(buf, "metadata.userlabels values:\n")
		for k, v := range latest.Metadata.UserLabels {
			fmt.Fprintf(buf, "\t%s:%v\n", k, v)
		}
		if latest.Metadata.SystemLabels != nil {
			fmt.Fprintf(buf, "metadata.systemlabels values:\n")
			for k, v := range latest.Metadata.SystemLabels.Fields {
				fmt.Fprintf(buf, "\t%s:%v\n", k, v)
			}
		}
	}
	if latest.Metric != nil {
		fmt.Fprintf(buf, "metric labels values:\n")
		for k, v := range latest.Metric.Labels {
			fmt.Fprintf(buf, "\t%s:%v\n", k, v)
		}
	}
	s.LatestLabels.Set(buf.String())
	s.Count.Set(strconv.Itoa(len(ts)))
	dur := time.Now().Sub(when).Seconds()
	freq := float64(len(ts)) / float64(dur)
	s.Frequency.Set(fmt.Sprintf("%.2f", freq))
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
	{
		rov := tvp.NewReadOnlyTextView(a, s.LatestLabels)
		rov.SetBorder(true)
		c.AddItem(rov, 6, 1, false)
	}
}
