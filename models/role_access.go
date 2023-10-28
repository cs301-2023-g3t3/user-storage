package models

type RoleAccess struct {
    RoleId    int    `json:"roleId" gorm:"column:role_id" validate:"required"`
    APId      int    `json:"apId" gorm:"column:ap_id" validate:"required"`
}

func (RoleAccess) TableName() string {
    return "role_access"
}