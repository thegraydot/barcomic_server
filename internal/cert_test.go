package barcomic

import (
	"crypto/x509"
	"testing"
	"time"
)

func TestGenerateTLSCertificate(t *testing.T) {
	const testAddr = "testhost"
	const testPort = "9999"

	// Use cert.go to generate a TLS certificate
	tlsCert := GenerateTLSCertificate(testAddr, testPort)

	// Parse the generated TLS certificate
	cert, err := x509.ParseCertificate(tlsCert.Certificate[0])
	if err != nil {
		panic(err)
	}

	// Check cert common name
	serverHost := testAddr + ":" + testPort
	if cert.Subject.CommonName != serverHost {
		t.Fatalf("Expected CommonName: %s, but got: %s", serverHost, cert.Subject.CommonName)
	}

	// Assert 1-year validity window
	expected := time.Now().AddDate(1, 0, 0)
	if cert.NotAfter.Before(expected.Add(-time.Minute)) ||
		cert.NotAfter.After(expected.Add(time.Minute)) {
		t.Fatalf("Expected NotAfter ~1 year from now, got %v", cert.NotAfter)
	}
}
