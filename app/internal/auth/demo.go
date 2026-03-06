package auth

import "net/http"

// DemoProvider reads role information from cookies set by the role picker UI.
// This is the default provider for demo and development environments.
type DemoProvider struct{}

func (p *DemoProvider) Authenticate(r *http.Request) *UserIdentity {
	role := "staff"
	if c, err := r.Cookie("heywood-role"); err == nil && c.Value != "" {
		role = c.Value
	}

	company := ""
	if c, err := r.Cookie("heywood-company"); err == nil && c.Value != "" {
		company = c.Value
	}

	studentID := ""
	if c, err := r.Cookie("heywood-student-id"); err == nil && c.Value != "" {
		studentID = c.Value
	}

	id := "demo-" + role
	if studentID != "" {
		id = "demo-" + studentID
	}

	return &UserIdentity{
		ID:      id,
		Role:    role,
		Company: company,
		Source:  "demo",
	}
}

func (p *DemoProvider) SupportsSwitch() bool { return true }
