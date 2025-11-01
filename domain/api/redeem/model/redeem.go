package model

import "time"

type (
	CustomerPointRedeem struct {
		ID               int64      `json:"id"`
		CustomerID       string     `json:"customer_id"`
		TotalRedeemPoint int        `json:"total_redeem_point"`
		CreatedAt        time.Time  `json:"created_at"`
		UpdatedAt        *time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:false;default:null"`
		DeletedAt        *time.Time `json:"deleted_at"`
		CreatedBy        string     `json:"created_by"`
		UpdatedBy        string     `json:"updated_by"`
		DeletedBy        string     `json:"deleted_by"`
	}

	CustomerPointReport struct {
		ID                    int64      `json:"id"`
		CustomerID            string     `json:"customer_id"`
		TransactionID         string     `json:"transaction_id"`
		CustomerPointRedeemID int64      `json:"customer_point_redeem_id"`
		Status                string     `json:"status"` // cashback, redeem
		Balance               int        `json:"balance"`
		PointIn               int        `json:"point_in"`
		PointOut              int        `json:"point_out"`
		CreatedAt             time.Time  `json:"created_at"`
		UpdatedAt             *time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:false;default:null"`
		DeletedAt             *time.Time `json:"deleted_at"`
		CreatedBy             string     `json:"created_by"`
		UpdatedBy             string     `json:"updated_by"`
		DeletedBy             string     `json:"deleted_by"`
	}

	CustomerRedeemedProduct struct {
		ID                    int64      `json:"id"`
		CustomerPointRedeemID int64      `json:"customer_point_redeem_id"`
		ProductID             string     `json:"product_id"`
		SizeID                string     `json:"size_id"`
		FlavorID              string     `json:"flavor_id"`
		Quantity              int        `json:"quantity"`
		CreatedAt             time.Time  `json:"created_at"`
		UpdatedAt             *time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:false;default:null"`
		DeletedAt             *time.Time `json:"deleted_at"`
		CreatedBy             string     `json:"created_by"`
		UpdatedBy             string     `json:"updated_by"`
		DeletedBy             string     `json:"deleted_by"`
	}

	RedeemPointRequest struct {
		CustomerID string `json:"-"`
		ProductID  string `json:"product_id"`
		SizeID     string `json:"size_id"`
		FlavorID   string `json:"flavor_id"`
		Quantity   int    `json:"quantity"`
	}
)
