package workflow

import (
	"testing"
	"time"
)

func TestNextScheduledRun_hourlyUTC(t *testing.T) {
	after := time.Date(2026, 5, 1, 10, 30, 0, 0, time.UTC)
	next, err := NextScheduledRun("0 * * * *", "UTC", after)
	if err != nil {
		t.Fatal(err)
	}
	want := time.Date(2026, 5, 1, 11, 0, 0, 0, time.UTC)
	if !next.Equal(want) {
		t.Fatalf("got %v want %v", next, want)
	}
}

func TestNextScheduledRun_invalidCron(t *testing.T) {
	_, err := NextScheduledRun("not-a-cron", "UTC", time.Now())
	if err == nil {
		t.Fatal("expected error")
	}
}
