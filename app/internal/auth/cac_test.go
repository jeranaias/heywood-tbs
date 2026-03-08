package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// testRosterEntry mirrors the unexported rosterEntry for writing test fixtures.
type testRosterEntry struct {
	EDIPI   string `json:"edipi"`
	Name    string `json:"name"`
	Role    string `json:"role"`
	Company string `json:"company"`
	Email   string `json:"email"`
}

// makeTestCert generates a self-signed X.509 certificate with the given Common
// Name and returns its PEM encoding, URL-encoded, ready for the
// X-ARR-ClientCert header.
func makeTestCert(t *testing.T, cn string) string {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: cn},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("create certificate: %v", err)
	}

	pemBlock := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	return url.QueryEscape(string(pemBlock))
}

// makeRosterFile writes a JSON roster to a temporary file and returns its path.
func makeRosterFile(t *testing.T, entries []testRosterEntry) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "roster.json")

	data, err := json.Marshal(entries)
	if err != nil {
		t.Fatalf("marshal roster: %v", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("write roster file: %v", err)
	}
	return path
}

func TestCACProvider_NoCertHeader(t *testing.T) {
	p := NewCACProvider("/nonexistent/roster.json")
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	id := p.Authenticate(r)

	assertStr(t, "Role", RoleUnauthorized, id.Role)
	assertStr(t, "ID", "unknown", id.ID)
	assertStr(t, "Source", "cac", id.Source)
}

func TestCACProvider_InvalidURLEncoding(t *testing.T) {
	p := NewCACProvider("/nonexistent/roster.json")
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("X-ARR-ClientCert", "%ZZnot-valid-url-encoding")

	id := p.Authenticate(r)

	assertStr(t, "Role", RoleUnauthorized, id.Role)
	assertStr(t, "ID", "unknown", id.ID)
}

func TestCACProvider_InvalidPEM(t *testing.T) {
	p := NewCACProvider("/nonexistent/roster.json")
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	// URL-encode garbage that is not valid PEM
	r.Header.Set("X-ARR-ClientCert", url.QueryEscape("this is not PEM data"))

	id := p.Authenticate(r)

	assertStr(t, "Role", RoleUnauthorized, id.Role)
	assertStr(t, "ID", "unknown", id.ID)
}

func TestCACProvider_ShortCommonName(t *testing.T) {
	// CN with fewer than 4 dot-separated parts
	tests := []struct {
		name string
		cn   string
	}{
		{"one_part", "SMITH"},
		{"two_parts", "SMITH.JOHN"},
		{"three_parts", "SMITH.JOHN.A"},
	}

	p := NewCACProvider("/nonexistent/roster.json")

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			certHeader := makeTestCert(t, tc.cn)
			r, _ := http.NewRequest(http.MethodGet, "/", nil)
			r.Header.Set("X-ARR-ClientCert", certHeader)

			id := p.Authenticate(r)

			assertStr(t, "Role", RoleUnauthorized, id.Role)
			assertStr(t, "Source", "cac", id.Source)
		})
	}
}

func TestCACProvider_EDIPIInRoster(t *testing.T) {
	roster := []testRosterEntry{
		{
			EDIPI:   "1234567890",
			Name:    "Maj Smith, John",
			Role:    RoleXO,
			Company: "Alpha",
			Email:   "john.smith@usmc.mil",
		},
		{
			EDIPI:   "9876543210",
			Name:    "SSgt Jones, Bob",
			Role:    RoleSPC,
			Company: "Bravo",
			Email:   "bob.jones@usmc.mil",
		},
	}
	rosterPath := makeRosterFile(t, roster)
	p := NewCACProvider(rosterPath)

	tests := []struct {
		name        string
		cn          string
		wantRole    string
		wantCompany string
		wantID      string
		wantName    string
	}{
		{
			name:        "xo_role",
			cn:          "SMITH.JOHN.A.1234567890",
			wantRole:    RoleXO,
			wantCompany: "Alpha",
			wantID:      "1234567890",
			wantName:    "Maj Smith, John",
		},
		{
			name:        "spc_role",
			cn:          "JONES.BOB.B.9876543210",
			wantRole:    RoleSPC,
			wantCompany: "Bravo",
			wantID:      "9876543210",
			wantName:    "SSgt Jones, Bob",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			certHeader := makeTestCert(t, tc.cn)
			r, _ := http.NewRequest(http.MethodGet, "/", nil)
			r.Header.Set("X-ARR-ClientCert", certHeader)

			id := p.Authenticate(r)

			assertStr(t, "Role", tc.wantRole, id.Role)
			assertStr(t, "Company", tc.wantCompany, id.Company)
			assertStr(t, "ID", tc.wantID, id.ID)
			assertStr(t, "Name", tc.wantName, id.Name)
			assertStr(t, "Source", "cac", id.Source)
		})
	}
}

func TestCACProvider_EDIPINotInRoster(t *testing.T) {
	roster := []testRosterEntry{
		{EDIPI: "1111111111", Name: "Only User", Role: RoleStaff, Company: "Alpha"},
	}
	rosterPath := makeRosterFile(t, roster)
	p := NewCACProvider(rosterPath)

	certHeader := makeTestCert(t, "DOE.JANE.M.9999999999")
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("X-ARR-ClientCert", certHeader)

	id := p.Authenticate(r)

	assertStr(t, "Role", RoleUnauthorized, id.Role)
	assertStr(t, "ID", "9999999999", id.ID)
	assertStr(t, "Name", "DOE, JANE", id.Name)
	assertStr(t, "Source", "cac", id.Source)
}

func TestCACProvider_EmptyRoster(t *testing.T) {
	// Provider created with a nonexistent roster file — no panic, all users unauthorized
	p := NewCACProvider(filepath.Join(t.TempDir(), "does-not-exist.json"))

	certHeader := makeTestCert(t, "SMITH.JOHN.A.1234567890")
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("X-ARR-ClientCert", certHeader)

	id := p.Authenticate(r)

	assertStr(t, "Role", RoleUnauthorized, id.Role)
	assertStr(t, "ID", "1234567890", id.ID)
	assertStr(t, "Source", "cac", id.Source)
}

// assertStr is a test helper that compares two strings.
func assertStr(t *testing.T, field, want, got string) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %q, want %q", field, got, want)
	}
}
