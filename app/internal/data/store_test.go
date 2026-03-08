package data

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"heywood-tbs/internal/models"
)

// testDataDir returns the absolute path to the app/data directory.
func testDataDir(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("unable to determine test file location")
	}
	// file is .../app/internal/data/store_test.go
	dir := filepath.Join(filepath.Dir(file), "..", "..", "data")
	abs, err := filepath.Abs(dir)
	if err != nil {
		t.Fatalf("abs path: %v", err)
	}
	if _, err := os.Stat(abs); err != nil {
		t.Fatalf("data directory not found at %s: %v", abs, err)
	}
	return abs
}

// newTestStore creates a Store in demo mode (in-memory mutable data).
func newTestStore(t *testing.T) *Store {
	t.Helper()
	// Force demo mode so mutable operations stay in-memory
	t.Setenv("AUTH_MODE", "")
	s, err := NewStore(testDataDir(t))
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return s
}

// --- Student tests ---

func TestStore_ListStudents(t *testing.T) {
	s := newTestStore(t)
	students := s.ListStudents("", "", "", false)
	if len(students) == 0 {
		t.Fatal("expected at least one student, got 0")
	}
	// Verify students have IDs
	for _, st := range students {
		if st.ID == "" {
			t.Error("student has empty ID")
		}
	}
}

func TestStore_ListStudents_FilterCompany(t *testing.T) {
	s := newTestStore(t)
	students := s.ListStudents("Alpha", "", "", false)
	if len(students) == 0 {
		t.Fatal("expected Alpha company students, got 0")
	}
	for _, st := range students {
		if !strings.EqualFold(st.Company, "Alpha") {
			t.Errorf("student %s has company %q, want Alpha", st.ID, st.Company)
		}
	}
}

func TestStore_ListStudents_FilterSearch(t *testing.T) {
	s := newTestStore(t)
	// Get a known student name to search for
	all := s.ListStudents("", "", "", false)
	if len(all) == 0 {
		t.Fatal("no students loaded")
	}
	target := all[0]
	query := strings.ToLower(target.LastName)

	results := s.ListStudents("", "", query, false)
	if len(results) == 0 {
		t.Fatalf("search for %q returned 0 results", query)
	}
	found := false
	for _, st := range results {
		if st.ID == target.ID {
			found = true
		}
	}
	if !found {
		t.Errorf("search for %q did not return student %s", query, target.ID)
	}
}

func TestStore_GetStudent(t *testing.T) {
	s := newTestStore(t)
	st, ok := s.GetStudent("STU-001")
	if !ok {
		t.Fatal("GetStudent(STU-001) returned not found")
	}
	if st.ID != "STU-001" {
		t.Errorf("got ID %q, want STU-001", st.ID)
	}
	if st.LastName == "" {
		t.Error("student has empty LastName")
	}
}

func TestStore_GetStudent_NotFound(t *testing.T) {
	s := newTestStore(t)
	_, ok := s.GetStudent("STU-NONEXISTENT")
	if ok {
		t.Error("GetStudent should return false for unknown ID")
	}
}

func TestStore_StudentStats(t *testing.T) {
	s := newTestStore(t)
	stats := s.StudentStats("")
	if stats.ActiveStudents == 0 {
		t.Fatal("expected non-zero ActiveStudents")
	}
	if stats.ByPhase == nil {
		t.Error("ByPhase map is nil")
	}
}

func TestStore_AtRiskStudents(t *testing.T) {
	s := newTestStore(t)
	atRisk := s.AtRiskStudents("")
	// Verify every returned student is actually at-risk
	for _, st := range atRisk {
		if !st.AtRisk {
			t.Errorf("student %s returned by AtRiskStudents but AtRisk is false", st.ID)
		}
	}
}

// --- Instructor tests ---

func TestStore_ListInstructors(t *testing.T) {
	s := newTestStore(t)
	instructors := s.ListInstructors("")
	if len(instructors) == 0 {
		t.Fatal("expected at least one instructor, got 0")
	}
	for _, inst := range instructors {
		if inst.ID == "" {
			t.Error("instructor has empty ID")
		}
	}
}

func TestStore_GetInstructor(t *testing.T) {
	s := newTestStore(t)
	inst, ok := s.GetInstructor("INST-001")
	if !ok {
		t.Fatal("GetInstructor(INST-001) returned not found")
	}
	if inst.ID != "INST-001" {
		t.Errorf("got ID %q, want INST-001", inst.ID)
	}
	if inst.LastName == "" {
		t.Error("instructor has empty LastName")
	}
}

// --- Schedule tests ---

func TestStore_ListSchedule(t *testing.T) {
	s := newTestStore(t)
	events := s.ListSchedule("")
	if len(events) == 0 {
		t.Fatal("expected at least one schedule event, got 0")
	}
	for _, evt := range events {
		if evt.ID == "" {
			t.Error("schedule event has empty ID")
		}
	}
}

// --- Task CRUD ---

func TestStore_TaskCRUD(t *testing.T) {
	s := newTestStore(t)

	// Create
	task := models.Task{
		Title:       "Test Task",
		Description: "Unit test task",
		AssignedTo:  "instructor-alpha",
		CreatedBy:   "heywood",
		Priority:    "high",
	}
	if err := s.CreateTask(task); err != nil {
		t.Fatalf("CreateTask: %v", err)
	}

	// List
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
	if tasks[0].CreatedAt == "" {
		t.Error("CreatedAt should be auto-populated")
	}

	// Get
	got, ok := s.GetTask(createdID)
	if !ok {
		t.Fatalf("GetTask(%s) returned not found", createdID)
	}
	if got.Title != "Test Task" {
		t.Errorf("Title = %q, want %q", got.Title, "Test Task")
	}

	// Update
	newStatus := "completed"
	newPriority := "low"
	err := s.UpdateTask(createdID, models.TaskUpdateRequest{
		Status:   &newStatus,
		Priority: &newPriority,
	})
	if err != nil {
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

	// Update non-existent task
	if err := s.UpdateTask("TSK-FAKE", models.TaskUpdateRequest{Status: &newStatus}); err == nil {
		t.Error("UpdateTask on non-existent ID should return error")
	}
}

// --- Message CRUD ---

func TestStore_MessageCRUD(t *testing.T) {
	s := newTestStore(t)

	msg := models.Message{
		From:    "heywood",
		To:      "co-alpha",
		Subject: "Test Message",
		Body:    "This is a test message body.",
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
	if msgID == "" {
		t.Fatal("message has empty ID")
	}
	if msgs[0].Read {
		t.Error("new message should not be read")
	}

	// List filtered by recipient
	filtered := s.ListMessages("co-alpha")
	if len(filtered) == 0 {
		t.Fatal("ListMessages(co-alpha) returned 0")
	}

	// Mark read
	if err := s.MarkMessageRead(msgID); err != nil {
		t.Fatalf("MarkMessageRead: %v", err)
	}

	// Verify read status
	updated := s.ListMessages("")
	var foundMsg *models.Message
	for i := range updated {
		if updated[i].ID == msgID {
			foundMsg = &updated[i]
			break
		}
	}
	if foundMsg == nil {
		t.Fatal("message not found after MarkMessageRead")
	}
	if !foundMsg.Read {
		t.Error("message should be marked as read")
	}
}

// --- Notification CRUD ---

func TestStore_NotificationCRUD(t *testing.T) {
	s := newTestStore(t)

	n := models.Notification{
		UserRole: "co-alpha",
		Type:     "alert",
		Title:    "Test Notification",
		Body:     "Notification body text.",
	}
	if err := s.CreateNotification(n); err != nil {
		t.Fatalf("CreateNotification: %v", err)
	}

	// List
	notifs := s.ListNotifications("co-alpha", false)
	if len(notifs) == 0 {
		t.Fatal("ListNotifications returned 0 after CreateNotification")
	}
	nID := notifs[0].ID
	if nID == "" {
		t.Fatal("notification has empty ID")
	}

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
		t.Errorf("unread count should decrease after MarkNotificationRead: before=%d, after=%d", count, afterCount)
	}

	// Unread-only filter should exclude the read notification
	unreadOnly := s.ListNotifications("co-alpha", true)
	for _, notif := range unreadOnly {
		if notif.ID == nID {
			t.Error("read notification appeared in unread-only list")
		}
	}
}

// --- TotalStudentCount ---

func TestStore_TotalStudentCount(t *testing.T) {
	s := newTestStore(t)
	count := s.TotalStudentCount()
	if count == 0 {
		t.Fatal("TotalStudentCount returned 0")
	}
	// Should match length of ListStudents with no filters
	all := s.ListStudents("", "", "", false)
	if count != len(all) {
		t.Errorf("TotalStudentCount=%d, but ListStudents returned %d", count, len(all))
	}
}
