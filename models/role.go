package models

type Role struct {
    Id        int       `json:"id"`
    Name      string    `json:"name" validate:"required"`
}
