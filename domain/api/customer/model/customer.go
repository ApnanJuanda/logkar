package model

import "time"

type (
	Customer struct {
		ID                string     `json:"id"`
		Name              string     `json:"name"`
		Email             string     `json:"email"`
		EncryptedPassword string     `json:"encrypted_password"`
		Phone             string     `json:"phone"`
		TotalPoints       int        `json:"total_points"`
		CreatedAt         time.Time  `json:"created_at"`
		UpdatedAt         *time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:false;default:null"`
		DeletedAt         *time.Time `json:"deleted_at"`
		CreatedBy         string     `json:"created_by"`
		UpdatedBy         string     `json:"updated_by"`
		DeletedBy         string     `json:"deleted_by"`
	}

	RegisterCustomerRequest struct {
		Name                 string `json:"name" binding:"required"`
		Email                string `json:"email" binding:"required"`
		Phone                string `json:"phone" binding:"required"`
		Password             string `json:"password" binding:"required"`
		PasswordConfirmation string `json:"password_confirmation" binding:"required"`
	}

	LoginCustomerRequest struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	LoginCustomerResponse struct {
		Token string `json:"token"`
	}

	GetCustomerRequest struct {
		CustomerID string
		Page       int
		Limit      int
		Offset     int
	}

	GetCustomerResponse struct {
		Name        string `json:"name"`
		TotalPoints int    `json:"total_points"`
	}
)
