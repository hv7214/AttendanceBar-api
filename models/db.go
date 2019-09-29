package models

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

var db *gorm.DB

func init() {

	e := godotenv.Load(".env")
	if e != nil {
		fmt.Println(e)
	}

	dbUri := os.Getenv("DATABASE_URL")
	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		fmt.Println(err)
	}

	db = conn
	db.Debug().AutoMigrate(&User{})
	db.Debug().AutoMigrate(&Admin{})
	db.Debug().AutoMigrate(&Course{})
	db.Debug().AutoMigrate(&Resource{})
	db.Debug().AutoMigrate(&UserCourse{})
	db.Debug().AutoMigrate(&ScheduleDB{})
	db.Debug().AutoMigrate(&UserScheduleDB{})
}

func GetDB() *gorm.DB {
	return db
}
