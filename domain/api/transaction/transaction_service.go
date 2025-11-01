package transaction

import (
	customerRepository "bsnack/domain/api/customer/repository"
	generalRepository "bsnack/domain/api/general/repository"
	productModel "bsnack/domain/api/product/model"
	productRepository "bsnack/domain/api/product/repository"
	redeemModel "bsnack/domain/api/redeem/model"
	redeemRepository "bsnack/domain/api/redeem/repository"
	"bsnack/domain/api/transaction/model"
	"bsnack/domain/api/transaction/repository"
	"bsnack/lib/utils"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"math"
	"net/http"
)

type TransactionServiceInterface interface {
	CreateTransaction(req model.TransactionRequest) (responseCode int, err error)
	GetListTransaction(req model.GetTransactionRequest) (resp []model.GetTransactionResponse, count int64, responseCode int, err error)
}

type transactionService struct {
	Repository         repository.TransactionRepositoryInterface
	ProductRepository  productRepository.ProductRepositoryInterface
	CustomerRepository customerRepository.CustomerRepositoryInterface
	GeneralRepository  generalRepository.GeneralRepositoryInterface
	RedeemRepository   redeemRepository.RedeemRepositoryInterface
}

func NewTransactionService(
	repository repository.TransactionRepositoryInterface,
	productRepository productRepository.ProductRepositoryInterface,
	customerRepository customerRepository.CustomerRepositoryInterface,
	generalRepository generalRepository.GeneralRepositoryInterface,
	repositoryInterface redeemRepository.RedeemRepositoryInterface) TransactionServiceInterface {
	return &transactionService{
		Repository:         repository,
		ProductRepository:  productRepository,
		CustomerRepository: customerRepository,
		GeneralRepository:  generalRepository,
		RedeemRepository:   repositoryInterface,
	}
}

func (s *transactionService) CreateTransaction(req model.TransactionRequest) (responseCode int, err error) {
	tx := s.GeneralRepository.BeginTrans()
	defer func() {
		if err != nil {
			s.GeneralRepository.RollbackTrans(tx)
			return
		}
	}()

	var totalAmount float64
	var listTransactionItem = []model.TransactionItem{}
	var mappingPickQuantity = make(map[int64]int)
	transactionID := uuid.NewString()
	localTime := utils.GetLocaltime()
	for _, item := range req.ListItem {
		params := productModel.ProductDetailParams{
			ProductID: item.ProductID,
			SizeID:    item.SizeID,
			FlavorID:  item.FlavorID,
		}
		productDetail, err := s.ProductRepository.GetProductDetailByParams(tx, params)
		if err != nil {
			log.Printf("ERROR GetProductDetailByParams %v : %v", item.ProductID, err)
			return http.StatusBadRequest, err
		}
		if productDetail.Stock < item.Quantity {
			errMessage := fmt.Sprintf("stock product %v is insufficient", item.ProductID)
			log.Printf(errMessage)
			err = errors.New(errMessage)
			return http.StatusBadRequest, err
		}

		subTotal := float64(item.Quantity) * productDetail.Price
		transactionItem := model.TransactionItem{
			TransactionID: transactionID,
			ProductID:     item.ProductID,
			SizeID:        item.SizeID,
			FlavorID:      item.FlavorID,
			Quantity:      item.Quantity,
			Subtotal:      subTotal,
			CreatedAt:     localTime,
			CreatedBy:     req.CustomerEmail,
		}
		listTransactionItem = append(listTransactionItem, transactionItem)

		mappingPickQuantity[productDetail.ID] = productDetail.Stock - item.Quantity
		totalAmount += subTotal
	}

	newTransaction := model.Transaction{
		ID:          transactionID,
		CustomerID:  req.CustomerID,
		TotalAmount: totalAmount,
		Status:      10,
		Note:        "",
		CreatedAt:   localTime,
		CreatedBy:   req.CustomerEmail,
	}

	err = s.Repository.CreateTransaction(tx, newTransaction)
	if err != nil {
		log.Printf("ERROR CreateTransaction customer %v : %v", req.CustomerEmail, err)
		return http.StatusInternalServerError, err
	}

	for productDetailId, remainingStock := range mappingPickQuantity {
		dataUpdate := map[string]any{
			"stock":      remainingStock,
			"updated_at": localTime,
			"updated_by": "system",
		}
		err = s.ProductRepository.UpdateFieldProductDetail(tx, productDetailId, dataUpdate)
	}

	cashbackPoint := math.Floor(totalAmount / 1000)
	if cashbackPoint > 0 {
		customer, err := s.CustomerRepository.GetCustomerByID(req.CustomerID)
		if err != nil {
			log.Printf("ERROR GetCustomerByID %v : %v", req.CustomerID, err)
			return http.StatusInternalServerError, err
		}
		dataUpdate := map[string]any{
			"total_points": customer.TotalPoints + int(cashbackPoint),
			"updated_at":   localTime,
			"updated_by":   "system",
		}
		err = s.CustomerRepository.UpdateFieldCustomer(tx, customer.ID, dataUpdate)

		// save customer_point_reports
		customerPointReport := redeemModel.CustomerPointReport{
			CustomerID:            req.CustomerID,
			TransactionID:         transactionID,
			CustomerPointRedeemID: 0,
			Status:                "cashback",
			Balance:               customer.TotalPoints + int(cashbackPoint),
			PointIn:               int(cashbackPoint),
			PointOut:              0,
			CreatedAt:             localTime,
			CreatedBy:             "system",
		}
		err = s.RedeemRepository.InsertCustomerPointReport(tx, customerPointReport)
	}
	s.GeneralRepository.CommitTrans(tx)

	err = s.Repository.SaveTransactionItems(nil, listTransactionItem)
	if err != nil {
		log.Printf("ERROR SaveTransactionItems customer %v : %v", req.CustomerEmail, err)
		return http.StatusInternalServerError, err
	}
	return
}

func (s *transactionService) GetListTransaction(req model.GetTransactionRequest) (resp []model.GetTransactionResponse, count int64, responseCode int, err error) {
	resp, count, err = s.Repository.GetListTransaction(req)
	if err != nil {
		log.Printf("ERROR GetListProduct: %v", err)
		responseCode = http.StatusInternalServerError
		return
	}
	responseCode = http.StatusOK
	return
}
