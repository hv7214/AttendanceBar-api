package models

import (
	u "api/utils"
	"fmt"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type Admin struct {
	gorm.Model
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Course struct {
	gorm.Model
	AdminEmail            string `json:"email"`
	CourseName            string `json:"coursename"`
	CourseCode            string `json:"coursecode"`
	EnrolledStudents      string `json:"enrolledstudents"`
	Attendance            string `json:"attendance"`
	TotalClasses          string `json:"totalclasses"`
	TotalClassesConducted string `json:"totalclassesconducted"`
	CourseToken           string `json:"CourseToken"`
}

type Resource struct {
	gorm.Model
	CourseName string
	CourseCode string
	Filename   string
	Filesize   string
}

type Location struct {
	Email      string `json:"email"`
	CourseCode string `json:"coursecode"`
	Longitude  string `json:"longitude"`
	Latitude   string `json:"latitude"`
}

type Day struct {
	CourseCode string `json:"coursecode"`
	Time       string `json:"time"`
}
type Schedule struct {
	AdminEmail string `json:"email"`
	Mon        []Day  `json:"monday"`
	Tue        []Day  `json:"tuesday"`
	Wedn       []Day  `json:"wednesday"`
	Thurs      []Day  `json:"thursday"`
	Fri        []Day  `json:"friday"`
}

type ScheduleDB struct {
	gorm.Model
	AdminEmail string
	Day        string
	CourseCode string
	Time       string
}

type AdminScheduleFetch struct {
	AdminEmail string `json:"email"`
	Day        string `json:"day"`
}

var adminLocation Location //global admins location

func (admin *Admin) AdminValidate() map[string]interface{} {

	if admin.Name == "" {
		return u.Message(false, "Invalid Name")
	}

	if len(admin.Password) < 6 {
		return u.Message(false, "Password length critical")
	}

	if admin.Email == "" || !strings.Contains(admin.Email, "@") {
		return u.Message(false, "Invalid Email")
	}

	existingAdmin := &Admin{}
	err := GetDB().Table("admins").Where("email=?", admin.Email).First(existingAdmin).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry")
	}

	err = GetDB().Table("admins").Where("name=?", admin.Name).First(existingAdmin).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry")
	}

	if existingAdmin.Email != "" {
		return u.Message(false, "Existing name or email")
	}

	return u.Message(true, "Validation success")
}

func (admin *Admin) AdminSignup() map[string]interface{} {

	if response := admin.AdminValidate(); response["status"] == false {
		return response
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	admin.Password = string(hashedPassword)

	GetDB().Create(admin)

	if admin.ID <= 0 {
		return u.Message(false, "Failed to create account, connection error.")
	}

	admin.Password = ""

	response := u.Message(true, "Account has been created")
	return response
}

func AdminLogin(email string, password string) map[string]interface{} {

	admin := &Admin{}
	err := GetDB().Table("admins").Where("email = ?", email).First(admin).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error, Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return u.Message(false, "Invalid login credentials.")
	}

	admin.Password = ""

	response := u.Message(true, "Logged In")
	return response
}

func (course *Course) AddCourse() map[string]interface{} {

	GetDB().Create(course)

	response := u.Message(true, "Course has been created")

	return response
}

func (admin *Admin) GetCourses() []Course {

	adminEmail := admin.Email

	var courses []Course
	GetDB().Table("courses").Where("admin_email = ?", adminEmail).Find(&courses)

	return courses
}

func (course *Course) GetCourse() map[string]interface{} {

	coursetemp := &Course{}
	err := GetDB().Where("course_code = ?", course.CourseCode).First(coursetemp).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Course not found")
		}
		return u.Message(false, "Connection error, Please retry")
	}

	return map[string]interface{}{
		"course": coursetemp,
	}

}

func AddResource(coursename string, coursecode string, filename string, filesize string) map[string]interface{} {

	resource := &Resource{
		CourseName: coursename,
		CourseCode: coursecode,
		Filename:   filename,
		Filesize:   filesize,
	}

	GetDB().Create(resource)

	return u.Message(true, "Uploaded")
}

func (location *Location) SetLocation() {

	adminLocation = *location
	GetDB().Table("courses").Where("course_code= ?", location.CourseCode).Update("attendance", 0)
}

func (location *Location) UpdateAttendance() {

	course := &Course{}

	err := GetDB().Table("courses").Where("course_code= ?", adminLocation.CourseCode).First(course)
	if err != nil {
		fmt.Println(err)
	}

	currentattendance, _ := strconv.Atoi(course.Attendance)
	currentattendance++

	GetDB().Table("courses").Where("course_code= ?", adminLocation.CourseCode).Update("attendance", currentattendance)
}

func (schedule *Schedule) AddSchedule() map[string]interface{} {

	scheduledb := &ScheduleDB{}
	scheduledb2 := &ScheduleDB{}
	scheduledb3 := &ScheduleDB{}
	scheduledb4 := &ScheduleDB{}
	scheduledb5 := &ScheduleDB{}

	for _, day := range schedule.Mon {
		scheduledb.AdminEmail = schedule.AdminEmail
		scheduledb.CourseCode = day.CourseCode
		scheduledb.Time = day.Time
		scheduledb.Day = "monday"
		GetDB().Table("schedule_dbs").Create(scheduledb)
	}

	for _, day := range schedule.Tue {
		scheduledb2.AdminEmail = schedule.AdminEmail
		scheduledb2.CourseCode = day.CourseCode
		scheduledb2.Time = day.Time
		scheduledb2.Day = "tuesday"
		GetDB().Table("schedule_dbs").Create(scheduledb2)
	}
	for _, day := range schedule.Wedn {
		scheduledb3.AdminEmail = schedule.AdminEmail
		scheduledb3.CourseCode = day.CourseCode
		scheduledb3.Time = day.Time
		scheduledb3.Day = "wednesday"
		GetDB().Table("schedule_dbs").Create(scheduledb3)
	}
	for _, day := range schedule.Thurs {
		scheduledb4.AdminEmail = schedule.AdminEmail
		scheduledb4.CourseCode = day.CourseCode
		scheduledb4.Time = day.Time
		scheduledb4.Day = "thursday"
		GetDB().Table("schedule_dbs").Create(scheduledb4)
	}
	for _, day := range schedule.Fri {
		scheduledb5.AdminEmail = schedule.AdminEmail
		scheduledb5.CourseCode = day.CourseCode
		scheduledb5.Time = day.Time
		scheduledb5.Day = "friday"
		GetDB().Table("schedule_dbs").Create(scheduledb5)
	}

	schedule.UserAddSchedule()

	return u.Message(true, "Schedule uploaded")
}

func (admin *AdminScheduleFetch) GetCurrentSchedule() map[string]interface{} {

	email := admin.AdminEmail
	day := admin.Day

	var courses []Day
	var scheduledbs []ScheduleDB

	tk := GetDB().Table("schedule_dbs").Where("admin_email = ?", email)
	err := tk.Where("day = ?", day).Find(&scheduledbs).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "No class today")
		}
		return u.Message(false, "Connection error, Please retry")
	}

	for _, x := range scheduledbs {
		courses = append(courses, Day{CourseCode: x.CourseCode, Time: x.Time})
	}

	response := make(map[string]interface{})
	response["courses"] = courses
	return response
}
