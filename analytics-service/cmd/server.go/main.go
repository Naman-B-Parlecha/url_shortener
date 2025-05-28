package main

import (
	"log"
	"net"

	"github.com/Naman-B-Parlecha/url-shortener/analytics-service/analytics"
	"google.golang.org/grpc"
)

type AnalyticsServer struct {
	analytics.UnimplementedAnalyticsServiceServer
}

func main() {
	listen, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Failed to listen on port 50053: %v", err)
		return
	}
	s := grpc.NewServer()
	analytics.RegisterAnalyticsServiceServer(s, &AnalyticsServer{})

	log.Println("Analytics service is running on port 50053")
	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
		return
	}
}
