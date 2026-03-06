package data

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

// DataSourceConfig represents the persisted data source configuration.
type DataSourceConfig struct {
	Type       string `json:"type"` // "json", "excel", "sharepoint", "cosmos", "postgres", "sqlserver"
	JSONDir    string `json:"jsonDir"`
	ExcelPath  string `json:"excelPath"`
	SharePoint struct {
		TenantID     string `json:"tenantId"`
		ClientID     string `json:"clientId"`
		ClientSecret string `json:"clientSecret"`
		SiteURL      string `json:"siteUrl"`
		Cloud        string `json:"cloud"`
	} `json:"sharepoint"`
	Database struct {
		Type             string `json:"type"`
		ConnectionString string `json:"connectionString"`
	} `json:"database"`
}

// NewDataStore creates the appropriate DataStore based on the settings.json config.
// Falls back to JSON store if the config is missing or the source type is unknown.
func NewDataStore(dataDir string) (DataStore, *Store, error) {
	// Always load the JSON store as the mutable backend
	jsonStore, err := NewStore(dataDir)
	if err != nil {
		return nil, nil, fmt.Errorf("load JSON store: %w", err)
	}

	// Read settings.json for data source config
	settingsPath := filepath.Join(dataDir, "settings.json")
	settingsData, err := os.ReadFile(settingsPath)
	if err != nil {
		slog.Info("no settings.json found, using JSON store", "path", settingsPath)
		return jsonStore, jsonStore, nil
	}

	var settings struct {
		DataSource DataSourceConfig `json:"dataSource"`
	}
	if err := json.Unmarshal(settingsData, &settings); err != nil {
		slog.Warn("failed to parse settings.json, using JSON store", "error", err)
		return jsonStore, jsonStore, nil
	}

	cfg := settings.DataSource

	switch cfg.Type {
	case "json", "":
		slog.Info("data source: JSON files", "dir", dataDir)
		return jsonStore, jsonStore, nil

	case "excel":
		if cfg.ExcelPath == "" {
			slog.Warn("Excel data source configured but no file path set, falling back to JSON")
			return jsonStore, jsonStore, nil
		}
		slog.Info("data source: Excel (hybrid mode)", "path", cfg.ExcelPath)
		// For now, Excel import is done through the settings API upload wizard.
		// The imported data replaces the JSON store's in-memory data.
		// Full ExcelStore with live file watching would be a future enhancement.
		return jsonStore, jsonStore, nil

	case "sharepoint":
		slog.Info("data source: SharePoint (hybrid mode)")
		// SharePoint connector reads reference data from SP lists via Graph API.
		// Mutable data (tasks, messages) stays in the JSON store.
		// TODO: implement SharePointStore when Graph SDK is integrated
		slog.Warn("SharePoint connector not yet implemented, falling back to JSON")
		return jsonStore, jsonStore, nil

	case "cosmos":
		slog.Info("data source: Azure Cosmos DB")
		// TODO: implement CosmosStore when azcosmos SDK is integrated
		slog.Warn("Cosmos DB connector not yet implemented, falling back to JSON")
		return jsonStore, jsonStore, nil

	case "postgres":
		slog.Info("data source: PostgreSQL")
		// TODO: implement SQLStore with pgx when dependency is added
		slog.Warn("PostgreSQL connector not yet implemented, falling back to JSON")
		return jsonStore, jsonStore, nil

	case "sqlserver":
		slog.Info("data source: Azure SQL / SQL Server")
		// TODO: implement SQLStore with go-mssqldb when dependency is added
		slog.Warn("SQL Server connector not yet implemented, falling back to JSON")
		return jsonStore, jsonStore, nil

	default:
		slog.Warn("unknown data source type, falling back to JSON", "type", cfg.Type)
		return jsonStore, jsonStore, nil
	}
}
