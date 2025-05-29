package internal

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Naman-B-Parlecha/url-shortener/broker-service/analytics"
	"github.com/Naman-B-Parlecha/url-shortener/broker-service/redirect"
	"github.com/Naman-B-Parlecha/url-shortener/broker-service/url"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RequestPayload struct {
	Action      string       `json:"action"`
	RedirectURL *RedirectURL `json:"redirect_url,omitempty"`
	ShortURL    *ShortURL    `json:"short_url,omitempty"`
	Analyatics  *Analytics   `json:"analytics,omitempty"`
}

type RedirectURL struct {
	ShortURL string `json:"short_url"`
}

type ShortURL struct {
	LongURL string `json:"long_url"`
}

type Analytics struct {
	ShortURL string `json:"short_url"`
}

func Handler(ctx *gin.Context) {
	var payload RequestPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	switch payload.Action {
	case "shorten":
		if payload.ShortURL == nil || payload.ShortURL.LongURL == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Short URL is required"})
			return
		}
		handleShorten(ctx, *payload.ShortURL)
	case "analytics":
		if payload.Analyatics == nil || payload.Analyatics.ShortURL == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Analytics data is required"})
			return
		}
		handleAnalytics(ctx, *payload.Analyatics)
	default:
		ctx.JSON(400, gin.H{"error": "Invalid action"})
	}
}

func Redirect(ctx *gin.Context) {
	shortURL := ctx.Param("shortURL")
	if shortURL == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Short URL is required"})
		return
	}

	handleRedirect(ctx, RedirectURL{ShortURL: shortURL})
}

func handleRedirect(ctx *gin.Context, redirectURL RedirectURL) {
	conn, err := grpc.NewClient("localhost:50002", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to shortener service"})
		return
	}
	defer conn.Close()

	grpcCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := redirect.NewRedirectServiceClient(conn)
	resp, err := client.GetRedirectURL(grpcCtx, &redirect.ShortURL{Shorturl: redirectURL.ShortURL})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve long URL", "message": err.Error()})
		return
	}
	ctx.Redirect(http.StatusFound, resp.Longurl)
}

func handleShorten(ctx *gin.Context, shortURL ShortURL) {
	conn, err := grpc.NewClient("localhost:50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to shortener service"})
		return
	}
	defer conn.Close()

	grpcCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client := url.NewURLServiceClient(conn)
	resp, err := client.GenerateShortURL(grpcCtx, &url.LongURL{Longurl: shortURL.LongURL})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate short URL", "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"short_url": fmt.Sprintf("http://localhost:8081/%s", resp.Shorturl)})
}

func handleAnalytics(ctx *gin.Context, a Analytics) {
	conn, err := grpc.NewClient("localhost:50003", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to analytics service"})
		return
	}
	defer conn.Close()

	grpcCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := analytics.NewAnalyticsServiceClient(conn)
	shorturl_id := strings.Split(a.ShortURL, "/")
	short_id := shorturl_id[len(shorturl_id)-1]
	resp, err := client.GetAnalytics(grpcCtx, &analytics.ShortURL{Shorturl: short_id})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record analytics", "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Analytics recorded successfully", "data": resp})
}
