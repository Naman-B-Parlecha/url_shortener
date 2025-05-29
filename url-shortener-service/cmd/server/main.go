package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net"

	"github.com/Naman-B-Parlecha/url-shortener/shortener-service/analytics"
	"github.com/Naman-B-Parlecha/url-shortener/shortener-service/db"
	"github.com/Naman-B-Parlecha/url-shortener/shortener-service/url"
	"github.com/Naman-B-Parlecha/url-shortener/shortener-service/util"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCServer struct {
	url.UnimplementedURLServiceServer
	DBConn *sql.DB
}

func main() {
	dbConn, err := db.ConnectDB()
	defer dbConn.Close()
	if err != nil {
		log.Println("Failed to connect to the database:", err)
		return
	}
	_, err = dbConn.Exec(`
	CREATE TABLE IF NOT EXISTS urls (
		id UUID PRIMARY KEY,
		long_url TEXT UNIQUE NOT NULL,
		short_url TEXT UNIQUE NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`)
	if err != nil {
		log.Println("Failed to create table:", err)
		panic(err)
	}
	listen, err := net.Listen("tcp", ":50001")
	if err != nil {
		log.Fatalln("Failed to listen: ", err)
		return
	}
	defer listen.Close()

	s := grpc.NewServer()
	url.RegisterURLServiceServer(s, &GRPCServer{
		DBConn: dbConn,
	})
	log.Println("Starting gRPC server on port :50001")

	if err := s.Serve(listen); err != nil {
		log.Fatalln("Failed to serve: ", err)
		return
	}
}

func (s *GRPCServer) GenerateShortURL(ctx context.Context, req *url.LongURL) (*url.ShortURL, error) {
	id := uuid.New()
	shortu_url := util.HashToBase62(req.Longurl, 5)
	query := `INSERT INTO urls (id, long_url, short_url) VALUES ($1, $2, $3)`
	_, err := s.DBConn.ExecContext(ctx, query, id, req.Longurl, shortu_url)
	if err != nil {
		log.Println("Failed to insert URL into database:", err)
		return nil, err
	}
	log.Println("Inserted URL into database successfully")

	conn, err := grpc.NewClient("analytics-service:50003", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Failed to connect to analytics service:", err)
		return nil, err
	}
	defer conn.Close()

	analyticsClient := analytics.NewAnalyticsServiceClient(conn)
	_, err = analyticsClient.AddAnalytics(ctx, &analytics.ShortURL{Shorturl: shortu_url})

	if err != nil {
		log.Println("Failed to add analytics for short URL:", err)
		return nil, err
	}
	log.Println("Added analytics for short URL successfully")

	return &url.ShortURL{
		Shorturl: shortu_url,
	}, nil

}

func (s *GRPCServer) GetLongURL(ctx context.Context, req *url.ShortURL) (*url.LongURL, error) {
	if req.Shorturl == "" {
		log.Println("Short URL is empty")
		return nil, errors.New("Short URL cannot be empty")
	}

	query := `SELECT long_url FROM urls WHERE short_url = $1`
	var long_url string
	resp, err := s.DBConn.ExecContext(ctx, query, req.Shorturl)
	if err != nil {
		log.Println("Failed to retrieve long URL from database:", err)
		return nil, err
	}
	if val, err := resp.RowsAffected(); val == 0 || err != nil {
		log.Println("No long URL found for the given short URL")
		return nil, err
	}
	err = s.DBConn.QueryRowContext(ctx, query, req.Shorturl).Scan(&long_url)
	if err != nil {
		log.Println("Failed to scan long URL:", err)
		return nil, err
	}
	log.Println("Retrieved long URL from database successfully")
	return &url.LongURL{
		Longurl: long_url,
	}, nil
}
