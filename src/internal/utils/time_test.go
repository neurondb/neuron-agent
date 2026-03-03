/*-------------------------------------------------------------------------
 *
 * time_test.go
 *    Tests for time formatting and parsing utilities.
 *
 *-------------------------------------------------------------------------
 */

package utils

import (
	"testing"
	"time"
)

func TestFormatTime(t *testing.T) {
	ts := time.Date(2025, 3, 3, 12, 0, 0, 0, time.UTC)
	s := FormatTime(ts)
	if s == "" {
		t.Error("FormatTime returned empty")
	}
	parsed, err := time.Parse(ISO8601Format, s)
	if err != nil {
		t.Errorf("FormatTime output %q not parseable: %v", s, err)
	}
	if !parsed.Equal(ts) {
		t.Errorf("FormatTime roundtrip: got %v", parsed)
	}
}

func TestParseTime(t *testing.T) {
	want := time.Date(2025, 3, 3, 12, 0, 0, 0, time.UTC)
	s := want.Format(ISO8601Format)
	got, err := ParseTime(s)
	if err != nil {
		t.Fatalf("ParseTime(%q) = %v", s, err)
	}
	if !got.Equal(want) {
		t.Errorf("ParseTime(%q) = %v, want %v", s, got, want)
	}
	_, err = ParseTime("not-a-date")
	if err == nil {
		t.Error("ParseTime(invalid) expected error")
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		d    time.Duration
		want string
	}{
		{500 * time.Millisecond, "500ms"},
		{2 * time.Second, "2.00s"},
		{90 * time.Second, "1.50m"},
		{2 * time.Hour, "2.00h"},
	}
	for _, tt := range tests {
		got := FormatDuration(tt.d)
		if got != tt.want {
			t.Errorf("FormatDuration(%v) = %q, want %q", tt.d, got, tt.want)
		}
	}
}

func TestParseDuration(t *testing.T) {
	d, err := ParseDuration("5s")
	if err != nil {
		t.Fatalf("ParseDuration(5s) = %v", err)
	}
	if d != 5*time.Second {
		t.Errorf("ParseDuration(5s) = %v", d)
	}
	_, err = ParseDuration("invalid")
	if err == nil {
		t.Error("ParseDuration(invalid) expected error")
	}
}

func TestUnixTimestamp(t *testing.T) {
	ts := time.Unix(1000000, 0).UTC()
	got := UnixTimestamp(ts)
	if got != 1000000 {
		t.Errorf("UnixTimestamp = %d, want 1000000", got)
	}
}

func TestFromUnixTimestamp(t *testing.T) {
	want := time.Unix(1000000, 0).UTC()
	got := FromUnixTimestamp(1000000)
	if !got.Equal(want) {
		t.Errorf("FromUnixTimestamp(1000000) = %v, want %v", got, want)
	}
}

func TestTimeAgo(t *testing.T) {
	now := time.Now()
	if got := TimeAgo(now.Add(-30 * time.Second)); got != "just now" {
		t.Errorf("TimeAgo(30s ago) = %q", got)
	}
	if got := TimeAgo(now.Add(-2 * time.Minute)); got != "2 minutes ago" {
		t.Errorf("TimeAgo(2m ago) = %q", got)
	}
	if got := TimeAgo(now.Add(-1 * time.Hour)); got != "1 hour ago" {
		t.Errorf("TimeAgo(1h ago) = %q", got)
	}
}
