package main

import (
	"context"
	"log"
	"net"

	"github.com/Naman-B-Parlecha/url-shortener/redirect-service/redirect"
	"google.golang.org/grpc"
)

type RedirectServiceServer struct {
	redirect.UnimplementedRedirectServiceServer
}

func main() {
	listen, err := net.Listen("tcp", ":50002")
	if err != nil {
		log.Fatalf("Failed to listen on port 50002: %v", err)
		return
	}
	defer listen.Close()

	s := grpc.NewServer()
	redirect.RegisterRedirectServiceServer(s, &RedirectServiceServer{})

	log.Println("Redirect service gRPC server is running on port 50002")
	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
		return
	}
}

func (s *RedirectServiceServer) GetRedirectURL(ctx context.Context, req *redirect.ShortURL) (*redirect.LongURL, error) {
	return &redirect.LongURL{
		Longurl: "https://github.com/Naman-B-Parlecha",
	}, nil
}
