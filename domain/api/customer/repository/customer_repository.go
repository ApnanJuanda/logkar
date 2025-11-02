package repository

import (
	"bsnack/domain/api/customer/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
	"math"
	"time"
)

type CustomerRepositoryInterface interface {
	Create(data model.Customer) (err error)
	Take(selectParams []string, conditions interface{}) (resp model.Customer, err error)
	GetCustomerByID(id string) (resp model.Customer, err error)
	UpdateFieldCustomer(tx *gorm.DB, customerID string, req map[string]interface{}) (err error)
	GetListCustomer(req model.GetCustomerRequest) (resp []model.GetCustomerResponse, count int64, err error)
}

type customerRepository struct {
	DB          *gorm.DB
	redisClient *redis.Client
}

func NewCustomerRepository(db *gorm.DB, redisClient *redis.Client) CustomerRepositoryInterface {
	return &customerRepository{
		DB:          db,
		redisClient: redisClient,
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
	ctx := context.Background()
	_, err = r.redisClient.Get(ctx, "accounts_all").Result()
	if req.Page == 0 && req.Limit == 0 && err == redis.Nil {
		req.Page = 1
		req.Limit = math.MaxInt64
	}
	if req.Page > 0 && req.Limit > 0 {
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
		if req.Limit == math.MaxInt64 {
			err = saveGetAccountToRedis(r.redisClient, resp)
			if err != nil {
				log.Printf("error save data getProduct to redis: %v", err)
			}
		}
	} else {
		dataFromRedis, err := getAllAccountFromRedis(r.redisClient)
		if err != nil {
			return nil, 0, err
		}
		return dataFromRedis, int64(len(dataFromRedis)), nil
	}
	return
}

func saveGetAccountToRedis(redisClient *redis.Client, datas []model.GetCustomerResponse) error {
	ctx := context.Background()

	jsonData, err := json.Marshal(datas)
	if err != nil {
		return err
	}
	return redisClient.Set(ctx, "accounts_all", jsonData, 24*time.Hour).Err()
}

func getAllAccountFromRedis(redisClient *redis.Client) ([]model.GetCustomerResponse, error) {
	ctx := context.Background()

	val, err := redisClient.Get(ctx, "accounts_all").Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("data is not found")
	} else if err != nil {
		return nil, err
	}

	var datas []model.GetCustomerResponse
	if err := json.Unmarshal([]byte(val), &datas); err != nil {
		return nil, err
	}
	return datas, nil
}
