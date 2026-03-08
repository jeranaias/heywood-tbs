package auth

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// rosterEntry maps an EDIPI to a role in data/user-roster.json.
type rosterEntry struct {
	EDIPI   string `json:"edipi"`
	Name    string `json:"name"`
	Role    string `json:"role"`
	Company string `json:"company"`
	Email   string `json:"email"`
}

// CACProvider authenticates users via X.509 client certificates forwarded by
// Azure App Service (or any reverse proxy) in the X-ARR-ClientCert header.
// It parses the cert's Common Name (LAST.FIRST.MI.EDIPI) and looks up the
// EDIPI in a roster file to determine the user's role.
type CACProvider struct {
	roster map[string]*rosterEntry // EDIPI -> entry
}

// NewCACProvider loads the user roster from the given JSON file path.
// If the file doesn't exist or is empty, all users get RoleUnauthorized.
func NewCACProvider(rosterPath string) *CACProvider {
	p := &CACProvider{roster: make(map[string]*rosterEntry)}

	data, err := os.ReadFile(rosterPath)
	if err != nil {
		slog.Warn("CAC roster not found, all users will be unauthorized", "path", rosterPath, "error", err)
		return p
	}

	var entries []rosterEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		slog.Error("failed to parse CAC roster", "path", rosterPath, "error", err)
		return p
	}

	for i := range entries {
		p.roster[entries[i].EDIPI] = &entries[i]
	}
	slog.Info("CAC roster loaded", "entries", len(p.roster))
	return p
}

func (p *CACProvider) Authenticate(r *http.Request) *UserIdentity {
	// Azure App Service forwards the client cert as URL-encoded PEM
	certHeader := r.Header.Get("X-ARR-ClientCert")
	if certHeader == "" {
		slog.Debug("no client cert header")
		return &UserIdentity{ID: "unknown", Role: RoleUnauthorized, Source: "cac"}
	}

	// Decode URL encoding
	decoded, err := url.QueryUnescape(certHeader)
	if err != nil {
		slog.Error("failed to URL-decode client cert", "error", err)
		return &UserIdentity{ID: "unknown", Role: RoleUnauthorized, Source: "cac"}
	}

	// Parse PEM block
	block, _ := pem.Decode([]byte(decoded))
	if block == nil {
		// Try raw DER (some proxies send base64-encoded DER without PEM wrapper)
		slog.Error("failed to decode PEM from client cert header")
		return &UserIdentity{ID: "unknown", Role: RoleUnauthorized, Source: "cac"}
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		slog.Error("failed to parse X.509 certificate", "error", err)
		return &UserIdentity{ID: "unknown", Role: RoleUnauthorized, Source: "cac"}
	}

	// DoD CAC Common Name format: LAST.FIRST.MI.EDIPI
	cn := cert.Subject.CommonName
	parts := strings.Split(cn, ".")
	if len(parts) < 4 {
		slog.Warn("unexpected CN format", "cn", cn)
		return &UserIdentity{ID: cn, Name: cn, Role: RoleUnauthorized, Source: "cac"}
	}

	edipi := parts[len(parts)-1]
	lastName := parts[0]
	firstName := parts[1]
	name := lastName + ", " + firstName

	// Extract email from SAN if available
	email := ""
	if len(cert.EmailAddresses) > 0 {
		email = cert.EmailAddresses[0]
	}

	// Look up roster
	if entry, ok := p.roster[edipi]; ok {
		return &UserIdentity{
			ID:      edipi,
			Name:    entry.Name,
			Role:    entry.Role,
			Company: entry.Company,
			Email:   entry.Email,
			Source:  "cac",
		}
	}

	// EDIPI not in roster — deny access
	slog.Warn("EDIPI not found in roster", "edipi", edipi)
	return &UserIdentity{
		ID:     edipi,
		Name:   name,
		Role:   RoleUnauthorized,
		Email:  email,
		Source: "cac",
	}
}

func (p *CACProvider) SupportsSwitch() bool { return false }
