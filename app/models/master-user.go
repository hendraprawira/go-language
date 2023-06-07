package models

import "time"

type MasterUser struct {
	ID        uint64    `json:"id" gorm:"primaryKey;auto_increment;not_null"`
	Username  string    `json:"username" binding:"required" gorm:"uniqueIndex;not null"`
	Email     string    `json:"email" binding:"required" gorm:"uniqueIndex;not null"`
	Password  string    `json:"password" binding:"required"`
	Fullname  string    `json:"fullname" binding:"required"`
	CreatedBy uint64    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedBy uint64    `json:"updated_by"`
	UpdatedAt time.Time `json:"updated_at"`
	IsDeleted bool      `json:"is_deleted"`
}

type Pagination struct {
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
	Sort  string `json:"sort"`
}

func (MasterUser) TableName() string {
	return "master_user"
}

type SignUpInput struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
	Photo           string `json:"photo" binding:"required"`
}

type SignInInput struct {
	Username string `json:"username"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}
