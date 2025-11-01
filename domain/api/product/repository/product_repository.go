package repository

import (
	"bsnack/domain/api/product/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
)

type ProductRepositoryInterface interface {
	InsertSize(datas []model.Size) (err error)
	GetAllSize() (resp []model.Size, err error)
	GetSizeByID(tx *gorm.DB, id string) (resp model.Size, err error)
	InsertFlavor(datas []model.Flavor) (err error)
	GetAllFlavor() (resp []model.Flavor, err error)
	GetFlavorByID(tx *gorm.DB, id string) (resp model.Size, err error)
	InsertProductType(data []model.ProductType) (err error)
	GetProductTypeByID(tx *gorm.DB, id string) (resp model.ProductType, err error)
	InsertProduct(tx *gorm.DB, data []model.Product) (err error)
	GetProductByID(tx *gorm.DB, id string) (resp model.Product, err error)
	InsertProductDetail(tx *gorm.DB, datas []model.ProductDetail) (err error)
	GetProductDetailByParams(tx *gorm.DB, params model.ProductDetailParams) (resp model.ProductDetailResponse, err error)
	UpdateFieldProductDetail(tx *gorm.DB, id int64, req map[string]interface{}) (err error)
	GetListProduct(req model.GetProductRequest) (resp []model.GetProductResponse, count int64, err error)
}

type productRepository struct {
	DB *gorm.DB
}

func NewProductRepository(DB *gorm.DB) ProductRepositoryInterface {
	return &productRepository{
		DB: DB,
	}
}

func (r *productRepository) InsertSize(datas []model.Size) (err error) {
	err = r.DB.Create(&datas).Error
	return
}

func (r *productRepository) GetAllSize() (resp []model.Size, err error) {
	err = r.DB.Find(&resp).Error
	return
}

func (r *productRepository) GetSizeByID(tx *gorm.DB, id string) (resp model.Size, err error) {
	if tx == nil {
		tx = r.DB
	}
	err = tx.Table("sizes").Where("id = ?", id).Take(&resp).Error
	return
}

func (r *productRepository) InsertFlavor(datas []model.Flavor) (err error) {
	err = r.DB.Create(&datas).Error
	return
}

func (r *productRepository) GetAllFlavor() (resp []model.Flavor, err error) {
	err = r.DB.Find(&resp).Error
	return
}

func (r *productRepository) GetFlavorByID(tx *gorm.DB, id string) (resp model.Size, err error) {
	if tx == nil {
		tx = r.DB
	}
	err = tx.Table("flavors").Where("id = ?", id).Take(&resp).Error
	return
}

func (r *productRepository) BeginTrans() *gorm.DB {
	return r.DB.Begin()
}

func (r *productRepository) CommitTrans(tx *gorm.DB) error {
	return tx.Commit().Error
}

func (r *productRepository) RollbackTrans(tx *gorm.DB) error {
	return tx.Rollback().Error
}

func (r *productRepository) InsertProductType(data []model.ProductType) (err error) {
	err = r.DB.Create(&data).Error
	return
}

func (r *productRepository) GetProductTypeByID(tx *gorm.DB, id string) (resp model.ProductType, err error) {
	if tx == nil {
		tx = r.DB
	}
	err = tx.Table("product_types").Where("id = ?", id).Take(&resp).Error
	return
}

func (r *productRepository) InsertProduct(tx *gorm.DB, data []model.Product) (err error) {
	err = tx.Create(&data).Error
	return
}

func (r *productRepository) GetProductByID(tx *gorm.DB, id string) (resp model.Product, err error) {
	if tx == nil {
		tx = r.DB
	}
	err = tx.Table("products").Where("id = ?", id).Take(&resp).Error
	return
}

func (r *productRepository) InsertProductDetail(tx *gorm.DB, datas []model.ProductDetail) (err error) {
	err = tx.Create(&datas).Error
	return
}

func (r *productRepository) GetProductDetailByParams(tx *gorm.DB, params model.ProductDetailParams) (resp model.ProductDetailResponse, err error) {
	if tx == nil {
		tx = r.DB
		err = tx.Table("product_details as pd").
			Joins("LEFT JOIN point_redeem_rules as prr ON prr.size_id = pd.size_id").
			Select("pd.*, prr.exchange_point").
			Where("pd.product_id = ? AND pd.size_id = ? AND pd.flavor_id = ?", params.ProductID, params.SizeID, params.FlavorID).Take(&resp).Error
		return
	} else {
		err = tx.Table("product_details as pd").
			Joins("INNER JOIN point_redeem_rules as prr ON prr.size_id = pd.size_id").
			Select("pd.*, prr.exchange_point").
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("pd.product_id = ? AND pd.size_id = ? AND pd.flavor_id = ?", params.ProductID, params.SizeID, params.FlavorID).
			Scan(&resp).Error
		if err != nil {
			log.Printf("error lock table product detail: %v", err.Error())
			return
		}
	}
	return
}

func (r *productRepository) UpdateFieldProductDetail(tx *gorm.DB, id int64, req map[string]interface{}) (err error) {
	err = tx.Model(&model.ProductDetail{}).Where("id = ?", id).Updates(req).Error
	return
}

func (r *productRepository) GetListProduct(req model.GetProductRequest) (resp []model.GetProductResponse, count int64, err error) {
	var (
		filter string
		args   []interface{}
	)
	if req.StartDate != "" && req.EndDate != "" {
		dateArgs := []interface{}{
			req.StartDate,
			req.EndDate,
		}
		args = append(args, dateArgs...)
		newFilter := `(p.created_at >= ? AND p.created_at <= ?)`
		if filter != "" {
			newFilter = ` AND ` + newFilter
		}
		filter += newFilter
	}

	queryList := r.DB.Debug().Table("products as p").
		Joins("LEFT JOIN product_types as pt ON pt.id = p.type_id").
		Joins("LEFT JOIN product_details as pd ON pd.product_id = p.id").
		Joins("LEFT JOIN flavors as f ON f.id = pd.flavor_id").
		Joins("LEFT JOIN sizes as s ON s.id = pd.size_id").
		Unscoped().
		Where(filter, args...)

	if err = queryList.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	err = queryList.
		Select(`p.name as name, 
		pt.name as type,
		f.name as flavor,
		pd.stock as quantity,
		s.name as size,
		pd.price as price 
		`).
		Limit(req.Limit).
		Offset(req.Offset).
		Order(`p.created_at DESC`).Find(&resp).Error

	if err != nil {
		return nil, 0, err
	}
	return
}
