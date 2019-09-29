package controllers

import (
	"api/models"
	u "api/utils"
	"encoding/json"
	"net/http"
)

func UserSignup(w http.ResponseWriter, r *http.Request) {

	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)

	if err != nil {
		u.Response(w, u.Message(false, "Invalid request"))
		return
	}

	response := user.UserSignup()
	u.Response(w, response)
}
