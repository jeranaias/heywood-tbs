package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	// Database drivers — imported for side-effect registration
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"
)

// DataSourceConfig represents the persisted data source configuration.
type DataSourceConfig struct {
	Type       string `json:"type"` // "json", "excel", "sqlite", "postgres", "sharepoint", "cosmos"
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
		Type             string `json:"type"`             // "sqlite", "postgres"
		ConnectionString string `json:"connectionString"` // DSN or file path
	} `json:"database"`
}

// NewDataStore creates the appropriate DataStore based on the settings.json config.
// Returns (active store, JSON store reference, error).
func NewDataStore(dataDir string) (DataStore, *Store, error) {
	// Always load the JSON store for reference data
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
		}
		return jsonStore, jsonStore, nil

	case "sqlite":
		dsn := cfg.Database.ConnectionString
		if dsn == "" {
			dsn = filepath.Join(dataDir, "heywood.db")
		}
		slog.Info("data source: SQLite", "path", dsn)
		sqlStore, err := NewSQLStore(dataDir, "sqlite", dsn)
		if err != nil {
			slog.Error("failed to open SQLite, falling back to JSON", "error", err)
			return jsonStore, jsonStore, nil
		}
		return sqlStore, jsonStore, nil

	case "postgres":
		dsn := cfg.Database.ConnectionString
		if dsn == "" {
			slog.Warn("PostgreSQL configured but no connection string set, falling back to JSON")
			return jsonStore, jsonStore, nil
		}
		slog.Info("data source: PostgreSQL")
		sqlStore, err := NewSQLStore(dataDir, "pgx", dsn)
		if err != nil {
			slog.Error("failed to connect to PostgreSQL, falling back to JSON", "error", err)
			return jsonStore, jsonStore, nil
		}
		return sqlStore, jsonStore, nil

	case "sharepoint":
		slog.Warn("SharePoint connector not yet implemented, falling back to JSON")
		return jsonStore, jsonStore, nil

	case "cosmos":
		slog.Warn("Cosmos DB connector not yet implemented, falling back to JSON")
		return jsonStore, jsonStore, nil

	default:
		slog.Warn("unknown data source type, falling back to JSON", "type", cfg.Type)
		return jsonStore, jsonStore, nil
	}
}

// TestConnection tests connectivity to a database.
func TestConnection(driver, dsn string) error {
	var driverName string
	switch driver {
	case "sqlite":
		driverName = "sqlite"
	case "postgres":
		driverName = "pgx"
	default:
		return fmt.Errorf("unsupported driver: %s", driver)
	}

	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}
	return nil
}
