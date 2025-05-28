package internal

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Naman-B-Parlecha/url-shortener/broker-service/url"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RequestPayload struct {
	Action      string       `json:"action"`
	RedirectURL *RedirectURL `json:"redirect_url,omitempty"`
	ShortURL    *ShortURL    `json:"short_url,omitempty"`
	Analyatics  *Analytics   `json:"analyatics,omitempty"`
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
	case "redirect":
		handleRedirect(ctx, *payload.RedirectURL)
	case "shorten":
		handleShorten(ctx, *payload.ShortURL)
	case "analytics":
		handleAnalytics(ctx, *payload.Analyatics)
	default:
		ctx.JSON(400, gin.H{"error": "Invalid action"})
	}
}

func handleRedirect(ctx *gin.Context, redirectURL RedirectURL) {
	// here i will make a client to call the redirect service
}

func handleShorten(ctx *gin.Context, shortURL ShortURL) {
	conn, err := grpc.NewClient("url-shortener-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	ctx.JSON(http.StatusOK, gin.H{"short_url": fmt.Sprintf("http://nig-url/%s", resp.Shorturl)})
}

func handleAnalytics(ctx *gin.Context, analytics Analytics) {}
