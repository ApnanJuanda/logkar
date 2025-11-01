package repository

import (
	"bsnack/domain/api/customer/model"
	"gorm.io/gorm"
)

type CustomerRepositoryInterface interface {
	Create(data model.Customer) (err error)
	Take(selectParams []string, conditions interface{}) (resp model.Customer, err error)
	GetCustomerByID(id string) (resp model.Customer, err error)
	UpdateFieldCustomer(tx *gorm.DB, customerID string, req map[string]interface{}) (err error)
	GetListCustomer(req model.GetCustomerRequest) (resp []model.GetCustomerResponse, count int64, err error)
}

type customerRepository struct {
	DB *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepositoryInterface {
	return &customerRepository{
		DB: db,
	}
}

func (r *customerRepository) Create(data model.Customer) (err error) {
	err = r.DB.Create(&data).Error
	return
}

func (r *customerRepository) Take(selectParams []string, conditions interface{}) (resp model.Customer, err error) {
	err = r.DB.Select(selectParams).Take(&resp, conditions).Error
	return
}

func (r *customerRepository) GetCustomerByID(id string) (resp model.Customer, err error) {
	err = r.DB.Table("customers").Where("id = ?", id).Take(&resp).Error
	return
}

func (r *customerRepository) UpdateFieldCustomer(tx *gorm.DB, customerID string, req map[string]interface{}) (err error) {
	err = tx.Model(&model.Customer{}).Where("id = ?", customerID).Updates(req).Error
	return
}

func (r *customerRepository) GetListCustomer(req model.GetCustomerRequest) (resp []model.GetCustomerResponse, count int64, err error) {
	var (
		filter string
		args   []interface{}
	)
	if req.CustomerID != "" {
		dateArgs := []interface{}{
			req.CustomerID,
		}
		args = append(args, dateArgs...)
		newFilter := `(c.id = ?)`
		if filter != "" {
			newFilter = ` AND ` + newFilter
		}
		filter += newFilter
	}

	queryList := r.DB.Debug().Table("customers as c").
		Unscoped().
		Where(filter, args...)

	if err = queryList.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	err = queryList.
		Select(`
		c.name, 
		c.total_points
		`).
		Limit(req.Limit).
		Offset(req.Offset).
		Order(`c.created_at DESC`).Find(&resp).Error

	if err != nil {
		return nil, 0, err
	}
	return
}
