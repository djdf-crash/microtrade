package db

import "github.com/jinzhu/gorm"

type Token struct {
	gorm.Model
	UserID uint `gorm:"index"`
	Token  string
	Expire int64
}

// set User's table name to be `profiles`
func (Token) TableName() string {
	return "tokens"
}

type User struct {
	gorm.Model
	Email    string
	Password string
	Admin    bool
	Tokens   []Token `gorm:"ForeignKey:UserID"`
}

// set User's table name to be `profiles`
func (User) TableName() string {
	return "users"
}
