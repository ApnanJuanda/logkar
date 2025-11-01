package repository

import (
	"bsnack/domain/api/transaction/model"
	"gorm.io/gorm"
)

type TransactionRepositoryInterface interface {
	CreateTransaction(tx *gorm.DB, data model.Transaction) (err error)
	SaveTransactionItems(tx *gorm.DB, datas []model.TransactionItem) (err error)
	GetListTransaction(req model.GetTransactionRequest) (resp []model.GetTransactionResponse, count int64, err error)
}

type transactionRepository struct {
	DB *gorm.DB
}

func NewTransactionRepository(DB *gorm.DB) TransactionRepositoryInterface {
	return &transactionRepository{
		DB: DB,
	}
}

func (r *transactionRepository) CreateTransaction(tx *gorm.DB, data model.Transaction) (err error) {
	err = tx.Create(&data).Error
	return
}

func (r *transactionRepository) SaveTransactionItems(tx *gorm.DB, datas []model.TransactionItem) (err error) {
	if tx == nil {
		tx = r.DB
	}
	err = tx.Create(&datas).Error
	return
}

func (r *transactionRepository) GetListTransaction(req model.GetTransactionRequest) (resp []model.GetTransactionResponse, count int64, err error) {
	var (
		filter string
		args   []interface{}
	)
	if req.CustomerID != "" {
		dateArgs := []interface{}{
			req.CustomerID,
		}
		args = append(args, dateArgs...)
		newFilter := `(t.customer_id = ?)`
		if filter != "" {
			newFilter = ` AND ` + newFilter
		}
		filter += newFilter
	}

	queryList := r.DB.Debug().Table("transactions as t").
		Joins("LEFT JOIN customers as c ON c.id = t.customer_id").
		Unscoped().
		Where(filter, args...)

	if err = queryList.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	var IdResponse []struct {
		ID string `json:"id"`
	}
	err = queryList.
		Select(`t.id as id`).
		Limit(req.Limit).
		Offset(req.Offset).
		Order(`t.created_at DESC`).Find(&IdResponse).Error

	var (
		filterItem string
		argsItem   []interface{}
	)

	var listId = []string{}
	if len(IdResponse) > 0 {
		for _, value := range IdResponse {
			listId = append(listId, value.ID)
		}
		argsItem = append(argsItem, listId)
		newFilter := `t.id in (?)`
		if filterItem != "" {
			newFilter = ` AND ` + newFilter
		}
		filterItem += newFilter
	}

	var getTransactionItem []model.GetTransactionItem
	err = r.DB.Debug().Table("transactions as t").
		Joins("LEFT JOIN customers as c ON c.id = t.customer_id").
		Joins("LEFT JOIN transaction_items as ti ON ti.transaction_id = t.id").
		Joins("LEFT JOIN products as p ON p.id = ti.product_id").
		Joins("LEFT JOIN sizes as s ON s.id = ti.size_id").
		Joins("LEFT JOIN flavors as f ON f.id = ti.flavor_id").
		Unscoped().
		Select(`
			t.id as transaction_id,
			c.name as customer_name,
			p.name as product_name,
			s.name as product_size,
			f.name as product_flavor,
			ti.quantity,
			ti.created_at
		`).
		Where(filterItem, argsItem...).Order(`t.created_at DESC`).Find(&getTransactionItem).Error
	if err != nil {
		return nil, 0, err
	}

	var mappingTransaction = make(map[string][]model.GetTransactionItem)
	var mappingTransactionCustomerName = make(map[string]string)
	for _, value := range getTransactionItem {
		transactionItemData := model.GetTransactionItem{
			TransactionId: value.TransactionId,
			ProductName:   value.ProductName,
			ProductSize:   value.ProductSize,
			ProductFlavor: value.ProductFlavor,
			Quantity:      value.Quantity,
			CreatedAt:     value.CreatedAt,
		}
		if _, exist := mappingTransaction[value.TransactionId]; exist {
			mappingTransaction[value.TransactionId] = append(mappingTransaction[value.TransactionId], transactionItemData)
		} else {
			mappingTransaction[value.TransactionId] = []model.GetTransactionItem{
				transactionItemData,
			}
		}

		mappingTransactionCustomerName[value.TransactionId] = value.CustomerName
	}

	for id, value := range mappingTransaction {
		resp = append(resp, model.GetTransactionResponse{
			ID:                  id,
			CustomerName:        mappingTransactionCustomerName[id],
			ListTransactionItem: value,
		})
	}
	return
}
