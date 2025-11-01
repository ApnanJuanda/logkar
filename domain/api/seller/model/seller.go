package model

import "time"

type (
	Seller struct {
		ID                string     `json:"id"`
		Name              string     `json:"name"`
		Email             string     `json:"email"`
		EncryptedPassword string     `json:"encrypted_password"`
		CreatedAt         time.Time  `json:"created_at"`
		UpdatedAt         *time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:false;default:null"`
		DeletedAt         *time.Time `json:"deleted_at"`
		CreatedBy         string     `json:"created_by"`
		UpdatedBy         string     `json:"updated_by"`
		DeletedBy         string     `json:"deleted_by"`
	}

	RegisterSellerRequest struct {
		Name                 string `json:"name" binding:"required"`
		Email                string `json:"email" binding:"required"`
		Password             string `json:"password" binding:"required"`
		PasswordConfirmation string `json:"password_confirmation" binding:"required"`
	}

	LoginSellerRequest struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	LoginSellerResponse struct {
		Token string `json:"token"`
	}
)
