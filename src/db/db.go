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

	DB.SingularTable(true)

	DB.AutoMigrate(&User{})

	return nil
}

func FindUserByName(email string) User {
	var user User

	DB.Where("email=?", email).First(&user)

	return user
}

func CheckUserByEmail(email string) bool {
	var user User

	user = FindUserByName(email)
	if !reflect.DeepEqual(user, User{}) {
		return true
	}
	return false
}

func AddUser(user *User) error {

	tx := DB.Begin()

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil

}

func UpdateUser(user *User) error {

	tx := DB.Begin()

	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil

}
