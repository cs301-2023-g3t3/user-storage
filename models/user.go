package models

type User struct {
    Id        string    `json:"id"`
    FirstName string    `json:"firstName"`
    LastName  string    `json:"lastName"`
    Email     string    `json:"email"`
    Role      string    `json:"role"`
}
