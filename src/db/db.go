package db

import (
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
		return err
	}

	if !DB.HasTable(&Users{}) {
		DB.CreateTable(&Users{})
	}

	return nil
}

func FindUserByName(userName string) Users {
	var user Users

	DB.Where("username=?", userName).First(&user)

	return user
}

func CheckUserByUserName(userName string) bool {
	var user Users

	user = FindUserByName(userName)
	if !reflect.DeepEqual(user, Users{}) {
		return true
	}
	return false
}

func AddUser(user *Users) {

	//DB.NewRecord(user)

	DB.Create(&user)
	//
	//DB.Save(&user)

}
