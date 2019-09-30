package main

import (
	"api/controllers"
	"fmt"
	"net/http"
	"os"
	"github.com/gorilla/mux"
)

func InitRouter() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/user/login", controllers.UserLogin).Methods("POST")
	router.HandleFunc("/user/signup", controllers.UserSignup).Methods("POST")
	router.HandleFunc("/admin/login", controllers.AdminLogin).Methods("POST")
	router.HandleFunc("/admin/signup", controllers.AdminSignup).Methods("POST")
	router.HandleFunc("/admin/getcourses", controllers.GetCourses).Methods("POST")
	router.HandleFunc("/admin/addcourse", controllers.AddCourse).Methods("POST")
	router.HandleFunc("/admin/getcourse", controllers.GetCourse).Methods("POST")
	router.HandleFunc("/admin/addresource", controllers.AddResource).Methods("POST")
	router.HandleFunc("/admin/postlocation", controllers.PostLocation).Methods("POST")
	router.HandleFunc("/user/markattendance", controllers.MarkAttendance).Methods("POST")
	router.HandleFunc("/user/addcourse", controllers.AddUserCourse).Methods("POST")
	router.HandleFunc("/user/getcourses", controllers.UserGetCourses).Methods("POST")
	router.HandleFunc("/admin/addschedule", controllers.AddSchedule).Methods("POST")
	router.HandleFunc("/admin/gettodayschedule", controllers.GetCurrentSchedule).Methods("POST")
	router.HandleFunc("/user/gettodayschedule", controllers.GetUserCurrentSchedule).Methods("POST")

	return router
}

func main() {

	router := InitRouter()

	port := os.Getenv("PORT")
	err := http.ListenAndServe(":"+port, router)
	if err != nil {

		fmt.Println(err)
	}
}
