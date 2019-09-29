package controllers

import (
	"api/models"
	u "api/utils"
	"encoding/json"
	"net/http"
)

func UserGetCourses(w http.ResponseWriter, r *http.Request) {

	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)

	if err != nil {
		u.Response(w, u.Message(false, "Invalid request"))
		return
	}

	courses := user.UserGetCourses()

	// fmt.Println(response)
	response := make(map[string][]models.UserCourse)
	response["courses"] = courses
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
