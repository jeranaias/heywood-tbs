package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"heywood-tbs/internal/auth"
)

func TestGraphEndpoints_NonPrivilegedBlocked(t *testing.T) {
	h := newTestHandler(t)

	tests := []struct {
		name    string
		method  string
		path    string
		handler func(http.ResponseWriter, *http.Request)
	}{
		{"handleGraphTest", "POST", "/api/v1/graph/test", h.handleGraphTest},
		{"handleSharePointSite", "GET", "/api/v1/graph/sharepoint/site?url=x", h.handleSharePointSite},
		{"handleSharePointLists", "GET", "/api/v1/graph/sharepoint/lists?siteId=x", h.handleSharePointLists},
		{"handleSharePointListItems", "GET", "/api/v1/graph/sharepoint/list-items?siteId=x&listId=x", h.handleSharePointListItems},
		{"handleSharePointDrives", "GET", "/api/v1/graph/sharepoint/drives?siteId=x", h.handleSharePointDrives},
		{"handleSharePointFiles", "GET", "/api/v1/graph/sharepoint/files?siteId=x&driveId=x", h.handleSharePointFiles},
		{"handleTeamsList", "GET", "/api/v1/graph/teams", h.handleTeamsList},
		{"handleTeamsChannels", "GET", "/api/v1/graph/teams/channels?teamId=x", h.handleTeamsChannels},
		{"handleTeamsFiles", "GET", "/api/v1/graph/teams/files?teamId=x&channelId=x", h.handleTeamsFiles},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var req *http.Request
			if tc.method == "POST" {
				req = httptest.NewRequest(tc.method, tc.path, strings.NewReader("{}"))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tc.method, tc.path, nil)
			}
			req = withRoleContext(req, auth.RoleSPC)
			rec := httptest.NewRecorder()

			tc.handler(rec, req)

			if rec.Code != 403 {
				t.Errorf("expected 403 for SPC role, got %d: %s", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestGraphEndpoints_NilClient503(t *testing.T) {
	h := newTestHandler(t)
	// graphClient, sharePointSvc, teamsSvc are nil by default in newTestHandler

	tests := []struct {
		name    string
		method  string
		path    string
		handler func(http.ResponseWriter, *http.Request)
	}{
		{"handleSharePointSite", "GET", "/api/v1/graph/sharepoint/site?url=test", h.handleSharePointSite},
		{"handleSharePointLists", "GET", "/api/v1/graph/sharepoint/lists?siteId=test", h.handleSharePointLists},
		{"handleSharePointDrives", "GET", "/api/v1/graph/sharepoint/drives?siteId=test", h.handleSharePointDrives},
		{"handleSharePointFiles", "GET", "/api/v1/graph/sharepoint/files?siteId=test&driveId=test", h.handleSharePointFiles},
		{"handleTeamsList", "GET", "/api/v1/graph/teams", h.handleTeamsList},
		{"handleTeamsChannels", "GET", "/api/v1/graph/teams/channels?teamId=test", h.handleTeamsChannels},
		{"handleTeamsFiles", "GET", "/api/v1/graph/teams/files?teamId=test&channelId=test", h.handleTeamsFiles},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			req = withRoleContext(req, auth.RoleXO)
			rec := httptest.NewRecorder()

			tc.handler(rec, req)

			if rec.Code != 503 {
				t.Errorf("expected 503 for nil service, got %d: %s", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestGraphTest_NotConfigured(t *testing.T) {
	h := newTestHandler(t)
	// graphClient is nil by default

	req := httptest.NewRequest("POST", "/api/v1/graph/test", strings.NewReader("{}"))
	req.Header.Set("Content-Type", "application/json")
	req = withRoleContext(req, auth.RoleXO)
	rec := httptest.NewRecorder()

	h.handleGraphTest(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	body := rec.Body.String()
	if !strings.Contains(body, "not_configured") {
		t.Errorf("expected response to contain %q, got %s", "not_configured", body)
	}
}
