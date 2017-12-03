package db

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Email        string
	Password     string
	Admin        bool
	EmailConfirm bool
	//Tokens   []Token `gorm:"ForeignKey:UserID"`
}

// set User's table name to be `profiles`
func (User) TableName() string {
	return "users"
}
