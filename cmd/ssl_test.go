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

package cmd

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
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestSSLCmd_StatusAndRetrieve(t *testing.T) {
	certPEM, keyPEM := testCertPairPEM(t)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ssl/retrieve/example.com" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status":           "SUCCESS",
			"certificatechain": certPEM,
			"privatekey":       keyPEM,
		})
	})

	// 1. steamer ssl status --json (must never leak private key).
	statusOut, _, err := executeWithMockHandler(t, handler, nil, "ssl", "status", "example.com", "--json")
	if err != nil {
		t.Fatalf("ssl status --json unexpected error: %v", err)
	}
	if strings.Contains(statusOut, "PRIVATE KEY") {
		t.Fatalf("ssl status --json leaked private key to stdout")
	}
	if !strings.Contains(statusOut, "*.example.com") {
		t.Errorf("ssl status --json stdout = %q, want SAN %q", statusOut, "*.example.com")
	}

	// 2. steamer ssl retrieve writes 0644 cert + 0600 key, then skips rewrite when unchanged.
	outDir := t.TempDir()
	retOut, _, err := executeWithMockHandler(t, handler, nil, "ssl", "retrieve", "example.com", "--out-dir", outDir, "--json")
	if err != nil {
		t.Fatalf("ssl retrieve #1 unexpected error: %v", err)
	}
	var res1 sslRetrieveResult
	if err := json.Unmarshal([]byte(retOut), &res1); err != nil {
		t.Fatalf("json.Unmarshal(ssl retrieve #1 %q) unexpected error: %v", retOut, err)
	}
	if !res1.CertUpdated || !res1.KeyUpdated {
		t.Errorf("ssl retrieve #1 result = %+v, want CertUpdated=true and KeyUpdated=true", res1)
	}
	keyInfo, err := os.Stat(filepath.Join(outDir, "privkey.pem"))
	if err != nil {
		t.Fatalf("os.Stat(privkey.pem) unexpected error: %v", err)
	}
	if got, want := keyInfo.Mode().Perm(), os.FileMode(0o600); got != want {
		t.Errorf("privkey.pem mode = %04o, want %04o", got, want)
	}

	// Second retrieve with identical bundle should report updated=false.
	retOut2, _, err := executeWithMockHandler(t, handler, nil, "ssl", "retrieve", "example.com", "--out-dir", outDir, "--json")
	if err != nil {
		t.Fatalf("ssl retrieve #2 unexpected error: %v", err)
	}
	var res2 sslRetrieveResult
	if err := json.Unmarshal([]byte(retOut2), &res2); err != nil {
		t.Fatalf("json.Unmarshal(ssl retrieve #2 %q) unexpected error: %v", retOut2, err)
	}
	if res2.CertUpdated || res2.KeyUpdated {
		t.Errorf("ssl retrieve #2 result = %+v, want CertUpdated=false and KeyUpdated=false", res2)
	}
}

func testCertPairPEM(t *testing.T) (string, string) {
	t.Helper()
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(99),
		Subject:      pkix.Name{CommonName: "*.example.com"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(45 * 24 * time.Hour),
		DNSNames:     []string{"*.example.com", "example.com"},
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &priv.PublicKey, priv)
	if err != nil {
		t.Fatal(err)
	}
	keyDER, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		t.Fatal(err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})),
		string(pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER}))
}
