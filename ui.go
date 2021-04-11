package main

import (
	"log"

	tvp "github.com/emicklei/tviewplus"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func start(mon *Monitor) {
	textBg := tcell.NewRGBColor(25, 28, 32)
	dropBg := tcell.NewRGBColor(20, 23, 27)
	tview.Styles.PrimaryTextColor = tcell.ColorGray
	tview.Styles.ContrastBackgroundColor = dropBg

	app := tview.NewApplication()
	foc := tvp.NewFocusGroup(app)

	projects := tvp.NewDropDownView(foc, mon.ProjectList)
	projects.SetTextOptions("", "", "", "▼", "---")
	projects.SetLabel(" projects ")

	metricDescriptors := tvp.NewDropDownView(foc, mon.MetricDescriptorList)
	metricDescriptors.SetTextOptions("", "", "", "▼", "---")
	metricDescriptors.SetLabel(" metric descriptors ")

	console := tvp.NewReadOnlyTextView(app, mon.Console)
	console.SetTextColor(tcell.ColorLightGray)
	console.SetBackgroundColor(textBg)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(NewStaticView(" [yellow]gcpmon - Google Cloud Monitoring Inspector"), 0, 1, false).
		AddItem(projects, 1, 1, true).
		AddItem(metricDescriptors, 1, 1, false).

		// console
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
