module github.com/emicklei/gcpmon

go 1.16

require (
	contrib.go.opencensus.io/exporter/stackdriver v0.13.5
	github.com/rivo/tview v0.0.0-20210312174852-ae9464cc3598
	github.com/emicklei/tviewplus v0.7.2
	github.com/gdamore/tcell/v2 v2.2.0 
	go.opencensus.io v0.23.0
	google.golang.org/api v0.44.0
	google.golang.org/genproto v0.0.0-20210406143921-e86de6bf7a46
	google.golang.org/grpc v1.36.1
	google.golang.org/protobuf v1.26.0
)

// replace github.com/emicklei/tviewplus => ../tviewplus