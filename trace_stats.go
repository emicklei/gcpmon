package main

import (
	"strconv"

	tvp "github.com/emicklei/tviewplus"
	"github.com/rivo/tview"
)

type TraceStats struct {
	Count *tvp.StringHolder
}

func NewTraceStats() *TraceStats {
	return &TraceStats{Count: new(tvp.StringHolder)}
}

func (s *TraceStats) update(mon *Monitor) {
	list := mon.store.getTracespans(mon.ProjectList.Selection.Value, mon.BatchWriteSpansList.Selection.Value)
	s.Count.Set(strconv.Itoa(len(list)))
}

func (s *TraceStats) addUITo(a *tview.Application, c *tview.Flex) {
	{
		row := tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(NewStaticView(" Count"), 12, 1, false).
			AddItem(tvp.NewReadOnlyTextView(a, s.Count), 12, 1, false)
		c.AddItem(row, 1, 1, false)
	}
}
