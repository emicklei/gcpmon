package main

import (
	"log"

	tvp "github.com/emicklei/tviewplus"
	"github.com/rivo/tview"
)

const noSelection = "---"

func start(mon *Monitor) {
	// textBg := tcell.NewRGBColor(25, 28, 32)
	// dropBg := tcell.NewRGBColor(20, 23, 27)
	// tview.Styles.PrimaryTextColor = tcell.ColorGray
	// tview.Styles.ContrastBackgroundColor = dropBg

	app := tview.NewApplication()
	foc := tvp.NewFocusGroup(app)

	projects := tvp.NewListView(foc, mon.ProjectList)
	projects.SetBorder(true)
	// change item makes new selection
	projects.SetChangedFunc(func(index int, mainText, secondText string, shortcut rune) {
		mon.ProjectList.Select(index)
	})

	metricDescriptors := tvp.NewListView(foc, mon.MetricDescriptorList)
	metricDescriptors.SetBorder(true)
	metricDescriptors.SetHighlightFullLine(true)
	// change item makes new selection
	metricDescriptors.SetChangedFunc(func(index int, mainText, secondText string, shortcut rune) {
		mon.MetricDescriptorList.Select(index)
	})

	labels := tvp.NewReadOnlyTextView(app, mon.Labels)
	labels.SetBorder(true)

	spans := tvp.NewListView(foc, mon.BatchWriteSpansList)
	spans.SetBorder(true)
	spans.SetHighlightFullLine(true)

	console := tvp.NewReadOnlyTextView(app, mon.Console)
	//console.SetTextColor(tcell.ColorLightGray)
	console.SetBorder(true)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(NewStaticView(" [yellow]gcpmon - Google Cloud Monitoring Inspector"), 1, 1, false).
		// proj
		AddItem(tview.NewBox().SetBorderPadding(1, 0, 0, 0), 1, 1, false).
		AddItem(NewStaticView(" [green]projects"), 1, 1, false).
		AddItem(projects, 4, 1, true).
		// desc
		AddItem(tview.NewBox().SetBorderPadding(1, 0, 0, 0), 1, 1, false).
		AddItem(NewStaticView(" [green]metrics"), 1, 1, false).
		AddItem(metricDescriptors, 6, 1, false)

	// labels
	flex.
		AddItem(tview.NewBox().SetBorderPadding(1, 0, 0, 0), 1, 1, false).
		AddItem(NewStaticView(" [green]metric label definitions"), 1, 1, false).
		AddItem(labels, 6, 1, false)

	flex.AddItem(NewStaticView(" [green]traces"), 1, 1, false)
	flex.AddItem(spans, 4, 1, false)

	flex.AddItem(NewStaticView(" [green]metric stats"), 1, 1, false)
	mon.metricStats.addUITo(app, flex)

	flex.AddItem(NewStaticView(" [green]trace stats"), 1, 1, false)
	mon.traceStats.addUITo(app, flex)

	// console
	flex.
		AddItem(tview.NewBox().SetBorderPadding(1, 0, 0, 0), 1, 1, false).
		AddItem(NewStaticView(" [yellow]console"), 1, 1, false).
		AddItem(console, 0, 4, false)

	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		log.Println(err)
	}
}

func NewStaticView(label string) *tview.TextView {
	return tview.NewTextView().SetDynamicColors(true).SetText(label)
}
