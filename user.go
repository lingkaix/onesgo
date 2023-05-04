package main

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint64 `gorm:"primaryKey" json:"ID"`
	Email     string `gorm:"unique;not null" json:"email"`
	Password  string `json:"password,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Phone     string `json:"phone,omitempty"`

	gorm.Model
}

func (u *User) BeforeSave(*gorm.DB) error {
	// ? need a better way to filter password to follow the best practice
	if u.Password == "" {
		return nil
	}
	if len(u.Password) < 3 {
		return fmt.Errorf("The password is too short")
	}
	// bcrypt' length limit
	if len(u.Password) > 72 {
		return fmt.Errorf("The password is too long")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// func (u *User) MarshalJSON() ([]byte, error) {
// 	type Alias User
// 	return json.Marshal(&struct {
// 		Password string `json:"password,omitempty"`
// 		*Alias
// 	}{
// 		Password: "",
// 		Alias:    (*Alias)(u),
// 	})
// }

func AddUser(db *gorm.DB, u *User) error {
	return db.Create(u).Error
}

func GetUser(db *gorm.DB, id uint64) (*User, error) {
	var user User
	err := db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByEmail(db *gorm.DB, email string) (*User, error) {
	var user User
	err := db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func UpdateUser(db *gorm.DB, u *User, fieldList []string) error {
	return db.Model(u).Select(fieldList).Updates(u).Error
}

func DeleteUser(db *gorm.DB, id uint64) error {
	return db.Delete(&User{}, id).Error
}
