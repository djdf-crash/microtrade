package db

import (
	"log"
	"reflect"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const (
	pathDB = "./microtrade.db3"
)

var DB *gorm.DB

func InitDB() error {

	var err error

	DB, err = gorm.Open("sqlite3", pathDB)
	if err != nil {
		log.Panic(err.Error())
		return err
	}

	if !DB.HasTable(&Users{}) {
		DB.CreateTable(&Users{})
	}

	//DB.AutoMigrate(&Users{})

	return nil
}

func FindUserByName(email string) Users {
	var user Users

	DB.Where("email=?", email).First(&user)

	return user
}

func CheckUserByEmail(email string) bool {
	var user Users

	user = FindUserByName(email)
	if !reflect.DeepEqual(user, Users{}) {
		return true
	}
	return false
}

func AddUser(user *Users) error {

	tx := DB.Begin()

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil

}
