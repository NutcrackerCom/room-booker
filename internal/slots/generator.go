package slots

import "time"

type GeneratedSlot struct {
	Start time.Time
	End   time.Time
}

func BuildSlotsForDate(date time.Time, daysOfWeek []int, startHHMM, endHHMM string) ([]GeneratedSlot, error) {
	weekday := int(date.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	allowed := false
	for _, d := range daysOfWeek {
		if d == weekday {
			allowed = true
			break
		}
	}
	if !allowed {
		return nil, nil
	}

	startParsed, err := time.Parse("15:04", startHHMM)
	if err != nil {
		return nil, err
	}
	endParsed, err := time.Parse("15:04", endHHMM)
	if err != nil {
		return nil, err
	}

	current := time.Date(date.Year(), date.Month(), date.Day(), startParsed.Hour(), startParsed.Minute(), 0, 0, time.UTC)
	finish := time.Date(date.Year(), date.Month(), date.Day(), endParsed.Hour(), endParsed.Minute(), 0, 0, time.UTC)

	var result []GeneratedSlot
	for current.Before(finish) {
		next := current.Add(30 * time.Minute)
		if next.After(finish) {
			break
		}

		result = append(result, GeneratedSlot{
			Start: current,
			End:   next,
		})
		current = next
	}

	return result, nil
}
