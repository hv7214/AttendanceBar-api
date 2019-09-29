package controllers

import (
	"api/models"
	u "api/utils"
	"encoding/json"
	"net/http"
)

func AddUserCourse(w http.ResponseWriter, r *http.Request) {

	usercourseKey := &models.UserCourseKey{}

	err := json.NewDecoder(r.Body).Decode(usercourseKey)
	if err != nil {
		u.Response(w, u.Message(false, "Invalid request"))
		return
	}

	response := usercourseKey.AddUserCourse()
	u.Response(w, response)
}
