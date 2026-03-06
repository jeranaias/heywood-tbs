package api

import (
	"encoding/json"
	"net/http"

	"heywood-tbs/internal/middleware"
	"heywood-tbs/internal/msgraph"
)

var (
	graphClient    *msgraph.Client
	sharePointSvc  *msgraph.SharePointService
	teamsSvc       *msgraph.TeamsService
)

// InitGraph sets up the shared Microsoft Graph services.
func InitGraph(client *msgraph.Client) {
	graphClient = client
	sharePointSvc = msgraph.NewSharePointService(client)
	teamsSvc = msgraph.NewTeamsService(client)
}

// handleGraphTest tests the Microsoft Graph connection.
func (h *Handler) handleGraphTest(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	if role != "xo" && role != "staff" {
		writeError(w, 403, "admin only")
		return
	}

	if graphClient == nil {
		writeJSON(w, 200, map[string]interface{}{
			"status":  "not_configured",
			"message": "Microsoft Graph credentials not set. Set GRAPH_TENANT_ID, GRAPH_CLIENT_ID, and GRAPH_CLIENT_SECRET environment variables.",
		})
		return
	}

	if err := graphClient.TestConnection(); err != nil {
		writeJSON(w, 200, map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	writeJSON(w, 200, map[string]interface{}{
		"status":  "ok",
		"message": "Microsoft Graph connection successful",
	})
}

// --- SharePoint endpoints ---

// handleSharePointSite resolves a SharePoint site by URL.
func (h *Handler) handleSharePointSite(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	if role != "xo" && role != "staff" {
		writeError(w, 403, "admin only")
		return
	}

	if sharePointSvc == nil {
		writeError(w, 503, "SharePoint not configured")
		return
	}

	siteURL := r.URL.Query().Get("url")
	if siteURL == "" {
		writeError(w, 400, "url parameter required")
		return
	}

	site, err := sharePointSvc.GetSiteByURL(siteURL)
	if err != nil {
		writeJSON(w, 200, map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	writeJSON(w, 200, site)
}

// handleSharePointLists returns lists from a SharePoint site.
func (h *Handler) handleSharePointLists(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	if role != "xo" && role != "staff" {
		writeError(w, 403, "admin only")
		return
	}

	if sharePointSvc == nil {
		writeError(w, 503, "SharePoint not configured")
		return
	}

	siteID := r.URL.Query().Get("siteId")
	if siteID == "" {
		writeError(w, 400, "siteId parameter required")
		return
	}

	lists, err := sharePointSvc.ListLists(siteID)
	if err != nil {
		writeJSON(w, 200, map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	writeJSON(w, 200, map[string]interface{}{"lists": lists})
}

// handleSharePointListItems returns items from a SharePoint list.
func (h *Handler) handleSharePointListItems(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	if role != "xo" && role != "staff" {
		writeError(w, 403, "admin only")
		return
	}

	if sharePointSvc == nil {
		writeError(w, 503, "SharePoint not configured")
		return
	}

	siteID := r.URL.Query().Get("siteId")
	listID := r.URL.Query().Get("listId")
	if siteID == "" || listID == "" {
		writeError(w, 400, "siteId and listId parameters required")
		return
	}

	items, err := sharePointSvc.GetListItems(siteID, listID)
	if err != nil {
		writeJSON(w, 200, map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	writeJSON(w, 200, map[string]interface{}{
		"items": items,
		"count": len(items),
	})
}

// handleSharePointDrives returns document libraries for a SharePoint site.
func (h *Handler) handleSharePointDrives(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	if role != "xo" && role != "staff" {
		writeError(w, 403, "admin only")
		return
	}

	if sharePointSvc == nil {
		writeError(w, 503, "SharePoint not configured")
		return
	}

	siteID := r.URL.Query().Get("siteId")
	if siteID == "" {
		writeError(w, 400, "siteId parameter required")
		return
	}

	drives, err := sharePointSvc.GetDrives(siteID)
	if err != nil {
		writeJSON(w, 200, map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	writeJSON(w, 200, map[string]interface{}{"drives": drives})
}

// handleSharePointFiles returns files from a SharePoint document library.
func (h *Handler) handleSharePointFiles(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	if role != "xo" && role != "staff" {
		writeError(w, 403, "admin only")
		return
	}

	if sharePointSvc == nil {
		writeError(w, 503, "SharePoint not configured")
		return
	}

	siteID := r.URL.Query().Get("siteId")
	driveID := r.URL.Query().Get("driveId")
	folder := r.URL.Query().Get("folder")
	if siteID == "" || driveID == "" {
		writeError(w, 400, "siteId and driveId parameters required")
		return
	}

	items, err := sharePointSvc.ListDriveItems(siteID, driveID, folder)
	if err != nil {
		writeJSON(w, 200, map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	writeJSON(w, 200, map[string]interface{}{
		"files": items,
		"count": len(items),
	})
}

// --- Teams endpoints ---

// handleTeamsList returns teams the app has access to.
func (h *Handler) handleTeamsList(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	if role != "xo" && role != "staff" {
		writeError(w, 403, "admin only")
		return
	}

	if teamsSvc == nil {
		writeError(w, 503, "Teams not configured")
		return
	}

	teams, err := teamsSvc.ListTeams()
	if err != nil {
		writeJSON(w, 200, map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	writeJSON(w, 200, map[string]interface{}{"teams": teams})
}

// handleTeamsChannels returns channels for a team.
func (h *Handler) handleTeamsChannels(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	if role != "xo" && role != "staff" {
		writeError(w, 403, "admin only")
		return
	}

	if teamsSvc == nil {
		writeError(w, 503, "Teams not configured")
		return
	}

	teamID := r.URL.Query().Get("teamId")
	if teamID == "" {
		writeError(w, 400, "teamId parameter required")
		return
	}

	channels, err := teamsSvc.ListChannels(teamID)
	if err != nil {
		writeJSON(w, 200, map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	writeJSON(w, 200, map[string]interface{}{"channels": channels})
}

// handleTeamsFiles returns files from a Teams channel's shared folder.
func (h *Handler) handleTeamsFiles(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	if role != "xo" && role != "staff" {
		writeError(w, 403, "admin only")
		return
	}

	if teamsSvc == nil {
		writeError(w, 503, "Teams not configured")
		return
	}

	teamID := r.URL.Query().Get("teamId")
	channelID := r.URL.Query().Get("channelId")
	if teamID == "" || channelID == "" {
		writeError(w, 400, "teamId and channelId parameters required")
		return
	}

	var req struct {
		TeamID    string `json:"teamId"`
		ChannelID string `json:"channelId"`
	}
	if r.Method == "POST" {
		json.NewDecoder(r.Body).Decode(&req)
		if req.TeamID != "" {
			teamID = req.TeamID
		}
		if req.ChannelID != "" {
			channelID = req.ChannelID
		}
	}

	files, err := teamsSvc.ListChannelFiles(teamID, channelID)
	if err != nil {
		writeJSON(w, 200, map[string]interface{}{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	writeJSON(w, 200, map[string]interface{}{
		"files": files,
		"count": len(files),
	})
}
