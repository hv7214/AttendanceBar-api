package controllers

import (
	"api/models"
	u "api/utils"
	"encoding/json"
	"net/http"
)

func PostLocation(w http.ResponseWriter, r *http.Request) {

	adminLocation := &models.Location{}

	err := json.NewDecoder(r.Body).Decode(adminLocation)

	if err != nil {
		u.Response(w, u.Message(false, "Invalid request"))
		return
	}
	adminLocation.SetLocation()
	u.Response(w, u.Message(true, "Attendance marked"))
}
