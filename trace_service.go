package main

import (
	"context"

	. "google.golang.org/genproto/googleapis/devtools/cloudtrace/v2"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TraceService struct {
	monitor *Monitor
}

// Sends new spans to new or existing traces. You cannot update
// existing spans.
func (s *TraceService) BatchWriteSpans(ctx context.Context, req *BatchWriteSpansRequest) (*emptypb.Empty, error) {
	// s.monitor.Printf("BatchWriteSpans:%s", req.Name)
	// for _, each := range req.Spans {
	// 	s.monitor.Printf("\tspan %v with delta:%d", each.DisplayName, each.EndTime.Nanos-each.StartTime.GetNanos())
	// }
	return new(emptypb.Empty), nil
}

// Creates a new span.
func (s *TraceService) CreateSpan(ctx context.Context, req *Span) (*Span, error) {
	//s.monitor.Printf("CreateSpan:%s", req.Name)
	return req, nil
}
