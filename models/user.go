package models

import (
	u "api/utils"
	"math"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserCourse struct {
	gorm.Model
	Email                 string `json:"email"`
	CourseName            string `json:"coursename"`
	CourseCode            string `json:"coursecode"`
	TotalClasses          string `json:"totalclasses"`
	Absent                string `json:"absent"`
	TotalClassesConducted string `json:"totalclassesconducted"`
}

type UserCourseKey struct {
	CourseCode  string `json:"coursecode"`
	CourseToken string `json:"coursetoken"`
	Email       string `json:"email"`
}

type UserScheduleDB struct {
	gorm.Model
	Email      string
	Day        string
	CourseCode string
	Time       string
}

type UserScheduleFetch struct {
	Email string `json:"email"`
	Day   string `json:"day"`
}

func (user *User) UserValidate() map[string]interface{} {

	if user.Name == "" {
		return u.Message(false, "Invalid Name")
	}

	if len(user.Password) < 6 {
		return u.Message(false, "Password length critical")
	}

	if user.Email == "" || !strings.Contains(user.Email, "@") {
		return u.Message(false, "Invalid Email")
	}

	existingUser := &User{}
	err := GetDB().Table("users").Where("email=?", user.Email).First(existingUser).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry")
	}

	err = GetDB().Table("users").Where("name=?", user.Name).First(existingUser).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry")
	}

	if existingUser.Email != "" {
		return u.Message(false, "Existing name or email")
	}

	return u.Message(true, "Validation success")
}

func (user *User) UserSignup() map[string]interface{} {

	if response := user.UserValidate(); response["status"] == false {
		return response
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	GetDB().Create(user)

	if user.ID <= 0 {
		return u.Message(false, "Failed to create account, connection error.")
	}

	user.Password = ""

	response := u.Message(true, "Account has been created")
	response["user"] = user
	return response
}

func UserLogin(email string, password string) map[string]interface{} {

	user := &User{}
	err := GetDB().Table("users").Where("email = ?", email).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error, Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return u.Message(false, "Invalid login credentials.")
	}

	user.Password = ""

	response := u.Message(true, "Logged In")
	response["user"] = user
	return response
}

func CheckDistance(longUser float64, latUser float64, latAdmin float64, longAdmin float64) bool {

	R := 6371000.00
	toRadian := 0.0174533

	phi1 := float64(latUser) * toRadian
	phi2 := float64(latAdmin) * toRadian
	delPhi := math.Abs(phi2 - phi1)

	delLembda := float64(longAdmin-longUser) * toRadian

	a := math.Sin(delPhi/2)*math.Sin(delPhi/2) + math.Cos(phi1)*math.Cos(phi2)*math.Sin(delLembda/2)*math.Sin(delLembda/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	d := R * c

	if d <= 10 {
		return true
	}
	return false
}

func (location *Location) MarkAttendance() map[string]interface{} {

	coursecode := location.CourseCode
	adminLong, _ := strconv.ParseFloat(adminLocation.Longitude, 64)
	adminLat, _ := strconv.ParseFloat(adminLocation.Latitude, 64)
	userLong, _ := strconv.ParseFloat(location.Longitude, 64)
	userLat, _ := strconv.ParseFloat(location.Latitude, 64)

	usercourse := &UserCourse{}
	err := GetDB().Table("user_courses").Where("course_code = ?", coursecode).First(usercourse).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Course not found!")
		}
		return u.Message(false, "Connection error, Please retry")
	}

	currentclasses, _ := strconv.Atoi(usercourse.TotalClassesConducted)
	GetDB().Table("user_courses").Where("course_code = ?", coursecode).Update("total_classes_conducted", strconv.Itoa(currentclasses+1))

	if !CheckDistance(userLong, userLat, adminLat, adminLong) {
		absent, _ := strconv.Atoi(usercourse.Absent)
		GetDB().Table("user_courses").Where("course_code = ?", coursecode).Update("absent", strconv.Itoa(absent+1))
		return u.Message(false, "Attendance not marked")
	}

	adminLocation.UpdateAttendance()
	return u.Message(true, "Attendance marked")

}

func (userCourseKey *UserCourseKey) VerifyCourseToken() map[string]interface{} {

	courseCode := userCourseKey.CourseCode
	courseToken := userCourseKey.CourseToken

	course1 := &Course{}
	course2 := &Course{}

	err := GetDB().Table("courses").Where("course_code = ?", courseCode).First(course1).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Course not found")
		}
		return u.Message(false, "Connection error, Please retry")
	}

	err2 := GetDB().Table("courses").Where("course_token = ?", courseToken).First(course2).Error
	if err2 != nil {
		if err2 == gorm.ErrRecordNotFound {
			return u.Message(false, "Course Token not found")
		}
		return u.Message(false, "Connection error, Please retry")
	}

	return u.Message(true, "Verified Course token")
}

func (userCourseKey *UserCourseKey) CheckIfAdded() bool {

	tk := GetDB().Table("user_courses").Where("email = ?", userCourseKey.Email)

	userCourse := &UserCourse{}
	err := tk.Where("course_code = ?", userCourseKey.CourseCode).First(userCourse).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false
		}
	}

	return true
}
func (userCourseKey *UserCourseKey) AddUserCourse() map[string]interface{} {

	checkResp := userCourseKey.VerifyCourseToken()

	if checkResp["status"] == false {
		return checkResp
	}

	if userCourseKey.CheckIfAdded() {
		return u.Message(false, "you are already enrolled in this course")
	}

	usercourse := &UserCourse{}
	usercourse.Email = userCourseKey.Email
	usercourse.CourseCode = userCourseKey.CourseCode
	usercourse.Absent = "0"
	usercourse.TotalClassesConducted = "0"

	course := &Course{}
	err := GetDB().Table("courses").Where("course_code = ?", userCourseKey.CourseCode).First(course).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Course not found")
		}
		return u.Message(false, "Connection error, Please retry")
	}

	usercourse.TotalClasses = course.TotalClasses
	usercourse.CourseName = course.CourseName

	GetDB().Table("user_courses").Create(usercourse)

	currentEnrolled, err := strconv.Atoi(course.EnrolledStudents)
	if err != nil {
		return u.Message(false, "Error in converting string to int")
	}

	newEnrolledstudents := strconv.Itoa(currentEnrolled + 1)

	GetDB().Table("courses").Where("course_code = ?", userCourseKey.CourseCode).Update("enrolled_students", newEnrolledstudents)

	return u.Message(true, "Course Added")
}

func (user *User) UserGetCourses() []UserCourse {

	userEmail := user.Email

	var courses []UserCourse
	GetDB().Table("user_courses").Where("email = ?", userEmail).Find(&courses)

	return courses
}

func (schedule *Schedule) UserAddSchedule() {

	scheduledb := &UserScheduleDB{}
	scheduledb2 := &UserScheduleDB{}
	scheduledb3 := &UserScheduleDB{}
	scheduledb4 := &UserScheduleDB{}
	scheduledb5 := &UserScheduleDB{}

	for _, day := range schedule.Mon {
		usercourses := []UserCourse{}
		GetDB().Table("user_courses").Where("course_code = ?", day.CourseCode).Find(&usercourses)
		for _, user := range usercourses {
			scheduledb.Email = user.Email
			scheduledb.CourseCode = day.CourseCode
			scheduledb.Time = day.Time
			scheduledb.Day = "monday"
			GetDB().Table("user_schedule_dbs").Create(scheduledb)
		}
	}

	for _, day := range schedule.Tue {
		usercourses := []UserCourse{}
		GetDB().Table("user_courses").Where("course_code = ?", day.CourseCode).Find(&usercourses)
		for _, user := range usercourses {
			scheduledb2.Email = user.Email
			scheduledb2.CourseCode = day.CourseCode
			scheduledb2.Time = day.Time
			scheduledb2.Day = "tuesday"
			GetDB().Table("user_schedule_dbs").Create(scheduledb2)
		}
	}

	for _, day := range schedule.Wedn {
		usercourses := []UserCourse{}
		GetDB().Table("user_courses").Where("course_code = ?", day.CourseCode).Find(&usercourses)
		for _, user := range usercourses {
			scheduledb3.Email = user.Email
			scheduledb3.CourseCode = day.CourseCode
			scheduledb3.Time = day.Time
			scheduledb3.Day = "wednesday"
			GetDB().Table("user_schedule_dbs").Create(scheduledb3)
		}
	}

	for _, day := range schedule.Thurs {
		usercourses := []UserCourse{}
		GetDB().Table("user_courses").Where("course_code = ?", day.CourseCode).Find(&usercourses)
		for _, user := range usercourses {
			scheduledb4.Email = user.Email
			scheduledb4.CourseCode = day.CourseCode
			scheduledb4.Time = day.Time
			scheduledb4.Day = "thursday"
			GetDB().Table("user_schedule_dbs").Create(scheduledb4)
		}
	}

	for _, day := range schedule.Fri {

		usercourses := []UserCourse{}
		GetDB().Table("user_courses").Where("course_code = ?", day.CourseCode).Find(&usercourses)
		for _, user := range usercourses {
			scheduledb5.Email = user.Email
			scheduledb5.CourseCode = day.CourseCode
			scheduledb5.Time = day.Time
			scheduledb5.Day = "friday"
			GetDB().Table("user_schedule_dbs").Create(scheduledb5)
		}
	}

}

func (user *UserScheduleFetch) GetUserCurrentSchedule() map[string]interface{} {

	email := user.Email
	day := user.Day

	var courses []Day
	var scheduledbs []UserScheduleDB

	tk := GetDB().Table("user_schedule_dbs").Where("email = ?", email)
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
