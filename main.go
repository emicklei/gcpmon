package main

import (
	"log"
	"net"

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

	svc := new(MetricsService)
	grpcServer := grpc.NewServer(grpc.Creds(cred))
	monitoringpb.RegisterMetricServiceServer(grpcServer, svc)
	log.Println("serving gRPC on :9443")
	grpcServer.Serve(lis)
}
