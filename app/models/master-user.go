package models

import "time"

type MasterUser struct {
	ID                 uint64    `json:"id" gorm:"primaryKey;auto_increment;not_null"`
	Username           string    `json:"username" binding:"required" gorm:"uniqueIndex;not null"`
	Email              string    `json:"email" binding:"required" gorm:"uniqueIndex;not null"`
	Password           string    `json:"password" gorm:"not null" binding:"required"`
	Fullname           string    `json:"fullname" gorm:"not null" binding:"required"`
	VerificationCode   string    `json:"verification_code"`
	IsVerified         bool      `json:"is_verified"`
	PasswordResetToken string    `json:"password_reset_token"`
	PasswordResetAt    time.Time `json:"password_reset_at"`
	CreatedBy          uint64    `json:"created_by" `
	CreatedAt          time.Time `json:"created_at"`
	UpdatedBy          uint64    `json:"updated_by"`
	UpdatedAt          time.Time `json:"updated_at"`
	IsDeleted          bool      `json:"is_deleted"`
}

type Pagination struct {
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
	Sort  string `json:"sort"`
}

func (MasterUser) TableName() string {
	return "master_user"
}

type MasterUserUpdateInput struct {
	Username           string    `json:"username" gorm:"uniqueIndex;not null"`
	Email              string    `json:"email" gorm:"uniqueIndex;not null"`
	Password           string    `json:"password" gorm:"not null"`
	Fullname           string    `json:"fullname" gorm:"not null"`
	VerificationCode   string    `json:"verification_code"`
	IsVerified         bool      `json:"is_verified"`
	PasswordResetToken string    `json:"password_reset_token"`
	PasswordResetAt    time.Time `json:"password_reset_at"`
	CreatedBy          uint64    `json:"created_by" `
	CreatedAt          time.Time `json:"created_at"`
	UpdatedBy          uint64    `json:"updated_by"`
	UpdatedAt          time.Time `json:"updated_at"`
	IsDeleted          bool      `json:"is_deleted"`
}

type SignUpInput struct {
	Username        string `json:"username" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
	Fullname        string `json:"fullname"  binding:"required"`
}

type SignInInput struct {
	Username string `json:"username"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}

type ForgotPasswordInput struct {
	Email string `json:"email" binding:"required"`
}

type ResetPasswordInput struct {
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
}
