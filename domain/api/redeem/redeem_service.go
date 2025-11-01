package redeem

import (
	customerModel "bsnack/domain/api/customer/model"
	customerRepository "bsnack/domain/api/customer/repository"
	generalRepository "bsnack/domain/api/general/repository"
	productModel "bsnack/domain/api/product/model"
	productRepository "bsnack/domain/api/product/repository"
	"bsnack/domain/api/redeem/model"
	"bsnack/domain/api/redeem/repository"
	"bsnack/lib/utils"
	"errors"
	"log"
	"net/http"
)

type RedeemServiceInterface interface {
	RedeemPoint(req model.RedeemPointRequest) (responseCode int, err error)
}

type redeemService struct {
	Repository   repository.RedeemRepositoryInterface
	GeneralRepo  generalRepository.GeneralRepositoryInterface
	ProductRepo  productRepository.ProductRepositoryInterface
	CustomerRepo customerRepository.CustomerRepositoryInterface
}

func NewRedeemService(repository repository.RedeemRepositoryInterface,
	generalRepo generalRepository.GeneralRepositoryInterface,
	productRepo productRepository.ProductRepositoryInterface,
	customerRepo customerRepository.CustomerRepositoryInterface) RedeemServiceInterface {
	return &redeemService{
		Repository:   repository,
		GeneralRepo:  generalRepo,
		ProductRepo:  productRepo,
		CustomerRepo: customerRepo,
	}
}

func (s *redeemService) RedeemPoint(req model.RedeemPointRequest) (responseCode int, err error) {
	tx := s.GeneralRepo.BeginTrans()
	defer func() {
		if err != nil {
			s.GeneralRepo.RollbackTrans(tx)
			return
		}
	}()

	localTime := utils.GetLocaltime()

	// get customer point
	getCustomer := customerModel.GetCustomerRequest{
		CustomerID: req.CustomerID,
		Page:       1,
		Limit:      10,
		Offset:     0,
	}
	customerPoint, _, err := s.CustomerRepo.GetListCustomer(getCustomer)
	if err != nil {
		log.Printf("ERROR GetCustomerPoint: %v", err)
		responseCode = http.StatusInternalServerError
		return
	}
	if len(customerPoint) < 1 {
		err = errors.New("CustomerPoint is not found")
		log.Printf("ERROR GetCustomerPoint: %v", err)
		responseCode = http.StatusInternalServerError
		return
	}

	// get exhange_point
	getProductRequest := productModel.ProductDetailParams{
		ProductID: req.ProductID,
		SizeID:    req.SizeID,
		FlavorID:  req.FlavorID,
	}
	productDetail, err := s.ProductRepo.GetProductDetailByParams(tx, getProductRequest)
	if err != nil {
		log.Printf("ERROR GetProductDetailByParams: %v", err)
		responseCode = http.StatusInternalServerError
		return
	}

	if productDetail.Stock <= 0 {
		err = errors.New("ProductStock is issuficient")
		log.Printf("ERROR RedeemPoint: %v", err)
		responseCode = http.StatusInternalServerError
		return
	}

	// validate minimum point
	if customerPoint[0].TotalPoints < (productDetail.ExchangePoint * req.Quantity) {
		err = errors.New("CustomerPoint is issuficient")
		log.Printf("ERROR RedeemPoint: %v", err)
		responseCode = http.StatusInternalServerError
		return
	}

	// decrease stock product_detail
	productUpdate := map[string]any{
		"stock":      productDetail.Stock - req.Quantity,
		"updated_at": localTime,
		"updated_by": "system",
	}
	err = s.ProductRepo.UpdateFieldProductDetail(tx, productDetail.ID, productUpdate)

	// decrease point customer
	pointOut := productDetail.ExchangePoint * req.Quantity
	balancePoint := customerPoint[0].TotalPoints - pointOut
	customerUpdate := map[string]any{
		"total_points": balancePoint,
		"updated_at":   localTime,
		"updated_by":   "system",
	}
	err = s.CustomerRepo.UpdateFieldCustomer(tx, req.CustomerID, customerUpdate)

	// save CustomerPointRedeem
	customerPointRedeem := model.CustomerPointRedeem{
		CustomerID:       req.CustomerID,
		TotalRedeemPoint: pointOut,
		CreatedAt:        localTime,
		CreatedBy:        "system",
	}
	err = s.Repository.InsertCustomerPointRedeem(tx, &customerPointRedeem)
	if err != nil {
		log.Printf("ERROR InsertCustomerPointRedeem: %v", err)
		responseCode = http.StatusInternalServerError
		return
	}
	s.GeneralRepo.CommitTrans(tx)

	redeemedProduct := model.CustomerRedeemedProduct{
		CustomerPointRedeemID: customerPointRedeem.ID,
		ProductID:             req.ProductID,
		SizeID:                req.SizeID,
		FlavorID:              req.FlavorID,
		Quantity:              req.Quantity,
		CreatedAt:             localTime,
		CreatedBy:             "system",
	}
	err = s.Repository.InsertRedeemedProduct(nil, redeemedProduct)
	if err != nil {
		log.Printf("ERROR InsertRedeemedProduct: %v", err)
		responseCode = http.StatusInternalServerError
		return
	}

	redeemReport := model.CustomerPointReport{
		CustomerID:            req.CustomerID,
		CustomerPointRedeemID: customerPointRedeem.ID,
		Status:                "redeem",
		Balance:               balancePoint,
		PointIn:               0,
		PointOut:              pointOut,
		CreatedAt:             localTime,
		CreatedBy:             "system",
	}
	err = s.Repository.InsertCustomerPointReport(nil, redeemReport)
	if err != nil {
		log.Printf("ERROR InsertCustomerPointReport: %v", err)
		responseCode = http.StatusInternalServerError
		return
	}
	responseCode = http.StatusOK
	return
}
