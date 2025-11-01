package model

import "time"

type (
	Transaction struct {
		ID          string     `json:"id"`
		CustomerID  string     `json:"customer_id"`
		TotalAmount float64    `json:"total_amount"`
		Status      int        `json:"status"`
		Note        string     `json:"note"`
		CreatedAt   time.Time  `json:"created_at"`
		UpdatedAt   *time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:false;default:null"`
		DeletedAt   *time.Time `json:"deleted_at"`
		CreatedBy   string     `json:"created_by"`
		UpdatedBy   string     `json:"updated_by"`
		DeletedBy   string     `json:"deleted_by"`
	}

	TransactionItem struct {
		ID            int64      `json:"id"`
		TransactionID string     `json:"transaction_id"`
		ProductID     string     `json:"product_id"`
		SizeID        string     `json:"size_id"`
		FlavorID      string     `json:"flavor_id"`
		Quantity      int        `json:"quantity"`
		Subtotal      float64    `json:"subtotal"`
		CreatedAt     time.Time  `json:"created_at"`
		UpdatedAt     *time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:false;default:null"`
		DeletedAt     *time.Time `json:"deleted_at"`
		CreatedBy     string     `json:"created_by"`
		UpdatedBy     string     `json:"updated_by"`
		DeletedBy     string     `json:"deleted_by"`
	}

	TransactionRequest struct {
		CustomerID    string                   `json:"-"`
		CustomerEmail string                   `json:"-"`
		ListItem      []TransactionItemRequest `json:"list_item" binding:"required"`
	}

	TransactionItemRequest struct {
		ProductID string `json:"product_id" binding:"required"`
		SizeID    string `json:"size_id" binding:"required"`
		FlavorID  string `json:"flavor_id" binding:"required"`
		Quantity  int    `json:"quantity" binding:"required"`
	}

	GetTransactionRequest struct {
		CustomerID string `json:"customer_id"`
		Page       int    `json:"page"`
		Limit      int    `json:"limit"`
		Offset     int    `json:"offset"`
	}

	GetTransactionResponse struct {
		ID                  string               `json:"id"`
		CustomerName        string               `json:"customer_name"`
		ListTransactionItem []GetTransactionItem `json:"list_transaction_item"`
	}

	GetTransactionItem struct {
		TransactionId string    `json:"-"`
		CustomerName  string    `json:"-"`
		ProductName   string    `json:"product_name"`
		ProductSize   string    `json:"product_size"`
		ProductFlavor string    `json:"product_flavor"`
		Quantity      int       `json:"quantity"`
		CreatedAt     time.Time `json:"created_at"`
	}
)
