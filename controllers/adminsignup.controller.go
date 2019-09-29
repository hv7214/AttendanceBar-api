package controllers

import (
	"api/models"
	u "api/utils"
	"encoding/json"
	"net/http"
)

func AdminSignup(w http.ResponseWriter, r *http.Request) {

	admin := &models.Admin{}
	err := json.NewDecoder(r.Body).Decode(admin)

	if err != nil {
		u.Response(w, u.Message(false, "Invalid request"))
		return
	}

	response := admin.AdminSignup()
	u.Response(w, response)
}
