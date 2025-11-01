package model

import "time"

type (
	Size struct {
		ID        string     `json:"id"`
		Name      string     `json:"name"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:false;default:null"`
		DeletedAt *time.Time `json:"deleted_at"`
		CreatedBy string     `json:"created_by"`
		UpdatedBy string     `json:"updated_by"`
		DeletedBy string     `json:"deleted_by"`
	}

	InsertSizeRequest struct {
		SellerEmail string   `json:"-"`
		ListName    []string `json:"list_name"`
	}

	InsertFlavorRequest struct {
		SellerEmail string   `json:"-"`
		ListName    []string `json:"list_name"`
	}

	Flavor struct {
		ID        string     `json:"id"`
		Name      string     `json:"name"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:false;default:null"`
		DeletedAt *time.Time `json:"deleted_at"`
		CreatedBy string     `json:"created_by"`
		UpdatedBy string     `json:"updated_by"`
		DeletedBy string     `json:"deleted_by"`
	}

	InsertProductTypeRequest struct {
		SellerEmail string   `json:"-"`
		ListName    []string `json:"list_name"`
	}

	ProductType struct {
		ID        string     `json:"id"`
		Name      string     `json:"name"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:false;default:null"`
		DeletedAt *time.Time `json:"deleted_at"`
		CreatedBy string     `json:"created_by"`
		UpdatedBy string     `json:"updated_by"`
		DeletedBy string     `json:"deleted_by"`
	}

	Product struct {
		ID        string     `json:"id"`
		Name      string     `json:"name"`
		SellerID  string     `json:"seller_id"`
		TypeID    string     `json:"type_id"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:false;default:null"`
		DeletedAt *time.Time `json:"deleted_at"`
		CreatedBy string     `json:"created_by"`
		UpdatedBy string     `json:"updated_by"`
		DeletedBy string     `json:"deleted_by"`
	}

	ProductDetail struct {
		ID        int64      `json:"id"`
		ProductID string     `json:"product_id"`
		SizeID    string     `json:"size_id"`
		FlavorID  string     `json:"flavor_id"`
		Price     float64    `json:"price"`
		Stock     int        `json:"stock"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:false;default:null"`
		DeletedAt *time.Time `json:"deleted_at"`
		CreatedBy string     `json:"created_by"`
		UpdatedBy string     `json:"updated_by"`
		DeletedBy string     `json:"deleted_by"`
	}

	ProductDetailResponse struct {
		ID            int64      `json:"id"`
		ProductID     string     `json:"product_id"`
		SizeID        string     `json:"size_id"`
		FlavorID      string     `json:"flavor_id"`
		Price         float64    `json:"price"`
		Stock         int        `json:"stock"`
		ExchangePoint int        `json:"exchange_point"`
		CreatedAt     time.Time  `json:"created_at"`
		UpdatedAt     *time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:false;default:null"`
		DeletedAt     *time.Time `json:"deleted_at"`
		CreatedBy     string     `json:"created_by"`
		UpdatedBy     string     `json:"updated_by"`
		DeletedBy     string     `json:"deleted_by"`
	}

	InsertProductRequest struct {
		SellerID    string               `json:"-"`
		SellerEmail string               `json:"-"`
		ListProduct []ProductItemRequest `json:"list_product" binding:"required"`
	}

	ProductItemRequest struct {
		Name   string `json:"name"`
		TypeId string `json:"type_id"`
	}

	InsertProductDetailRequest struct {
		SellerID        string               `json:"-"`
		SellerEmail     string               `json:"-"`
		ListProductInfo []ProductInfoRequest `json:"list_product_info" binding:"required"`
	}

	ProductInfoRequest struct {
		ProductID string  `json:"product_id" binding:"required"`
		SizeID    string  `json:"size_id" binding:"required"`
		FlavorID  string  `json:"flavor_id" binding:"required"`
		Price     float64 `json:"price" binding:"required"`
		Stock     int     `json:"stock" binding:"required"`
	}

	ProductDetailParams struct {
		ProductID string `json:"product_id"`
		SizeID    string `json:"size_id"`
		FlavorID  string `json:"flavor_id"`
	}

	GetProductRequest struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		Page      int
		Limit     int
		Offset    int
	}

	GetProductResponse struct {
		Name     string  `json:"name"`
		Type     string  `json:"type"`
		Flavor   string  `json:"flavor"`
		Quantity string  `json:"quantity"`
		Size     string  `json:"size"`
		Price    float64 `json:"price"`
	}
)
