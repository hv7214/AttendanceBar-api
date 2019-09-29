package controllers

import (
	"api/models"
	u "api/utils"
	"encoding/json"
	"net/http"
)

func GetCurrentSchedule(w http.ResponseWriter, r *http.Request) {

	admin := &models.AdminScheduleFetch{}

	err := json.NewDecoder(r.Body).Decode(admin)

	if err != nil {
		u.Response(w, u.Message(false, "Invalid request"))
		return
	}

	response := admin.GetCurrentSchedule()
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
