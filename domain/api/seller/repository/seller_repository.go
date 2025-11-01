package repository

import (
	"bsnack/domain/api/seller/model"
	"gorm.io/gorm"
)

type SellerRepositoryInterface interface {
	Create(data model.Seller) (err error)
	Take(selectParams []string, conditions interface{}) (resp model.Seller, err error)
}

type sellerRepository struct {
	DB *gorm.DB
}

func NewSellerRepository(db *gorm.DB) SellerRepositoryInterface {
	return &sellerRepository{
		DB: db,
	}
}

func (r *sellerRepository) Create(data model.Seller) (err error) {
	err = r.DB.Create(&data).Error
	return
}

func (r *sellerRepository) Take(selectParams []string, conditions interface{}) (resp model.Seller, err error) {
	err = r.DB.Select(selectParams).Take(&resp, conditions).Error
	return
}
