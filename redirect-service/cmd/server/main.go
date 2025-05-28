package main

import (
	"context"
	"errors"
	"log"
	"net"
	"time"

	"github.com/Naman-B-Parlecha/url-shortener/redirect-service/redirect"
	"github.com/Naman-B-Parlecha/url-shortener/redirect-service/url"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	conn, err := grpc.NewClient("url-shortener-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to broker service: %v", err)
		return nil, err
	}
	defer conn.Close()

	gprcCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client := url.NewURLServiceClient(conn)
	resp, err := client.GetLongURL(gprcCtx, &url.ShortURL{Shorturl: req.Shorturl})
	if err != nil {
		log.Printf("Failed to get long URL from broker service: %v", err)
		return nil, err
	}

	if resp.Longurl == "" {
		log.Printf("No long URL found for short URL: %s", req.Shorturl)
		return nil, errors.New("no long URL found for the provided short URL")
	}

	return &redirect.LongURL{
		Longurl: resp.Longurl,
	}, nil
}
