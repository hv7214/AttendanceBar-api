package controllers

import (
	"api/models"
	u "api/utils"
	"encoding/json"
	"net/http"
)

func AddSchedule(w http.ResponseWriter, r *http.Request) {

	schedule := &models.Schedule{}

	err := json.NewDecoder(r.Body).Decode(schedule)

	if err != nil {
		u.Response(w, u.Message(false, "Invalid request"))
		return
	}

	response := schedule.AddSchedule()
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
