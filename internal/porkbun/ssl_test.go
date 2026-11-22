// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package porkbun

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
	"testing"
	"time"
)

func TestClient_RetrieveSSLBundleAndParse(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	certPEM, keyPEM := generateTestCertPEM(t, []string{"*.example.com", "example.com"}, now.Add(60*24*time.Hour))

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ssl/retrieve/example.com" {
			t.Errorf("request URL.Path = %q, want %q", r.URL.Path, "/ssl/retrieve/example.com")
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status":           "SUCCESS",
			"certificatechain": certPEM,
			"privatekey":       keyPEM,
		})
	})

	client := newTestClient(handler)

	bundle, err := client.RetrieveSSLBundle(t.Context(), "example.com")
	if err != nil {
		t.Fatalf("RetrieveSSLBundle(ctx, %q) unexpected error: %v", "example.com", err)
	}
	info, err := ParseCertificateChain(bundle.CertificateChain, now)
	if err != nil {
		t.Fatalf("ParseCertificateChain(chainPEM, %v) unexpected error: %v", now, err)
	}
	if got, want := info.ChainCount, 1; got != want {
		t.Errorf("ParseCertificateChain(...).ChainCount = %d, want %d", got, want)
	}
	if info.DaysRemaining < 59 || info.DaysRemaining > 60 {
		t.Errorf("ParseCertificateChain(...).DaysRemaining = %d, want 59 or 60", info.DaysRemaining)
	}
	if len(info.DNSNames) != 2 || info.DNSNames[0] != "*.example.com" {
		t.Errorf("ParseCertificateChain(...).DNSNames = %v, want [*.example.com example.com]", info.DNSNames)
	}
	if info.SHA256Fingerprint == "" {
		t.Errorf("ParseCertificateChain(...).SHA256Fingerprint = %q, want non-empty hex digest", info.SHA256Fingerprint)
	}
}

func generateTestCertPEM(t *testing.T, dnsNames []string, notAfter time.Time) (string, string) {
	t.Helper()
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(42),
		Subject:      pkix.Name{CommonName: dnsNames[0]},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     notAfter,
		DNSNames:     dnsNames,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &priv.PublicKey, priv)
	if err != nil {
		t.Fatal(err)
	}
	keyDER, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		t.Fatal(err)
	}
	certPEM := string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}))
	keyPEM := string(pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER}))
	return certPEM, keyPEM
}
