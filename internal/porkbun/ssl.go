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
	"context"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"math"
	"net/url"
	"time"
)

// SSLBundleResponse is the response from the Porkbun ssl/retrieve endpoint.
type SSLBundleResponse struct {
	APIResponse
	IntermediateCertificate string `json:"intermediatecertificate,omitempty"`
	CertificateChain        string `json:"certificatechain"`
	PrivateKey              string `json:"privatekey"`
	PublicKey               string `json:"publickey,omitempty"`
}

// SSLCertInfo holds parsed X.509 metadata for the leaf certificate and chain.
type SSLCertInfo struct {
	Subject           string    `json:"subject"`
	Issuer            string    `json:"issuer"`
	DNSNames          []string  `json:"dns_names"`
	NotBefore         time.Time `json:"not_before"`
	NotAfter          time.Time `json:"not_after"`
	DaysRemaining     int       `json:"days_remaining"`
	ChainCount        int       `json:"chain_count"`
	SHA256Fingerprint string    `json:"sha256_fingerprint"`
}

// RetrieveSSLBundle fetches the managed SSL certificate chain and private key for a domain.
func (c *Client) RetrieveSSLBundle(ctx context.Context, domain string) (*SSLBundleResponse, error) {
	if err := validatePathSegment("domain", domain); err != nil {
		return nil, err
	}
	var res SSLBundleResponse
	endpoint := fmt.Sprintf("ssl/retrieve/%s", url.PathEscape(domain))
	if err := c.post(ctx, endpoint, c.baseRequest(), &res); err != nil {
		return nil, err
	}
	if err := checkAPIStatus("retrieving ssl bundle failed", res.APIResponse); err != nil {
		return nil, err
	}
	if res.CertificateChain == "" || res.PrivateKey == "" {
		return nil, errors.New("ssl bundle response missing certificate chain or private key")
	}
	return &res, nil
}

// ParseCertificateChain parses a PEM-encoded certificate chain and returns metadata for the leaf certificate.
func ParseCertificateChain(chainPEM string, now time.Time) (*SSLCertInfo, error) {
	if now.IsZero() {
		now = time.Now().UTC()
	}
	rest := []byte(chainPEM)
	var certs []*x509.Certificate
	for len(rest) > 0 {
		var blk *pem.Block
		blk, rest = pem.Decode(rest)
		if blk == nil {
			break
		}
		if blk.Type != "CERTIFICATE" {
			continue
		}
		cert, err := x509.ParseCertificate(blk.Bytes)
		if err != nil {
			return nil, fmt.Errorf("parsing x509 certificate: %w", err)
		}
		certs = append(certs, cert)
	}
	if len(certs) == 0 {
		return nil, errors.New("no PEM certificates found in certificate chain")
	}

	leaf := certs[0]
	sum := sha256.Sum256(leaf.Raw)
	days := int(math.Floor(leaf.NotAfter.Sub(now).Hours() / 24.0))
	names := append([]string(nil), leaf.DNSNames...)
	if names == nil {
		names = []string{}
	}

	return &SSLCertInfo{
		Subject:           leaf.Subject.CommonName,
		Issuer:            leaf.Issuer.String(),
		DNSNames:          names,
		NotBefore:         leaf.NotBefore.UTC(),
		NotAfter:          leaf.NotAfter.UTC(),
		DaysRemaining:     days,
		ChainCount:        len(certs),
		SHA256Fingerprint: hex.EncodeToString(sum[:]),
	}, nil
}
