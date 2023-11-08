package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
    Id        string    `json:"id" gorm:"primaryKey;"`
    FirstName string    `json:"firstName" validate:"required"`
    LastName  string    `json:"lastName" validate:"required"`
    Email     string    `json:"email" validate:"required,email"`
    Role      *uint     `json:"role" gorm:"default:null"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (error){
    u.Id = uuid.NewString()
    return nil
}
