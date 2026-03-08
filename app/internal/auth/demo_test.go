package auth

import (
	"net/http"
	"testing"
)

func TestDemoProvider_NoCookies(t *testing.T) {
	p := &DemoProvider{}
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	id := p.Authenticate(r)

	assertStr(t, "Role", RoleStaff, id.Role)
	assertStr(t, "ID", "demo-staff", id.ID)
	assertStr(t, "Source", "demo", id.Source)
}

func TestDemoProvider_ValidRoleCookie(t *testing.T) {
	tests := []struct {
		role string
		want string
	}{
		{"xo", RoleXO},
		{"staff", RoleStaff},
		{"spc", RoleSPC},
		{"student", RoleStudent},
	}

	p := &DemoProvider{}

	for _, tc := range tests {
		t.Run(tc.role, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, "/", nil)
			r.AddCookie(&http.Cookie{Name: "heywood-role", Value: tc.role})

			id := p.Authenticate(r)

			assertStr(t, "Role", tc.want, id.Role)
			assertStr(t, "ID", "demo-"+tc.role, id.ID)
		})
	}
}

func TestDemoProvider_InvalidRoleCookie(t *testing.T) {
	invalidRoles := []string{"admin", "root", "superuser", ""}

	p := &DemoProvider{}

	for _, role := range invalidRoles {
		t.Run("role_"+role, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, "/", nil)
			r.AddCookie(&http.Cookie{Name: "heywood-role", Value: role})

			id := p.Authenticate(r)

			// Invalid roles fall back to RoleStaff
			assertStr(t, "Role", RoleStaff, id.Role)
		})
	}
}

func TestDemoProvider_CompanyAndStudentID(t *testing.T) {
	p := &DemoProvider{}
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(&http.Cookie{Name: "heywood-role", Value: "student"})
	r.AddCookie(&http.Cookie{Name: "heywood-company", Value: "Bravo"})
	r.AddCookie(&http.Cookie{Name: "heywood-student-id", Value: "STU-042"})

	id := p.Authenticate(r)

	assertStr(t, "Role", RoleStudent, id.Role)
	assertStr(t, "Company", "Bravo", id.Company)
	// When student-id is set, ID becomes "demo-{studentId}"
	assertStr(t, "ID", "demo-STU-042", id.ID)
}

func TestDemoProvider_SupportsSwitch(t *testing.T) {
	p := &DemoProvider{}
	if !p.SupportsSwitch() {
		t.Error("DemoProvider.SupportsSwitch() = false, want true")
	}
}
