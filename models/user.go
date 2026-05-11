package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        int64 `gorm:"primarykey"`
	FirstName string
	LastName  string
	Email     string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	IsOwner   bool
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (u *User) SetPassword(password string) error {
	hashedPass, err := hashPassword(password)
	u.Password = hashedPass
	return err
}

func CreateDemoAdmin(db *gorm.DB) {
	newUser := User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "demo@useseesaw.dev",
		IsOwner:   true,
	}
	newUser.SetPassword("Demo@1234")
	db.FirstOrCreate(&newUser, User{Email: newUser.Email})
}

func (u *User) CheckPassword(password string) bool {
	return checkPasswordHash(password, u.Password)
}
