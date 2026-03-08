package api

import (
	"net/http"
	"time"
)

var serverStartTime = time.Now()

// handleHealthz returns server health status for k8s/Docker probes.
// Bypasses auth middleware — suitable for unauthenticated liveness checks.
func (h *Handler) handleHealthz(w http.ResponseWriter, r *http.Request) {
	dataStatus := "ok"
	if h.store == nil {
		dataStatus = "unavailable"
	}

	aiStatus := "not_configured"
	if h.chatSvc != nil {
		aiStatus = "configured"
	}

	graphStatus := "not_configured"
	if h.graphClient != nil {
		graphStatus = "configured"
	}

	status := "ok"
	httpCode := 200
	if dataStatus != "ok" {
		status = "degraded"
		httpCode = 503
	}

	writeJSON(w, httpCode, map[string]interface{}{
		"status":    status,
		"version":   "1.0.0",
		"uptime":    time.Since(serverStartTime).String(),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"checks": map[string]string{
			"dataStore": dataStatus,
			"ai":        aiStatus,
			"graph":     graphStatus,
		},
	})
}
