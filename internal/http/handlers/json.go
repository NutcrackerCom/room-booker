package handlers

import (
	"encoding/json"
	"net/http"
)

func decodeJSON(r *http.Request, dst any) error {
	return json.NewDecoder(r.Body).Decode(dst)
}
