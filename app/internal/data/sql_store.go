package data

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"heywood-tbs/internal/models"
)

// SQLStore implements DataStore backed by a SQL database (SQLite or PostgreSQL).
// Reference data (students, instructors, schedule, etc.) is loaded from JSON files
// into memory (same as Store). Mutable data (tasks, messages, notifications) and
// chat history are stored in the database with proper per-user isolation.
type SQLStore struct {
	*Store // embedded JSON store for reference data
	db     *sql.DB
	mu     sync.RWMutex
	nextID int
	driver string // "sqlite" or "pgx"
}

// NewSQLStore creates a SQL-backed store. Reference data is loaded from dataDir JSON files.
// Mutable data is stored in the SQL database at the given DSN.
// driver: "sqlite" for SQLite, "pgx" for PostgreSQL
func NewSQLStore(dataDir, driver, dsn string) (*SQLStore, error) {
	// Load reference data from JSON files
	jsonStore, err := NewStore(dataDir)
	if err != nil {
		return nil, fmt.Errorf("load reference data: %w", err)
	}

	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	s := &SQLStore{
		Store:  jsonStore,
		db:     db,
		nextID: 1,
		driver: driver,
	}

	if err := s.migrate(); err != nil {
		return nil, fmt.Errorf("migrate schema: %w", err)
	}

	// Load nextID from database
	s.loadNextID()

	slog.Info("SQL store initialized", "driver", driver)
	return s, nil
}

// migration represents a numbered schema migration.
type migration struct {
	version int
	stmts   []string
}

// migrations returns all schema migrations in order.
// Version 1 uses IF NOT EXISTS so it is safe to re-run on existing databases
// that were created before the schema_version table was introduced.
func (s *SQLStore) migrations() []migration {
	return []migration{
		{
			version: 1,
			stmts: []string{
				`CREATE TABLE IF NOT EXISTS tasks (
					id TEXT PRIMARY KEY,
					title TEXT NOT NULL,
					description TEXT DEFAULT '',
					assigned_to TEXT NOT NULL,
					created_by TEXT NOT NULL DEFAULT 'heywood',
					priority TEXT NOT NULL DEFAULT 'medium',
					status TEXT NOT NULL DEFAULT 'pending',
					due_date TEXT DEFAULT '',
					related_id TEXT DEFAULT '',
					created_at TEXT NOT NULL,
					updated_at TEXT NOT NULL
				)`,
				`CREATE TABLE IF NOT EXISTS messages (
					id TEXT PRIMARY KEY,
					from_user TEXT NOT NULL,
					to_user TEXT NOT NULL,
					subject TEXT NOT NULL,
					body TEXT NOT NULL,
					is_read INTEGER NOT NULL DEFAULT 0,
					related_id TEXT DEFAULT '',
					created_at TEXT NOT NULL
				)`,
				`CREATE TABLE IF NOT EXISTS notifications (
					id TEXT PRIMARY KEY,
					user_role TEXT NOT NULL,
					type TEXT NOT NULL,
					title TEXT NOT NULL,
					body TEXT DEFAULT '',
					is_read INTEGER NOT NULL DEFAULT 0,
					action_url TEXT DEFAULT '',
					created_at TEXT NOT NULL
				)`,
				`CREATE TABLE IF NOT EXISTS chat_sessions (
					id TEXT PRIMARY KEY,
					user_id TEXT NOT NULL,
					user_role TEXT NOT NULL,
					company TEXT DEFAULT '',
					title TEXT DEFAULT 'New conversation',
					created_at TEXT NOT NULL,
					updated_at TEXT NOT NULL
				)`,
				`CREATE TABLE IF NOT EXISTS chat_messages (
					id INTEGER PRIMARY KEY,
					session_id TEXT NOT NULL,
					role TEXT NOT NULL,
					content TEXT NOT NULL,
					created_at TEXT NOT NULL
				)`,
				`CREATE TABLE IF NOT EXISTS id_counter (
					key TEXT PRIMARY KEY,
					value INTEGER NOT NULL DEFAULT 0
				)`,
				// Indexes for common queries
				`CREATE INDEX IF NOT EXISTS idx_tasks_assigned ON tasks(assigned_to)`,
				`CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status)`,
				`CREATE INDEX IF NOT EXISTS idx_messages_to ON messages(to_user)`,
				`CREATE INDEX IF NOT EXISTS idx_notifications_role ON notifications(user_role)`,
				`CREATE INDEX IF NOT EXISTS idx_notifications_read ON notifications(is_read)`,
				`CREATE INDEX IF NOT EXISTS idx_chat_sessions_user ON chat_sessions(user_id)`,
				`CREATE INDEX IF NOT EXISTS idx_chat_messages_session ON chat_messages(session_id)`,
			},
		},
		// Future migrations go here as {version: 2, stmts: [...]}, etc.
	}
}

// migrate applies versioned schema migrations. It creates the schema_version
// tracking table if it does not exist, reads the current version, and applies
// only migrations newer than that version. Existing databases without a
// schema_version table are handled gracefully because version 1 migrations
// use IF NOT EXISTS.
func (s *SQLStore) migrate() error {
	// Ensure the schema_version table exists
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS schema_version (
		version INTEGER NOT NULL
	)`)
	if err != nil {
		return fmt.Errorf("create schema_version table: %w", err)
	}

	// Read current schema version
	currentVersion := 0
	row := s.db.QueryRow(`SELECT COALESCE(MAX(version), 0) FROM schema_version`)
	if err := row.Scan(&currentVersion); err != nil {
		// Table exists but might be empty; treat as version 0
		currentVersion = 0
	}

	slog.Info("schema migration check", "currentVersion", currentVersion)

	for _, m := range s.migrations() {
		if m.version <= currentVersion {
			continue
		}
		slog.Info("applying migration", "version", m.version)
		for _, stmt := range m.stmts {
			if _, err := s.db.Exec(stmt); err != nil {
				return fmt.Errorf("migration v%d exec %q: %w", m.version, stmt[:min(60, len(stmt))], err)
			}
		}
		// Record the newly applied version
		_, err := s.db.Exec(`INSERT INTO schema_version (version) VALUES ($1)`, m.version)
		if err != nil {
			return fmt.Errorf("record migration v%d: %w", m.version, err)
		}
	}

	// Initialize id counter if not present (safe for both SQLite and PostgreSQL)
	if s.driver == "pgx" {
		s.db.Exec(`INSERT INTO id_counter (key, value) VALUES ('next_id', 1) ON CONFLICT (key) DO NOTHING`)
	} else {
		s.db.Exec(`INSERT OR IGNORE INTO id_counter (key, value) VALUES ('next_id', 1)`)
	}

	return nil
}

func (s *SQLStore) loadNextID() {
	row := s.db.QueryRow(`SELECT value FROM id_counter WHERE key = 'next_id'`)
	if err := row.Scan(&s.nextID); err != nil {
		s.nextID = 1
	}
}

func (s *SQLStore) incrementID() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := s.nextID
	s.nextID++
	s.db.Exec(`UPDATE id_counter SET value = $1 WHERE key = 'next_id'`, s.nextID)
	return id
}

// Close closes the database connection.
func (s *SQLStore) Close() error {
	return s.db.Close()
}

// --- Task operations (override Store's in-memory versions) ---

func (s *SQLStore) CreateTask(task models.Task) error {
	if task.ID == "" {
		task.ID = fmt.Sprintf("TSK-%03d", s.incrementID())
	}
	now := time.Now().Format(time.RFC3339)
	if task.CreatedAt == "" {
		task.CreatedAt = now
	}
	task.UpdatedAt = now
	if task.Status == "" {
		task.Status = "pending"
	}

	_, err := s.db.Exec(
		`INSERT INTO tasks (id, title, description, assigned_to, created_by, priority, status, due_date, related_id, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		task.ID, task.Title, task.Description, task.AssignedTo, task.CreatedBy,
		task.Priority, task.Status, task.DueDate, task.RelatedID,
		task.CreatedAt, task.UpdatedAt,
	)
	return err
}

func (s *SQLStore) ListTasks(assignedTo string) []models.Task {
	var rows *sql.Rows
	var err error

	if assignedTo == "" {
		rows, err = s.db.Query(
			`SELECT id, title, description, assigned_to, created_by, priority, status, due_date, related_id, created_at, updated_at
			 FROM tasks ORDER BY created_at DESC`)
	} else {
		rows, err = s.db.Query(
			`SELECT id, title, description, assigned_to, created_by, priority, status, due_date, related_id, created_at, updated_at
			 FROM tasks WHERE LOWER(assigned_to) = LOWER($1) ORDER BY created_at DESC`, assignedTo)
	}
	if err != nil {
		slog.Error("list tasks", "error", err)
		return nil
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.AssignedTo, &t.CreatedBy,
			&t.Priority, &t.Status, &t.DueDate, &t.RelatedID, &t.CreatedAt, &t.UpdatedAt); err != nil {
			slog.Error("scan task", "error", err)
			continue
		}
		tasks = append(tasks, t)
	}
	return tasks
}

func (s *SQLStore) GetTask(id string) (*models.Task, bool) {
	var t models.Task
	err := s.db.QueryRow(
		`SELECT id, title, description, assigned_to, created_by, priority, status, due_date, related_id, created_at, updated_at
		 FROM tasks WHERE id = $1`, id,
	).Scan(&t.ID, &t.Title, &t.Description, &t.AssignedTo, &t.CreatedBy,
		&t.Priority, &t.Status, &t.DueDate, &t.RelatedID, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, false
	}
	return &t, true
}

func (s *SQLStore) UpdateTask(id string, req models.TaskUpdateRequest) error {
	var sets []string
	var args []interface{}
	argN := 1

	if req.Status != nil {
		sets = append(sets, fmt.Sprintf("status = $%d", argN))
		args = append(args, *req.Status)
		argN++
	}
	if req.Priority != nil {
		sets = append(sets, fmt.Sprintf("priority = $%d", argN))
		args = append(args, *req.Priority)
		argN++
	}
	if req.AssignedTo != nil {
		sets = append(sets, fmt.Sprintf("assigned_to = $%d", argN))
		args = append(args, *req.AssignedTo)
		argN++
	}

	if len(sets) == 0 {
		return nil
	}

	sets = append(sets, fmt.Sprintf("updated_at = $%d", argN))
	args = append(args, time.Now().Format(time.RFC3339))
	argN++

	args = append(args, id)
	query := fmt.Sprintf("UPDATE tasks SET %s WHERE id = $%d", strings.Join(sets, ", "), argN)

	result, err := s.db.Exec(query, args...)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return fmt.Errorf("task %s not found", id)
	}
	return nil
}

// --- Message operations ---

func (s *SQLStore) CreateMessage(msg models.Message) error {
	if msg.ID == "" {
		msg.ID = fmt.Sprintf("MSG-%03d", s.incrementID())
	}
	if msg.CreatedAt == "" {
		msg.CreatedAt = time.Now().Format(time.RFC3339)
	}

	_, err := s.db.Exec(
		`INSERT INTO messages (id, from_user, to_user, subject, body, is_read, related_id, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		msg.ID, msg.From, msg.To, msg.Subject, msg.Body,
		boolToInt(msg.Read), msg.RelatedID, msg.CreatedAt,
	)
	return err
}

func (s *SQLStore) ListMessages(userRole string) []models.Message {
	var rows *sql.Rows
	var err error

	if userRole == "" {
		rows, err = s.db.Query(
			`SELECT id, from_user, to_user, subject, body, is_read, related_id, created_at
			 FROM messages ORDER BY created_at DESC`)
	} else {
		rows, err = s.db.Query(
			`SELECT id, from_user, to_user, subject, body, is_read, related_id, created_at
			 FROM messages WHERE LOWER(to_user) = LOWER($1) ORDER BY created_at DESC`, userRole)
	}
	if err != nil {
		slog.Error("list messages", "error", err)
		return nil
	}
	defer rows.Close()

	var msgs []models.Message
	for rows.Next() {
		var m models.Message
		var isRead int
		if err := rows.Scan(&m.ID, &m.From, &m.To, &m.Subject, &m.Body,
			&isRead, &m.RelatedID, &m.CreatedAt); err != nil {
			slog.Error("scan message", "error", err)
			continue
		}
		m.Read = isRead != 0
		msgs = append(msgs, m)
	}
	return msgs
}

func (s *SQLStore) MarkMessageRead(id string) error {
	result, err := s.db.Exec(`UPDATE messages SET is_read = 1 WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return fmt.Errorf("message %s not found", id)
	}
	return nil
}

// --- Notification operations ---

func (s *SQLStore) CreateNotification(n models.Notification) error {
	if n.ID == "" {
		n.ID = fmt.Sprintf("NTF-%03d", s.incrementID())
	}
	if n.CreatedAt == "" {
		n.CreatedAt = time.Now().Format(time.RFC3339)
	}

	_, err := s.db.Exec(
		`INSERT INTO notifications (id, user_role, type, title, body, is_read, action_url, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		n.ID, n.UserRole, n.Type, n.Title, n.Body,
		boolToInt(n.Read), n.ActionURL, n.CreatedAt,
	)
	return err
}

func (s *SQLStore) ListNotifications(userRole string, unreadOnly bool) []models.Notification {
	var rows *sql.Rows
	var err error

	if userRole == "" && !unreadOnly {
		rows, err = s.db.Query(
			`SELECT id, user_role, type, title, body, is_read, action_url, created_at
			 FROM notifications ORDER BY created_at DESC`)
	} else if userRole == "" {
		rows, err = s.db.Query(
			`SELECT id, user_role, type, title, body, is_read, action_url, created_at
			 FROM notifications WHERE is_read = 0 ORDER BY created_at DESC`)
	} else if !unreadOnly {
		rows, err = s.db.Query(
			`SELECT id, user_role, type, title, body, is_read, action_url, created_at
			 FROM notifications WHERE LOWER(user_role) = LOWER($1) ORDER BY created_at DESC`, userRole)
	} else {
		rows, err = s.db.Query(
			`SELECT id, user_role, type, title, body, is_read, action_url, created_at
			 FROM notifications WHERE LOWER(user_role) = LOWER($1) AND is_read = 0 ORDER BY created_at DESC`, userRole)
	}
	if err != nil {
		slog.Error("list notifications", "error", err)
		return nil
	}
	defer rows.Close()

	var notifs []models.Notification
	for rows.Next() {
		var n models.Notification
		var isRead int
		if err := rows.Scan(&n.ID, &n.UserRole, &n.Type, &n.Title, &n.Body,
			&isRead, &n.ActionURL, &n.CreatedAt); err != nil {
			slog.Error("scan notification", "error", err)
			continue
		}
		n.Read = isRead != 0
		notifs = append(notifs, n)
	}
	return notifs
}

func (s *SQLStore) MarkNotificationRead(id string) error {
	result, err := s.db.Exec(`UPDATE notifications SET is_read = 1 WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return fmt.Errorf("notification %s not found", id)
	}
	return nil
}

func (s *SQLStore) UnreadNotificationCount(userRole string) int {
	var count int
	var err error
	if userRole == "" {
		err = s.db.QueryRow(`SELECT COUNT(*) FROM notifications WHERE is_read = 0`).Scan(&count)
	} else {
		err = s.db.QueryRow(`SELECT COUNT(*) FROM notifications WHERE LOWER(user_role) = LOWER($1) AND is_read = 0`, userRole).Scan(&count)
	}
	if err != nil {
		return 0
	}
	return count
}

// --- Chat history operations (implements data.ChatPersister) ---

func (s *SQLStore) CreateChatSession(session models.ChatSession) error {
	if session.CreatedAt == "" {
		session.CreatedAt = time.Now().Format(time.RFC3339)
	}
	if session.UpdatedAt == "" {
		session.UpdatedAt = session.CreatedAt
	}
	if session.Title == "" {
		session.Title = "New conversation"
	}

	_, err := s.db.Exec(
		`INSERT INTO chat_sessions (id, user_id, user_role, company, title, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		session.ID, session.UserID, session.UserRole, session.Company,
		session.Title, session.CreatedAt, session.UpdatedAt,
	)
	return err
}

func (s *SQLStore) ListChatSessions(userID, userRole string) []models.ChatSession {
	var rows *sql.Rows
	var err error

	if userID != "" {
		rows, err = s.db.Query(
			`SELECT id, user_id, user_role, company, title, created_at, updated_at
			 FROM chat_sessions WHERE user_id = $1 ORDER BY updated_at DESC`, userID)
	} else {
		rows, err = s.db.Query(
			`SELECT id, user_id, user_role, company, title, created_at, updated_at
			 FROM chat_sessions WHERE user_role = $1 ORDER BY updated_at DESC`, userRole)
	}
	if err != nil {
		slog.Error("list chat sessions", "error", err)
		return nil
	}
	defer rows.Close()

	var sessions []models.ChatSession
	for rows.Next() {
		var cs models.ChatSession
		if err := rows.Scan(&cs.ID, &cs.UserID, &cs.UserRole, &cs.Company,
			&cs.Title, &cs.CreatedAt, &cs.UpdatedAt); err != nil {
			continue
		}
		sessions = append(sessions, cs)
	}
	return sessions
}

func (s *SQLStore) GetChatSession(id string) (*models.ChatSession, bool) {
	var cs models.ChatSession
	err := s.db.QueryRow(
		`SELECT id, user_id, user_role, company, title, created_at, updated_at
		 FROM chat_sessions WHERE id = $1`, id,
	).Scan(&cs.ID, &cs.UserID, &cs.UserRole, &cs.Company,
		&cs.Title, &cs.CreatedAt, &cs.UpdatedAt)
	if err != nil {
		return nil, false
	}
	return &cs, true
}

func (s *SQLStore) UpdateChatSessionTitle(id, title string) error {
	_, err := s.db.Exec(
		`UPDATE chat_sessions SET title = $1, updated_at = $2 WHERE id = $3`,
		title, time.Now().Format(time.RFC3339), id,
	)
	return err
}

func (s *SQLStore) AddChatMessage(sessionID string, msg models.ChatMessage) error {
	now := time.Now().Format(time.RFC3339)
	_, err := s.db.Exec(
		`INSERT INTO chat_messages (session_id, role, content, created_at)
		 VALUES ($1, $2, $3, $4)`,
		sessionID, msg.Role, msg.Content, now,
	)
	if err != nil {
		return err
	}
	// Touch session updated_at
	s.db.Exec(`UPDATE chat_sessions SET updated_at = $1 WHERE id = $2`, now, sessionID)
	return nil
}

func (s *SQLStore) GetChatMessages(sessionID string) []models.ChatMessage {
	rows, err := s.db.Query(
		`SELECT role, content FROM chat_messages WHERE session_id = $1 ORDER BY id ASC`, sessionID)
	if err != nil {
		slog.Error("get chat messages", "error", err)
		return nil
	}
	defer rows.Close()

	var msgs []models.ChatMessage
	for rows.Next() {
		var m models.ChatMessage
		if err := rows.Scan(&m.Role, &m.Content); err != nil {
			continue
		}
		msgs = append(msgs, m)
	}
	return msgs
}

func (s *SQLStore) DeleteChatSession(id string) error {
	s.db.Exec(`DELETE FROM chat_messages WHERE session_id = $1`, id)
	_, err := s.db.Exec(`DELETE FROM chat_sessions WHERE id = $1`, id)
	return err
}

// --- Helpers ---

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
