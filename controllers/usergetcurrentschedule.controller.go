package controllers

import (
	"api/models"
	u "api/utils"
	"encoding/json"
	"net/http"
)

func GetUserCurrentSchedule(w http.ResponseWriter, r *http.Request) {

	user := &models.UserScheduleFetch{}

	err := json.NewDecoder(r.Body).Decode(user)

	if err != nil {
		u.Response(w, u.Message(false, "Invalid request"))
		return
	}

	response := user.GetUserCurrentSchedule()
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
