package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

const baseURL = "http://localhost:8080"

type tokenResponse struct {
	Token string `json:"token"`
}

type roomResponse struct {
	Room struct {
		ID string `json:"id"`
	} `json:"room"`
}

type slotResponse struct {
	Slots []struct {
		ID string `json:"id"`
	} `json:"slots"`
}

type bookingResponse struct {
	Booking struct {
		ID     string `json:"id"`
		SlotID string `json:"slotId"`
		UserID string `json:"userId"`
		Status string `json:"status"`
	} `json:"booking"`
}

func TestBookingFlow(t *testing.T) {
	adminToken := getToken(t, "admin")
	userToken := getToken(t, "user")

	roomID := createRoom(t, adminToken)
	createSchedule(t, adminToken, roomID)

	date := nextWeekdayDate()
	slotID := getFirstAvailableSlot(t, userToken, roomID, date)

	booking := createBooking(t, userToken, slotID)
	if booking.Booking.Status != "active" {
		t.Fatalf("expected booking status active, got %q", booking.Booking.Status)
	}
	if booking.Booking.SlotID != slotID {
		t.Fatalf("expected slotId %q, got %q", slotID, booking.Booking.SlotID)
	}
}

func getToken(t *testing.T, role string) string {
	t.Helper()

	status, body := doJSON(t, http.MethodPost, baseURL+"/dummyLogin", "", map[string]any{
		"role": role,
	})

	if status != http.StatusOK {
		t.Fatalf("dummyLogin failed: status=%d body=%s", status, string(body))
	}

	var resp tokenResponse
	mustUnmarshal(t, body, &resp)

	if resp.Token == "" {
		t.Fatal("empty token")
	}

	return resp.Token
}

func createRoom(t *testing.T, adminToken string) string {
	t.Helper()

	name := fmt.Sprintf("E2E Room %d", time.Now().UnixNano())

	status, body := doJSON(t, http.MethodPost, baseURL+"/rooms/create", adminToken, map[string]any{
		"name":        name,
		"description": "e2e room",
		"capacity":    4,
	})

	if status != http.StatusCreated {
		t.Fatalf("create room failed: status=%d body=%s", status, string(body))
	}

	var resp roomResponse
	mustUnmarshal(t, body, &resp)

	if resp.Room.ID == "" {
		t.Fatal("empty room id")
	}

	return resp.Room.ID
}

func createSchedule(t *testing.T, adminToken, roomID string) {
	t.Helper()

	status, body := doJSON(
		t,
		http.MethodPost,
		baseURL+"/rooms/"+roomID+"/schedule/create",
		adminToken,
		map[string]any{
			"daysOfWeek": []int{1, 2, 3, 4, 5},
			"startTime":  "09:00",
			"endTime":    "11:00",
		},
	)

	if status != http.StatusCreated {
		t.Fatalf("create schedule failed: status=%d body=%s", status, string(body))
	}
}

func getFirstAvailableSlot(t *testing.T, userToken, roomID, date string) string {
	t.Helper()

	url := fmt.Sprintf("%s/rooms/%s/slots/list?date=%s", baseURL, roomID, date)
	status, body := doJSON(t, http.MethodGet, url, userToken, nil)

	if status != http.StatusOK {
		t.Fatalf("list slots failed: status=%d body=%s", status, string(body))
	}

	var resp slotResponse
	mustUnmarshal(t, body, &resp)

	if len(resp.Slots) == 0 {
		t.Fatalf("no slots returned for date %s", date)
	}

	return resp.Slots[0].ID
}

func createBooking(t *testing.T, userToken, slotID string) bookingResponse {
	t.Helper()

	status, body := doJSON(t, http.MethodPost, baseURL+"/bookings/create", userToken, map[string]any{
		"slotId": slotID,
	})

	if status != http.StatusCreated {
		t.Fatalf("create booking failed: status=%d body=%s", status, string(body))
	}

	var resp bookingResponse
	mustUnmarshal(t, body, &resp)
	return resp
}

func nextWeekdayDate() string {
	now := time.Now().UTC()

	for i := 1; i <= 7; i++ {
		d := now.AddDate(0, 0, i)
		if d.Weekday() >= time.Monday && d.Weekday() <= time.Friday {
			return d.Format("2006-01-02")
		}
	}

	return now.AddDate(0, 0, 1).Format("2006-01-02")
}

func doJSON(t *testing.T, method, url, token string, payload any) (int, []byte) {
	t.Helper()

	var bodyReader io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("marshal payload: %v", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}

	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}

	return resp.StatusCode, body
}

func mustUnmarshal(t *testing.T, data []byte, target any) {
	t.Helper()

	if err := json.Unmarshal(data, target); err != nil {
		t.Fatalf("unmarshal response: %v; body=%s", err, string(data))
	}
}
