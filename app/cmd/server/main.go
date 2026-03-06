package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"heywood-tbs/internal/ai"
	"heywood-tbs/internal/api"
	"heywood-tbs/internal/auth"
	"heywood-tbs/internal/calendar"
	"heywood-tbs/internal/data"
	"heywood-tbs/internal/middleware"
	"heywood-tbs/internal/msgraph"
)

func main() {
	port := flag.String("port", "8080", "server port")
	dev := flag.Bool("dev", false, "development mode (CORS enabled, no embedded SPA)")
	dataDir := flag.String("data", "data", "path to JSON data directory")
	flag.Parse()

	// Load data via connector factory (reads settings.json to determine source)
	store, jsonStore, err := data.NewDataStore(*dataDir)
	if err != nil {
		slog.Error("failed to load data", "error", err)
		os.Exit(1)
	}
	_ = jsonStore // mutable store ref available for hybrid mode
	slog.Info("data loaded", "studentCount", store.TotalStudentCount())

	// Initialize chat service (auto-detects Azure vs OpenAI from env vars)
	chatSvc := api.NewChatService()
	if chatSvc == nil {
		slog.Warn("No AI API keys configured — using mock responses")
	}

	// Initialize live data services
	weatherSvc := &ai.WeatherService{}
	newsSvc := &ai.NewsService{}
	trafficSvc := &ai.TrafficService{}

	// Select auth provider based on AUTH_MODE env var
	var authProvider auth.IdentityProvider
	authMode := os.Getenv("AUTH_MODE")
	switch authMode {
	case "cac":
		rosterPath := filepath.Join(*dataDir, "user-roster.json")
		authProvider = auth.NewCACProvider(rosterPath)
		slog.Info("auth mode: CAC/PKI")
	default:
		authProvider = &auth.DemoProvider{}
		slog.Info("auth mode: Demo (role picker)")
	}

	// Initialize settings
	api.InitSettings(*dataDir)

	// Initialize Microsoft Graph client (for Outlook, SharePoint, Teams)
	graphTenantID := os.Getenv("GRAPH_TENANT_ID")
	graphClientID := os.Getenv("GRAPH_CLIENT_ID")
	graphClientSecret := os.Getenv("GRAPH_CLIENT_SECRET")
	graphCloud := os.Getenv("GRAPH_CLOUD") // "commercial", "gcc-high", "dod"
	if graphCloud == "" {
		graphCloud = "commercial"
	}

	if graphTenantID != "" && graphClientID != "" && graphClientSecret != "" {
		graphClient := msgraph.NewClient(graphTenantID, graphClientID, graphClientSecret, graphCloud)
		masterCalID := os.Getenv("GRAPH_MASTER_CALENDAR_ID")

		// Wire up real Outlook calendar
		outlookCal := calendar.NewOutlookCalendar(graphClient, masterCalID, nil)
		api.InitCalendar(outlookCal)

		// Wire up SharePoint and Teams services
		api.InitGraph(graphClient)

		slog.Info("Microsoft Graph connected", "cloud", graphCloud, "tenantID", graphTenantID[:8]+"...")
	} else {
		slog.Info("Microsoft Graph not configured — using mock calendar/mail")
	}

	// Build handler and router
	handler := api.NewHandler(store, chatSvc, weatherSvc, newsSvc, trafficSvc, authProvider, *dev)
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
		middleware.AuthWithProvider(authProvider),
	)

	addr := ":" + *port
	slog.Info("Heywood TBS starting", "addr", addr, "dev", *dev)
	if err := http.ListenAndServe(addr, chain(mux)); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}
