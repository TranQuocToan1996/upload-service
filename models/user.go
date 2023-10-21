package models

import "upload_service/repositories/db"

type User struct {
	db.BaseModel

	UserName      string `gorm:"column:user_name"`
	Password      string `gorm:"column:password"`
	Salt          string `gorm:"column:salt"`
	RevokeTokenAt int64  `gorm:"column:revoke_token_at"`
}

func (User) TableName() string {
	return "users"
}
