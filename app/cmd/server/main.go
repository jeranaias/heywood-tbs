package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"heywood-tbs/internal/ai"
	"heywood-tbs/internal/api"
	"heywood-tbs/internal/auth"
	"heywood-tbs/internal/calendar"
	"heywood-tbs/internal/config"
	"heywood-tbs/internal/data"
	"heywood-tbs/internal/middleware"
	"heywood-tbs/internal/msgraph"
)

func main() {
	cfg := config.Load()

	// Load data via connector factory (reads settings.json to determine source)
	store, jsonStore, err := data.NewDataStore(cfg.DataDir)
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

	// Select auth provider
	var authProvider auth.IdentityProvider
	switch cfg.AuthMode {
	case "cac":
		rosterPath := filepath.Join(cfg.DataDir, "user-roster.json")
		authProvider = auth.NewCACProvider(rosterPath)
		slog.Info("auth mode: CAC/PKI")
	default:
		authProvider = &auth.DemoProvider{}
		slog.Info("auth mode: Demo (role picker)")
	}

	// Initialize Microsoft Graph client (for Outlook, SharePoint, Teams)
	var calProvider calendar.CalendarProvider
	var graphClient *msgraph.Client
	var sharePointSvc *msgraph.SharePointService
	var teamsSvc *msgraph.TeamsService

	if cfg.GraphTenantID != "" && cfg.GraphClientID != "" && cfg.GraphClientSecret != "" {
		graphClient = msgraph.NewClient(cfg.GraphTenantID, cfg.GraphClientID, cfg.GraphClientSecret, cfg.GraphCloud)

		// Wire up real Outlook calendar
		calProvider = calendar.NewOutlookCalendar(graphClient, cfg.GraphMasterCalID, nil)

		// Wire up SharePoint and Teams services
		sharePointSvc = msgraph.NewSharePointService(graphClient)
		teamsSvc = msgraph.NewTeamsService(graphClient)

		slog.Info("Microsoft Graph connected", "cloud", cfg.GraphCloud, "tenantID", cfg.GraphTenantID[:8]+"...")
	} else {
		slog.Info("Microsoft Graph not configured — using mock calendar/mail")
	}

	// Settings file path
	settingsPath := filepath.Join(cfg.DataDir, "settings.json")

	// Build handler and router
	handler := api.NewHandler(store, chatSvc, weatherSvc, newsSvc, trafficSvc, authProvider, cfg.Dev,
		calProvider, graphClient, sharePointSvc, teamsSvc, settingsPath)
	mux := api.SetupRouter(handler)

	// Serve static files in production (embedded SPA)
	if !cfg.Dev {
		// In production, serve the built React SPA from web/dist
		fs := http.FileServer(http.Dir("web/dist"))
		mux.Handle("GET /assets/", fs)
		mux.Handle("GET /favicon.ico", fs)
		// SPA fallback: serve index.html for all non-API, non-asset routes
		mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "web/dist/index.html")
		})
	}

	// Apply middleware (outermost first)
	apiLimiter := middleware.NewRateLimiter(60, 120) // 60 req/s per IP, burst 120
	chain := middleware.Chain(
		middleware.Recovery,
		middleware.MaxBodySize(1<<20), // 1MB default body limit
		apiLimiter.Middleware,
		middleware.SecurityHeaders,
		middleware.CORS(cfg.Dev),
		middleware.AuthWithProvider(authProvider),
	)

	addr := ":" + cfg.Port
	srv := &http.Server{
		Addr:         addr,
		Handler:      chain(mux),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 120 * time.Second, // SSE streams need long writes
		IdleTimeout:  120 * time.Second,
	}

	// Graceful shutdown on SIGTERM/SIGINT
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("Heywood TBS starting", "addr", addr, "dev", cfg.Dev)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down gracefully (15s timeout)...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("forced shutdown", "error", err)
	}
	slog.Info("server stopped")
}
