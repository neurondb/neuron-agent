/*-------------------------------------------------------------------------
 *
 * url_test.go
 *    Tests for SSRF-safe URL validation.
 *
 *-------------------------------------------------------------------------
 */

package validation

import (
	"net"
	"testing"
)

func TestValidateURLForSSRF_InvalidURL(t *testing.T) {
	if err := ValidateURLForSSRF("://bad", nil); err == nil {
		t.Error("expected error for invalid URL")
	}
	if err := ValidateURLForSSRF("", nil); err == nil {
		t.Error("expected error for empty URL")
	}
}

func TestValidateURLForSSRF_SchemeNotAllowed(t *testing.T) {
	if err := ValidateURLForSSRF("ftp://example.com", nil); err == nil {
		t.Error("expected error for ftp scheme")
	}
	if err := ValidateURLForSSRF("file:///etc/passwd", nil); err == nil {
		t.Error("expected error for file scheme")
	}
}

func TestValidateURLForSSRF_PrivateIPRejected(t *testing.T) {
	for _, u := range []string{
		"http://127.0.0.1/",
		"http://10.0.0.1/",
		"http://192.168.1.1/",
		"http://172.16.0.1/",
	} {
		if err := ValidateURLForSSRF(u, nil); err == nil {
			t.Errorf("expected error for private URL %q", u)
		}
	}
}

func TestValidateURLForSSRF_PublicIPAllowed(t *testing.T) {
	if err := ValidateURLForSSRF("http://8.8.8.8/", nil); err != nil {
		t.Errorf("public IP should be allowed: %v", err)
	}
}

func TestValidateURLForSSRF_NoHost(t *testing.T) {
	if err := ValidateURLForSSRF("http://", nil); err == nil {
		t.Error("expected error for URL with no host")
	}
}

func TestIsPrivateOrInternal(t *testing.T) {
	if !isPrivateOrInternal(net.ParseIP("127.0.0.1")) {
		t.Error("127.0.0.1 should be private")
	}
	if !isPrivateOrInternal(net.ParseIP("10.0.0.1")) {
		t.Error("10.0.0.1 should be private")
	}
	if isPrivateOrInternal(net.ParseIP("8.8.8.8")) {
		t.Error("8.8.8.8 should not be private")
	}
}
