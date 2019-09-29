package controllers

import (
	"api/models"
	u "api/utils"
	"encoding/json"
	"net/http"
)

func GetCourses(w http.ResponseWriter, r *http.Request) {

	admin := &models.Admin{}
	err := json.NewDecoder(r.Body).Decode(admin)

	if err != nil {
		u.Response(w, u.Message(false, "Invalid request"))
		return
	}

	courses := admin.GetCourses()

	// fmt.Println(response)
	response := make(map[string][]models.Course)
	response["courses"] = courses
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
