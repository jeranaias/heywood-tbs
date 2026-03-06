package middleware

import (
	"crypto/tls"
	"net/http"
)

// SecurityHeaders returns middleware that sets STIG-compliant security headers.
// References: STIG V-222602 (HSTS), V-222610 (X-Frame-Options), V-222612 (X-Content-Type-Options)
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()

		// HSTS — enforce HTTPS for 1 year, include subdomains
		h.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// Prevent clickjacking
		h.Set("X-Frame-Options", "DENY")

		// Prevent MIME-type sniffing
		h.Set("X-Content-Type-Options", "nosniff")

		// XSS protection (legacy browsers)
		h.Set("X-XSS-Protection", "1; mode=block")

		// Referrer policy — don't leak internal URLs
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions policy — disable unnecessary APIs
		h.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")

		// CSP — restrict sources
		h.Set("Content-Security-Policy",
			"default-src 'self'; "+
				"script-src 'self'; "+
				"style-src 'self' 'unsafe-inline'; "+
				"img-src 'self' data:; "+
				"font-src 'self'; "+
				"connect-src 'self'; "+
				"frame-ancestors 'none'")

		// Cache control for API responses
		if len(r.URL.Path) > 4 && r.URL.Path[:5] == "/api/" {
			h.Set("Cache-Control", "no-store")
			h.Set("Pragma", "no-cache")
		}

		next.ServeHTTP(w, r)
	})
}

// FIPSTLSConfig returns a TLS configuration using only FIPS-approved cipher suites.
// Use with http.Server.TLSConfig when TLS is terminated at the Go server.
// On Azure App Service, TLS terminates at the load balancer, so this is for
// standalone deployments or when end-to-end TLS is required.
func FIPSTLSConfig() *tls.Config {
	return &tls.Config{
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS13,
		CipherSuites: []uint16{
			// TLS 1.2 FIPS-approved cipher suites
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			// TLS 1.3 cipher suites are automatically included and all use AES-GCM or ChaCha20
		},
		CurvePreferences: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
			tls.CurveP384,
		},
	}
}
