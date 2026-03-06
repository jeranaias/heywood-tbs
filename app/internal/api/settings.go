package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"heywood-tbs/internal/data"
	"heywood-tbs/internal/middleware"
)

// AppSettings represents the full Heywood configuration persisted to settings.json.
type AppSettings struct {
	DataSource DataSourceSettings `json:"dataSource"`
	AI         AISettings         `json:"ai"`
	Outlook    OutlookSettings    `json:"outlook"`
	Auth       AuthSettings       `json:"auth"`
}

type DataSourceSettings struct {
	Type       string              `json:"type"` // "json", "excel", "sharepoint", "cosmos", "postgres", "sqlserver"
	JSONDir    string              `json:"jsonDir"`
	ExcelPath  string              `json:"excelPath"`
	SharePoint SharePointSettings  `json:"sharepoint"`
	Database   DatabaseSettings    `json:"database"`
}

type SharePointSettings struct {
	TenantID     string `json:"tenantId"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	SiteURL      string `json:"siteUrl"`
	Cloud        string `json:"cloud"` // "commercial", "gcc-high", "dod"
}

type DatabaseSettings struct {
	Type             string `json:"type"` // "cosmos", "postgres", "sqlserver"
	ConnectionString string `json:"connectionString"`
}

type AISettings struct {
	Provider   string `json:"provider"` // "openai", "azure-openai"
	Model      string `json:"model"`
	SearxngURL string `json:"searxngUrl"`
}

type OutlookSettings struct {
	Enabled             bool   `json:"enabled"`
	TenantID            string `json:"tenantId"`
	ClientID            string `json:"clientId"`
	ClientSecret        string `json:"clientSecret"`
	Cloud               string `json:"cloud"`
	MasterCalendarID    string `json:"masterCalendarId"`
	SyncIntervalMinutes int    `json:"syncIntervalMinutes"`
}

type AuthSettings struct {
	Mode string `json:"mode"` // "demo", "cac"
}

var (
	settingsMu   sync.RWMutex
	settingsPath string
)

func InitSettings(dataDir string) {
	settingsPath = filepath.Join(dataDir, "settings.json")
}

func loadSettings() (*AppSettings, error) {
	settingsMu.RLock()
	defer settingsMu.RUnlock()

	data, err := os.ReadFile(settingsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return defaultSettings(), nil
		}
		return nil, err
	}

	var s AppSettings
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func saveSettings(s *AppSettings) error {
	settingsMu.Lock()
	defer settingsMu.Unlock()

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(settingsPath, data, 0644)
}

func defaultSettings() *AppSettings {
	return &AppSettings{
		DataSource: DataSourceSettings{
			Type:    "json",
			JSONDir: "data",
		},
		AI: AISettings{
			Provider:   "openai",
			Model:      "gpt-4o",
			SearxngURL: "http://localhost:8888",
		},
		Outlook: OutlookSettings{
			SyncIntervalMinutes: 5,
			Cloud:               "commercial",
		},
		Auth: AuthSettings{
			Mode: "demo",
		},
	}
}

// handleGetSettings returns the current settings (secrets masked).
func (h *Handler) handleGetSettings(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	if role != "xo" && role != "staff" {
		writeError(w, 403, "settings are only accessible to XO and staff")
		return
	}

	s, err := loadSettings()
	if err != nil {
		writeError(w, 500, "failed to load settings")
		return
	}

	// Mask secrets before sending to client
	resp := maskSettings(s)
	writeJSON(w, 200, resp)
}

// handleUpdateSettings updates the persisted settings.
func (h *Handler) handleUpdateSettings(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	if role != "xo" && role != "staff" {
		writeError(w, 403, "settings are only accessible to XO and staff")
		return
	}

	var incoming AppSettings
	if err := json.NewDecoder(r.Body).Decode(&incoming); err != nil {
		writeError(w, 400, "invalid settings JSON")
		return
	}

	// Load current settings to preserve masked secrets
	current, err := loadSettings()
	if err != nil {
		writeError(w, 500, "failed to load current settings")
		return
	}

	// If client sent masked secrets, keep the existing values
	if incoming.DataSource.SharePoint.ClientSecret == "••••••••" || incoming.DataSource.SharePoint.ClientSecret == "" {
		incoming.DataSource.SharePoint.ClientSecret = current.DataSource.SharePoint.ClientSecret
	}
	if incoming.DataSource.Database.ConnectionString == "••••••••" || incoming.DataSource.Database.ConnectionString == "" {
		incoming.DataSource.Database.ConnectionString = current.DataSource.Database.ConnectionString
	}
	if incoming.Outlook.ClientSecret == "••••••••" || incoming.Outlook.ClientSecret == "" {
		incoming.Outlook.ClientSecret = current.Outlook.ClientSecret
	}

	if err := saveSettings(&incoming); err != nil {
		writeError(w, 500, "failed to save settings")
		return
	}

	writeJSON(w, 200, map[string]string{"status": "saved", "note": "changes take effect on next restart"})
}

// handleTestConnection tests connectivity for a data source.
func (h *Handler) handleTestConnection(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	if role != "xo" && role != "staff" {
		writeError(w, 403, "settings are only accessible to XO and staff")
		return
	}

	var req struct {
		Type             string `json:"type"`
		ConnectionString string `json:"connectionString,omitempty"`
		TenantID         string `json:"tenantId,omitempty"`
		ClientID         string `json:"clientId,omitempty"`
		ClientSecret     string `json:"clientSecret,omitempty"`
		SiteURL          string `json:"siteUrl,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, 400, "invalid request body")
		return
	}

	// For now, return a simulated result based on type
	switch req.Type {
	case "json":
		writeJSON(w, 200, map[string]interface{}{"status": "ok", "message": "JSON data directory is accessible"})
	case "sharepoint":
		if req.TenantID == "" || req.ClientID == "" {
			writeJSON(w, 200, map[string]interface{}{"status": "error", "message": "Tenant ID and Client ID are required"})
			return
		}
		writeJSON(w, 200, map[string]interface{}{"status": "pending", "message": "SharePoint connector not yet configured — will be available after connector setup"})
	case "cosmos", "postgres", "sqlserver":
		if req.ConnectionString == "" {
			writeJSON(w, 200, map[string]interface{}{"status": "error", "message": "Connection string is required"})
			return
		}
		writeJSON(w, 200, map[string]interface{}{"status": "pending", "message": "Database connector not yet configured — will be available after connector setup"})
	default:
		writeJSON(w, 200, map[string]interface{}{"status": "error", "message": "Unknown data source type"})
	}
}

// handleUpload handles Excel/CSV file uploads.
func (h *Handler) handleUpload(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	if role != "xo" && role != "staff" {
		writeError(w, 403, "uploads are only accessible to XO and staff")
		return
	}

	r.ParseMultipartForm(32 << 20) // 32MB max
	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, 400, "no file provided")
		return
	}
	defer file.Close()

	// Validate file type
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".xlsx" && ext != ".csv" {
		writeError(w, 400, "only .xlsx and .csv files are supported")
		return
	}

	// Save to uploads directory
	uploadDir := filepath.Join(filepath.Dir(settingsPath), "uploads")
	os.MkdirAll(uploadDir, 0755)
	destPath := filepath.Join(uploadDir, header.Filename)

	dest, err := os.Create(destPath)
	if err != nil {
		writeError(w, 500, "failed to save uploaded file")
		return
	}
	defer dest.Close()

	if _, err := io.Copy(dest, file); err != nil {
		writeError(w, 500, "failed to write file")
		return
	}

	writeJSON(w, 200, map[string]interface{}{
		"status":   "uploaded",
		"filename": header.Filename,
		"path":     destPath,
		"size":     header.Size,
		"type":     ext[1:], // "xlsx" or "csv"
	})
}

// handleColumnMap auto-detects column mappings for an uploaded Excel file.
func (h *Handler) handleColumnMap(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	if role != "xo" && role != "staff" {
		writeError(w, 403, "settings are only accessible to XO and staff")
		return
	}

	var req struct {
		FilePath string `json:"filePath"`
		Sheet    string `json:"sheet"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, 400, "invalid request body")
		return
	}

	reader := data.NewExcelReader(req.FilePath)

	// List sheets if no sheet specified
	if req.Sheet == "" {
		sheets, err := reader.ListSheets()
		if err != nil {
			writeError(w, 400, err.Error())
			return
		}
		writeJSON(w, 200, map[string]interface{}{"sheets": sheets})
		return
	}

	// Read headers and auto-map
	headers, err := reader.ReadHeaders(req.Sheet)
	if err != nil {
		writeError(w, 400, err.Error())
		return
	}

	mappings := data.AutoMapColumns(headers)
	autoMapped := 0
	for _, m := range mappings {
		if m.AutoMatch {
			autoMapped++
		}
	}

	writeJSON(w, 200, map[string]interface{}{
		"headers":       headers,
		"mappings":      mappings,
		"autoMapped":    autoMapped,
		"totalColumns":  len(headers),
		"availableFields": data.AvailableFields(),
	})
}

// handleUploadPreview returns a preview of parsed data before import.
func (h *Handler) handleUploadPreview(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	if role != "xo" && role != "staff" {
		writeError(w, 403, "settings are only accessible to XO and staff")
		return
	}

	var req struct {
		FilePath string                `json:"filePath"`
		Sheet    string                `json:"sheet"`
		Mappings []data.ColumnMapping  `json:"mappings"`
		DataType string                `json:"dataType"` // "students" or "instructors"
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, 400, "invalid request body")
		return
	}

	reader := data.NewExcelReader(req.FilePath)

	switch req.DataType {
	case "students":
		students, err := reader.ReadStudents(req.Sheet, req.Mappings)
		if err != nil {
			writeError(w, 400, err.Error())
			return
		}
		// Return first 10 for preview
		preview := students
		if len(preview) > 10 {
			preview = preview[:10]
		}
		writeJSON(w, 200, map[string]interface{}{
			"dataType":    "students",
			"totalRows":   len(students),
			"previewRows": len(preview),
			"data":        preview,
		})
	case "instructors":
		instructors, err := reader.ReadInstructors(req.Sheet, req.Mappings)
		if err != nil {
			writeError(w, 400, err.Error())
			return
		}
		preview := instructors
		if len(preview) > 10 {
			preview = preview[:10]
		}
		writeJSON(w, 200, map[string]interface{}{
			"dataType":    "instructors",
			"totalRows":   len(instructors),
			"previewRows": len(preview),
			"data":        preview,
		})
	default:
		writeError(w, 400, "dataType must be 'students' or 'instructors'")
	}
}

// handleSystemInfo returns system status information.
func (h *Handler) handleSystemInfo(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	if role != "xo" && role != "staff" {
		writeError(w, 403, "system info is only accessible to XO and staff")
		return
	}

	s, err := loadSettings()
	if err != nil {
		s = defaultSettings()
	}

	// Detect AI configuration
	aiStatus := "not configured"
	aiKeyHint := ""
	if key := os.Getenv("OPENAI_API_KEY"); key != "" {
		aiStatus = "active"
		if len(key) > 4 {
			aiKeyHint = "••••" + key[len(key)-4:]
		}
	}
	if key := os.Getenv("AZURE_OPENAI_KEY"); key != "" {
		aiStatus = "active (Azure)"
		if len(key) > 4 {
			aiKeyHint = "••••" + key[len(key)-4:]
		}
	}

	info := map[string]interface{}{
		"version":      "1.0.0",
		"authMode":     s.Auth.Mode,
		"dataSource":   s.DataSource.Type,
		"studentCount": h.store.TotalStudentCount(),
		"ai": map[string]interface{}{
			"status":  aiStatus,
			"keyHint": aiKeyHint,
			"model":   s.AI.Model,
		},
		"searxng": map[string]interface{}{
			"url":    s.AI.SearxngURL,
			"status": "configured",
		},
		"outlook": map[string]interface{}{
			"enabled": s.Outlook.Enabled,
		},
	}

	writeJSON(w, 200, info)
}

// maskSettings returns a copy with secrets masked.
func maskSettings(s *AppSettings) *AppSettings {
	copy := *s
	if copy.DataSource.SharePoint.ClientSecret != "" {
		copy.DataSource.SharePoint.ClientSecret = "••••••••"
	}
	if copy.DataSource.Database.ConnectionString != "" {
		cs := copy.DataSource.Database.ConnectionString
		if len(cs) > 8 {
			copy.DataSource.Database.ConnectionString = cs[:4] + "••••••••" + cs[len(cs)-4:]
		} else {
			copy.DataSource.Database.ConnectionString = "••••••••"
		}
	}
	if copy.Outlook.ClientSecret != "" {
		copy.Outlook.ClientSecret = "••••••••"
	}
	return &copy
}
