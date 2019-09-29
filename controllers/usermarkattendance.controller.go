package controllers

import (
	"api/models"
	u "api/utils"
	"encoding/json"
	"net/http"
)

func MarkAttendance(w http.ResponseWriter, r *http.Request) {

	location := &models.Location{}
	err := json.NewDecoder(r.Body).Decode(location)

	if err != nil {
		u.Response(w, u.Message(false, "Invalid request"))
		return
	}

	reponse := location.MarkAttendance()
	u.Response(w, reponse)
}
