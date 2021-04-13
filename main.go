package main

import (
	"log"
	"net"

	"google.golang.org/genproto/googleapis/devtools/cloudtrace/v2"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	lis, err := net.Listen("tcp", ":9443")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	cred, err := credentials.NewServerTLSFromFile("cert/server.crt", "cert/server.key")
	if err != nil {
		log.Fatal("failed to load TLS", err)
	}

	mon := NewMonitor(newEventStore())
	mon.setup()

	ms := &MetricsService{monitor: mon}
	ts := &TraceService{monitor: mon}
	grpcServer := grpc.NewServer(grpc.Creds(cred))
	log.Println("register metrics service")
	monitoringpb.RegisterMetricServiceServer(grpcServer, ms)
	log.Println("register tracing service")
	cloudtrace.RegisterTraceServiceServer(grpcServer, ts)
	log.Println("serving gRPC on :9443")
	go grpcServer.Serve(lis)

	start(mon)
}
