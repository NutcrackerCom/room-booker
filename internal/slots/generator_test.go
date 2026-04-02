package slots

import (
	"testing"
	"time"
)

func TestBuildSlotsForDate_ReturnsSlotsForAllowedWeekday(t *testing.T) {
	date := time.Date(2026, 4, 3, 0, 0, 0, 0, time.UTC) // Friday

	result, err := BuildSlotsForDate(date, []int{1, 2, 3, 4, 5}, "09:00", "11:00")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 4 {
		t.Fatalf("expected 4 slots, got %d", len(result))
	}

	if result[0].Start.Format(time.RFC3339) != "2026-04-03T09:00:00Z" {
		t.Fatalf("unexpected first slot start: %s", result[0].Start.Format(time.RFC3339))
	}
	if result[3].End.Format(time.RFC3339) != "2026-04-03T11:00:00Z" {
		t.Fatalf("unexpected last slot end: %s", result[3].End.Format(time.RFC3339))
	}
}

func TestBuildSlotsForDate_ReturnsEmptyForDisallowedWeekday(t *testing.T) {
	date := time.Date(2026, 4, 4, 0, 0, 0, 0, time.UTC) // Saturday

	result, err := BuildSlotsForDate(date, []int{1, 2, 3, 4, 5}, "09:00", "11:00")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 0 {
		t.Fatalf("expected 0 slots, got %d", len(result))
	}
}

func TestBuildSlotsForDate_InvalidStartTime(t *testing.T) {
	date := time.Date(2026, 4, 3, 0, 0, 0, 0, time.UTC)

	_, err := BuildSlotsForDate(date, []int{1, 2, 3, 4, 5}, "bad", "11:00")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
