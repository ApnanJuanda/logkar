package repository

import (
	"bsnack/domain/api/redeem/model"
	"gorm.io/gorm"
)

type RedeemRepositoryInterface interface {
	InsertCustomerPointRedeem(tx *gorm.DB, data *model.CustomerPointRedeem) (err error)
	InsertRedeemedProduct(tx *gorm.DB, data model.CustomerRedeemedProduct) (err error)
	InsertCustomerPointReport(tx *gorm.DB, data model.CustomerPointReport) (err error)
}

type redeemRepository struct {
	DB *gorm.DB
}

func NewRedeemRepository(DB *gorm.DB) RedeemRepositoryInterface {
	return &redeemRepository{
		DB: DB,
	}
}

func (r *redeemRepository) InsertCustomerPointRedeem(tx *gorm.DB, data *model.CustomerPointRedeem) (err error) {
	if tx == nil {
		tx = r.DB
	}
	err = tx.Create(data).Error
	return
}

func (r *redeemRepository) InsertRedeemedProduct(tx *gorm.DB, data model.CustomerRedeemedProduct) (err error) {
	if tx == nil {
		tx = r.DB
	}
	err = tx.Create(&data).Error
	return
}

func (r *redeemRepository) InsertCustomerPointReport(tx *gorm.DB, data model.CustomerPointReport) (err error) {
	if tx == nil {
		tx = r.DB
	}
	err = tx.Create(&data).Error
	return
}
