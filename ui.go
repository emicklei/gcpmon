package main

import (
	"log"

	tvp "github.com/emicklei/tviewplus"
	"github.com/gdamore/tcell/v2"
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
	//projects.S SetLabel(" projects")

	metricDescriptors := tvp.NewListView(foc, mon.MetricDescriptorList)
	metricDescriptors.SetBorder(true)
	//metricDescriptors.SetLabel(" metrics")

	console := tvp.NewReadOnlyTextView(app, mon.Console)
	console.SetTextColor(tcell.ColorLightGray)
	console.SetBorder(true)
	//console.SetBackgroundColor(textBg)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(NewStaticView(" [yellow]gcpmon - Google Cloud Monitoring Inspector"), 1, 1, false).
		// proj
		AddItem(tview.NewBox().SetBorderPadding(1, 0, 0, 0), 1, 1, false).
		AddItem(NewStaticView(" [green]projects"), 1, 1, false).
		AddItem(projects, 3, 1, true).
		// desc
		AddItem(tview.NewBox().SetBorderPadding(1, 0, 0, 0), 1, 1, false).
		AddItem(NewStaticView(" [green] metrics"), 1, 1, false).
		AddItem(metricDescriptors, 9, 1, false).
		// console
		AddItem(tview.NewBox().SetBorderPadding(1, 0, 0, 0), 1, 1, false).
		AddItem(NewStaticView(" [yellow]console"), 1, 1, false).
		AddItem(console, 0, 4, false)

	mon.metricStats.addUITo(app, flex)

	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		log.Println(err)
	}
}

func NewStaticView(label string) *tview.TextView {
	return tview.NewTextView().SetDynamicColors(true).SetText(label)
}
