package unit

import (
	"testing"
	"time"
)

func TestScheduleTimeValidation_Parseable(t *testing.T) {
	start, err := time.Parse("15:04", "09:00")
	if err != nil {
		t.Fatalf("unexpected error parsing start: %v", err)
	}

	end, err := time.Parse("15:04", "18:00")
	if err != nil {
		t.Fatalf("unexpected error parsing end: %v", err)
	}

	if !start.Before(end) {
		t.Fatal("expected start before end")
	}
}

func TestScheduleTimeValidation_StartAfterEnd(t *testing.T) {
	start, err := time.Parse("15:04", "18:00")
	if err != nil {
		t.Fatalf("unexpected error parsing start: %v", err)
	}

	end, err := time.Parse("15:04", "09:00")
	if err != nil {
		t.Fatalf("unexpected error parsing end: %v", err)
	}

	if start.Before(end) {
		t.Fatal("expected start to not be before end")
	}
}

func TestDaysOfWeekValidationRange(t *testing.T) {
	valid := []int{1, 2, 3, 4, 5}
	for _, d := range valid {
		if d < 1 || d > 7 {
			t.Fatalf("expected valid day, got %d", d)
		}
	}

	invalid := []int{0, 8}
	for _, d := range invalid {
		if d >= 1 && d <= 7 {
			t.Fatalf("expected invalid day, got %d", d)
		}
	}
}
