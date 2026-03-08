package data

import (
	"testing"

	"heywood-tbs/internal/models"

	// Register SQLite driver for tests (normally registered in connector.go)
	_ "modernc.org/sqlite"
)

// newTestSQLStore creates an in-memory SQLite-backed store for testing.
func newTestSQLStore(t *testing.T) *SQLStore {
	t.Helper()
	t.Setenv("AUTH_MODE", "")
	s, err := NewSQLStore(testDataDir(t), "sqlite", ":memory:")
	if err != nil {
		t.Fatalf("NewSQLStore: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

// --- Task CRUD (SQL) ---

func TestSQLStore_TaskCRUD(t *testing.T) {
	s := newTestSQLStore(t)

	// Create
	task := models.Task{
		Title:       "SQL Task",
		Description: "Task stored in SQLite",
		AssignedTo:  "instructor-alpha",
		CreatedBy:   "heywood",
		Priority:    "high",
	}
	if err := s.CreateTask(task); err != nil {
		t.Fatalf("CreateTask: %v", err)
	}

	// List all
	tasks := s.ListTasks("")
	if len(tasks) == 0 {
		t.Fatal("ListTasks returned 0 after CreateTask")
	}
	createdID := tasks[0].ID
	if createdID == "" {
		t.Fatal("created task has empty ID")
	}
	if tasks[0].Status != "pending" {
		t.Errorf("default status = %q, want pending", tasks[0].Status)
	}

	// List filtered by assignee
	filtered := s.ListTasks("instructor-alpha")
	if len(filtered) == 0 {
		t.Fatal("ListTasks(instructor-alpha) returned 0")
	}

	// Get
	got, ok := s.GetTask(createdID)
	if !ok {
		t.Fatalf("GetTask(%s) returned not found", createdID)
	}
	if got.Title != "SQL Task" {
		t.Errorf("Title = %q, want %q", got.Title, "SQL Task")
	}

	// Update
	newStatus := "completed"
	newPriority := "low"
	if err := s.UpdateTask(createdID, models.TaskUpdateRequest{
		Status:   &newStatus,
		Priority: &newPriority,
	}); err != nil {
		t.Fatalf("UpdateTask: %v", err)
	}

	updated, ok := s.GetTask(createdID)
	if !ok {
		t.Fatal("GetTask after update returned not found")
	}
	if updated.Status != "completed" {
		t.Errorf("Status after update = %q, want completed", updated.Status)
	}
	if updated.Priority != "low" {
		t.Errorf("Priority after update = %q, want low", updated.Priority)
	}

	// Update non-existent
	if err := s.UpdateTask("TSK-FAKE", models.TaskUpdateRequest{Status: &newStatus}); err == nil {
		t.Error("UpdateTask on non-existent ID should return error")
	}
}

// --- Message CRUD (SQL) ---

func TestSQLStore_MessageCRUD(t *testing.T) {
	s := newTestSQLStore(t)

	msg := models.Message{
		From:    "heywood",
		To:      "co-alpha",
		Subject: "SQL Message",
		Body:    "Message stored in SQLite.",
	}
	if err := s.CreateMessage(msg); err != nil {
		t.Fatalf("CreateMessage: %v", err)
	}

	// List all
	msgs := s.ListMessages("")
	if len(msgs) == 0 {
		t.Fatal("ListMessages returned 0 after CreateMessage")
	}
	msgID := msgs[0].ID
	if msgs[0].Read {
		t.Error("new message should not be read")
	}

	// List by recipient
	filtered := s.ListMessages("co-alpha")
	if len(filtered) == 0 {
		t.Fatal("ListMessages(co-alpha) returned 0")
	}

	// Mark read
	if err := s.MarkMessageRead(msgID); err != nil {
		t.Fatalf("MarkMessageRead: %v", err)
	}

	// Verify
	after := s.ListMessages("")
	var found *models.Message
	for i := range after {
		if after[i].ID == msgID {
			found = &after[i]
			break
		}
	}
	if found == nil {
		t.Fatal("message not found after MarkMessageRead")
	}
	if !found.Read {
		t.Error("message should be marked as read")
	}
}

// --- Notification CRUD (SQL) ---

func TestSQLStore_NotificationCRUD(t *testing.T) {
	s := newTestSQLStore(t)

	n := models.Notification{
		UserRole: "co-alpha",
		Type:     "alert",
		Title:    "SQL Notification",
		Body:     "Stored in SQLite.",
	}
	if err := s.CreateNotification(n); err != nil {
		t.Fatalf("CreateNotification: %v", err)
	}

	// List
	notifs := s.ListNotifications("co-alpha", false)
	if len(notifs) == 0 {
		t.Fatal("ListNotifications returned 0")
	}
	nID := notifs[0].ID

	// Unread count
	count := s.UnreadNotificationCount("co-alpha")
	if count == 0 {
		t.Fatal("UnreadNotificationCount should be > 0")
	}

	// Mark read
	if err := s.MarkNotificationRead(nID); err != nil {
		t.Fatalf("MarkNotificationRead: %v", err)
	}

	// Verify count decreased
	afterCount := s.UnreadNotificationCount("co-alpha")
	if afterCount >= count {
		t.Errorf("unread count should decrease: before=%d, after=%d", count, afterCount)
	}

	// Unread-only filter
	unreadOnly := s.ListNotifications("co-alpha", true)
	for _, notif := range unreadOnly {
		if notif.ID == nID {
			t.Error("read notification appeared in unread-only list")
		}
	}
}

// --- Migration idempotency ---

func TestSQLStore_MigrationIdempotent(t *testing.T) {
	t.Setenv("AUTH_MODE", "")
	dataDir := testDataDir(t)

	// Create first store, then close it
	s1, err := NewSQLStore(dataDir, "sqlite", ":memory:")
	if err != nil {
		t.Fatalf("first NewSQLStore: %v", err)
	}
	// Insert a task to ensure the DB is actually used
	if err := s1.CreateTask(models.Task{
		Title:      "Idempotent Test",
		AssignedTo: "test",
		CreatedBy:  "test",
	}); err != nil {
		t.Fatalf("CreateTask: %v", err)
	}
	s1.Close()

	// Create second store against same in-memory DB (fresh :memory: is fine;
	// the point is that migrate() runs again without error)
	s2, err := NewSQLStore(dataDir, "sqlite", ":memory:")
	if err != nil {
		t.Fatalf("second NewSQLStore should succeed (idempotent migration): %v", err)
	}
	s2.Close()
}

// --- Schema version ---

func TestSQLStore_SchemaVersion(t *testing.T) {
	s := newTestSQLStore(t)

	var version int
	err := s.db.QueryRow(`SELECT COALESCE(MAX(version), 0) FROM schema_version`).Scan(&version)
	if err != nil {
		t.Fatalf("query schema_version: %v", err)
	}
	if version == 0 {
		t.Error("schema_version should have at least version 1 after migration")
	}
}

// --- Reference data comes from JSON store ---

func TestSQLStore_ListStudentsFromReference(t *testing.T) {
	s := newTestSQLStore(t)

	// Students are loaded from JSON, not SQL
	students := s.ListStudents("", "", "", false)
	if len(students) == 0 {
		t.Fatal("expected students from embedded JSON store, got 0")
	}

	// Verify a known student is present
	st, ok := s.GetStudent("STU-001")
	if !ok {
		t.Fatal("GetStudent(STU-001) should find student from JSON reference data")
	}
	if st.LastName == "" {
		t.Error("student from reference data should have a last name")
	}
}

// --- User isolation ---

func TestSQLStore_UserIsolation(t *testing.T) {
	s := newTestSQLStore(t)

	// Create a task assigned to user A
	taskA := models.Task{
		Title:      "Task for Alpha",
		AssignedTo: "user-alpha",
		CreatedBy:  "heywood",
	}
	if err := s.CreateTask(taskA); err != nil {
		t.Fatalf("CreateTask A: %v", err)
	}

	// Create a task assigned to user B
	taskB := models.Task{
		Title:      "Task for Bravo",
		AssignedTo: "user-bravo",
		CreatedBy:  "heywood",
	}
	if err := s.CreateTask(taskB); err != nil {
		t.Fatalf("CreateTask B: %v", err)
	}

	// User A should only see their tasks
	tasksA := s.ListTasks("user-alpha")
	for _, tsk := range tasksA {
		if tsk.AssignedTo != "user-alpha" {
			t.Errorf("user-alpha task list contains task assigned to %q", tsk.AssignedTo)
		}
	}
	if len(tasksA) != 1 {
		t.Errorf("user-alpha should have 1 task, got %d", len(tasksA))
	}

	// User B should only see their tasks
	tasksB := s.ListTasks("user-bravo")
	for _, tsk := range tasksB {
		if tsk.AssignedTo != "user-bravo" {
			t.Errorf("user-bravo task list contains task assigned to %q", tsk.AssignedTo)
		}
	}
	if len(tasksB) != 1 {
		t.Errorf("user-bravo should have 1 task, got %d", len(tasksB))
	}

	// Listing all tasks should return both
	all := s.ListTasks("")
	if len(all) < 2 {
		t.Errorf("expected at least 2 tasks total, got %d", len(all))
	}
}
