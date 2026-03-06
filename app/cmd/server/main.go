package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"

	"heywood-tbs/internal/ai"
	"heywood-tbs/internal/api"
	"heywood-tbs/internal/data"
	"heywood-tbs/internal/middleware"
)

func main() {
	port := flag.String("port", "8080", "server port")
	dev := flag.Bool("dev", false, "development mode (CORS enabled, no embedded SPA)")
	dataDir := flag.String("data", "data", "path to JSON data directory")
	flag.Parse()

	// Load data
	store, err := data.NewStore(*dataDir)
	if err != nil {
		slog.Error("failed to load data", "error", err)
		os.Exit(1)
	}
	slog.Info("data loaded",
		"students", len(store.Students),
		"instructors", len(store.Instructors),
		"schedule", len(store.Schedule),
		"qualifications", len(store.Qualifications),
		"qualRecords", len(store.QualRecords),
		"feedback", len(store.Feedback),
	)

	// Initialize chat service (auto-detects Azure vs OpenAI from env vars)
	chatSvc := api.NewChatService()
	if chatSvc == nil {
		slog.Warn("No AI API keys configured — using mock responses")
	}

	// Initialize weather service
	weatherSvc := &ai.WeatherService{}

	// Build handler and router
	handler := api.NewHandler(store, chatSvc, weatherSvc, *dev)
	mux := api.SetupRouter(handler)

	// Serve static files in production (embedded SPA)
	if !*dev {
		// In production, serve the built React SPA from web/dist
		fs := http.FileServer(http.Dir("web/dist"))
		mux.Handle("GET /assets/", fs)
		mux.Handle("GET /favicon.ico", fs)
		// SPA fallback: serve index.html for all non-API, non-asset routes
		mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "web/dist/index.html")
		})
	}

	// Apply middleware
	chain := middleware.Chain(
		middleware.SecurityHeaders,
		middleware.CORS(*dev),
		middleware.Auth,
	)

	addr := ":" + *port
	slog.Info("Heywood TBS starting", "addr", addr, "dev", *dev)
	if err := http.ListenAndServe(addr, chain(mux)); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}
