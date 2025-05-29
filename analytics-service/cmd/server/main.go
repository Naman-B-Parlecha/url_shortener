package main

import (
	"context"
	"log"
	"net"

	"github.com/Naman-B-Parlecha/url-shortener/analytics-service/analytics"
	"github.com/Naman-B-Parlecha/url-shortener/analytics-service/db"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AnalyticsServer struct {
	analytics.UnimplementedAnalyticsServiceServer
	db *mongo.Client
}

func main() {

	db, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return
	}
	listen, err := net.Listen("tcp", ":50003")
	if err != nil {
		log.Fatalf("Failed to listen on port 50003: %v", err)
		return
	}
	s := grpc.NewServer()
	analytics.RegisterAnalyticsServiceServer(s, &AnalyticsServer{
		db: db,
	})

	log.Println("Analytics service is running on port 50003")
	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
		return
	}

}

func (s *AnalyticsServer) GetAnalytics(ctx context.Context, req *analytics.ShortURL) (*analytics.AnalyticsResponse, error) {
	resp := s.db.Database("UrlShort").Collection("Analytics").FindOne(ctx, map[string]interface{}{
		"shorturl": req.Shorturl,
	})
	if resp.Err() != nil {
		log.Printf("Error fetching analytics for short URL %s: %v", req.Shorturl, resp.Err())
		return nil, resp.Err()
	}

	var result struct {
		TotalClicks int32 `bson:"total_clicks"`
	}
	if err := resp.Decode(&result); err != nil {
		log.Printf("Error decoding analytics result: %v", err)
		return nil, err
	}

	return &analytics.AnalyticsResponse{
		Shorturl:    req.Shorturl,
		TotalClicks: result.TotalClicks,
	}, nil
}

func (s *AnalyticsServer) RecordAnalytics(ctx context.Context, req *analytics.ShortURL) (*emptypb.Empty, error) {
	_, err := s.db.Database("UrlShort").Collection("Analytics").UpdateOne(
		ctx,
		map[string]interface{}{"shorturl": req.Shorturl},
		map[string]interface{}{
			"$inc": map[string]interface{}{"total_clicks": 1},
		},
	)
	if err != nil {
		log.Printf("Error updating analytics for short URL %s: %v", req.Shorturl, err)
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

func (s *AnalyticsServer) AddAnalytics(ctx context.Context, req *analytics.ShortURL) (*emptypb.Empty, error) {
	_, err := s.db.Database("UrlShort").Collection("Analytics").InsertOne(ctx, map[string]interface{}{
		"shorturl":     req.Shorturl,
		"total_clicks": 0,
	})
	if err != nil {
		log.Printf("Error inserting analytics for short URL %s: %v", req.Shorturl, err)
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
