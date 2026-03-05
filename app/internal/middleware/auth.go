package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const (
	RoleKey      contextKey = "role"
	CompanyKey   contextKey = "company"
	StudentIDKey contextKey = "studentId"
)

// Auth reads role info from cookies and injects into request context.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		ctx := context.WithValue(r.Context(), RoleKey, role)
		ctx = context.WithValue(ctx, CompanyKey, company)
		ctx = context.WithValue(ctx, StudentIDKey, studentID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRole extracts the role from request context.
func GetRole(ctx context.Context) string {
	if v, ok := ctx.Value(RoleKey).(string); ok {
		return v
	}
	return "staff"
}

// GetCompany extracts the company filter from request context.
func GetCompany(ctx context.Context) string {
	if v, ok := ctx.Value(CompanyKey).(string); ok {
		return v
	}
	return ""
}

// GetStudentID extracts the student ID from request context.
func GetStudentID(ctx context.Context) string {
	if v, ok := ctx.Value(StudentIDKey).(string); ok {
		return v
	}
	return ""
}
