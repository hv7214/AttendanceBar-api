package controllers

import (
	"api/models"
	u "api/utils"
	"encoding/json"
	"net/http"
)

func GetCourse(w http.ResponseWriter, r *http.Request) {

	course := &models.Course{}

	err := json.NewDecoder(r.Body).Decode(course)

	if err != nil {
		u.Response(w, u.Message(false, "Invalid request"))
		return
	}

	response := course.GetCourse()
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
