package api

import (
	"log/slog"
	"time"
)

// StartAtRiskMonitor launches a background goroutine that periodically checks
// students against at-risk thresholds and broadcasts SSE alerts when status changes.
// The monitor runs every interval and stops when done is closed.
func (h *Handler) StartAtRiskMonitor(interval time.Duration, done <-chan struct{}) {
	if h.sseBroker == nil {
		return
	}

	go func() {
		// Track previously known at-risk students to detect changes
		knownAtRisk := make(map[string]bool)

		// Initialize with current at-risk students
		for _, s := range h.store.ListStudents("", "", "", false) {
			if s.AtRisk {
				knownAtRisk[s.ID] = true
			}
		}

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				slog.Info("at-risk monitor stopped")
				return
			case <-ticker.C:
				h.checkAtRiskChanges(knownAtRisk)
			}
		}
	}()
}

// checkAtRiskChanges compares current at-risk status against known state
// and broadcasts SSE events for newly at-risk students.
func (h *Handler) checkAtRiskChanges(knownAtRisk map[string]bool) {
	students := h.store.ListStudents("", "", "", false)

	for _, s := range students {
		if s.AtRisk && !knownAtRisk[s.ID] {
			// Newly at-risk
			slog.Info("at-risk alert: student flagged", "student", s.ID, "name", s.LastName)

			h.sseBroker.Broadcast("xo", SSEEvent{
				Type: "at-risk-alert",
				Data: map[string]interface{}{
					"studentId":   s.ID,
					"studentName": s.Rank + " " + s.LastName,
					"riskFlags":   s.RiskFlags,
					"message":     s.Rank + " " + s.LastName + " has been flagged at-risk",
				},
			})
			h.sseBroker.Broadcast("staff", SSEEvent{
				Type: "at-risk-alert",
				Data: map[string]interface{}{
					"studentId":   s.ID,
					"studentName": s.Rank + " " + s.LastName,
					"riskFlags":   s.RiskFlags,
					"message":     s.Rank + " " + s.LastName + " has been flagged at-risk",
				},
			})
			// Also alert the SPC for that company
			if s.Company != "" {
				h.sseBroker.Broadcast("spc", SSEEvent{
					Type: "at-risk-alert",
					Data: map[string]interface{}{
						"studentId":   s.ID,
						"studentName": s.Rank + " " + s.LastName,
						"company":     s.Company,
						"riskFlags":   s.RiskFlags,
						"message":     s.Rank + " " + s.LastName + " has been flagged at-risk",
					},
				})
			}

			knownAtRisk[s.ID] = true
		} else if !s.AtRisk && knownAtRisk[s.ID] {
			// No longer at-risk
			delete(knownAtRisk, s.ID)
		}
	}
}
