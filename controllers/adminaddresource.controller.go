package controllers

import (
	"api/models"
	u "api/utils"
	"net/http"
)

func AddResource(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(10 << 20)
	r.ParseForm()

	file, handler, err := r.FormFile("resource")
	if err != nil {
		response := u.Message(false, "Error in retrieving file")
		u.Response(w, response)
	}

	defer file.Close()
	// fmt.Println(handler.Filename, handler.Size, handler.Header)

	coursename := r.PostFormValue("coursename")
	coursecode := r.PostFormValue("coursecode")
	response := models.AddResource(coursename, coursecode, handler.Filename, string(handler.Size/1000000))
	// fileBytes, err := ioutil.ReadAll(file)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Println(fileBytes)
	u.Response(w, response)
}
