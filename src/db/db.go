package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const (
	pathDB = "./microtrade.db"
)

var DB *gorm.DB

func InitDB() error {

	var err error

	DB, err = gorm.Open("sqlite3", pathDB)
	if err != nil {
		return err
	}
	return nil
}
