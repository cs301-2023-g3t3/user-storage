package models

type User struct {
    Id        string    `json:"id"`
    FirstName string    `json:"firstName" validate:"required"`
    LastName  string    `json:"lastName" validate:"required"`
    Email     string    `json:"email" validate:"required,email"`
    Role      int       `json:"role" gorm:"default:null"`
}
