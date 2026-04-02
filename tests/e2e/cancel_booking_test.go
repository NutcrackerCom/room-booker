package e2e

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestCancelBookingFlow(t *testing.T) {
	adminToken := getToken(t, "admin")
	userToken := getToken(t, "user")

	roomID := createRoom(t, adminToken)
	createSchedule(t, adminToken, roomID)

	date := nextWeekdayDate()
	slotID := getFirstAvailableSlot(t, userToken, roomID, date)

	booking := createBooking(t, userToken, slotID)
	bookingID := booking.Booking.ID
	if bookingID == "" {
		t.Fatal("empty booking id")
	}

	cancelResp := cancelBooking(t, userToken, bookingID)
	if cancelResp.Booking.Status != "cancelled" {
		t.Fatalf("expected cancelled status after first cancel, got %q", cancelResp.Booking.Status)
	}

	cancelResp2 := cancelBooking(t, userToken, bookingID)
	if cancelResp2.Booking.Status != "cancelled" {
		t.Fatalf("expected cancelled status after second cancel, got %q", cancelResp2.Booking.Status)
	}

	slotIDAgain := getFirstAvailableSlot(t, userToken, roomID, date)
	if slotIDAgain != slotID {
		t.Fatalf("expected slot %q to become available again, got %q", slotID, slotIDAgain)
	}
}

func cancelBooking(t *testing.T, userToken, bookingID string) bookingResponse {
	t.Helper()

	url := fmt.Sprintf("%s/bookings/%s/cancel", baseURL, bookingID)
	status, body := doJSON(t, http.MethodPost, url, userToken, nil)

	if status != http.StatusOK {
		t.Fatalf("cancel booking failed: status=%d body=%s", status, string(body))
	}

	var resp bookingResponse
	mustUnmarshal(t, body, &resp)
	return resp
}

func uniqueRoomName(prefix string) string {
	return fmt.Sprintf("%s %d", prefix, time.Now().UnixNano())
}

