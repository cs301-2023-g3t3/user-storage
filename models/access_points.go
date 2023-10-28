package models

type AccessPoint struct {
    Id        int       `json:"id"`
    Name      string    `json:"name" validate:"required"`
	EndPoint  string	`json:"endpoint" gorm:"column:endpoint" validate:"required"`
}

func (AccessPoint) TableName() string {
    return "access_points"
}