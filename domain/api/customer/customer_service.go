package customer

import (
	accountModel "bsnack/domain/api/account/model"
	"bsnack/domain/api/customer/model"
	"bsnack/domain/api/customer/repository"
	generalRepository "bsnack/domain/api/general/repository"
	"bsnack/lib/constant"
	"bsnack/lib/encrypt"
	"bsnack/lib/utils"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
)

type CustomerServiceInterface interface {
	Register(req model.RegisterCustomerRequest) (responseCode int, err error)
	Login(req model.LoginCustomerRequest) (resp model.LoginCustomerResponse, responseCode int, err error)
	GetListCustomer(req model.GetCustomerRequest) (resp []model.GetCustomerResponse, count int64, responseCode int, err error)
}

type customerService struct {
	Repository  repository.CustomerRepositoryInterface
	GeneralRepo generalRepository.GeneralRepositoryInterface
}

func NewCustomerService(repository repository.CustomerRepositoryInterface, generalRepository generalRepository.GeneralRepositoryInterface) CustomerServiceInterface {
	return &customerService{
		Repository:  repository,
		GeneralRepo: generalRepository,
	}
}

func (s *customerService) Register(req model.RegisterCustomerRequest) (responseCode int, err error) {
	req.Email = strings.ToLower(req.Email)
	if _, err = s.Repository.Take([]string{"id"}, &model.Customer{Email: req.Email}); err == nil {
		responseCode = http.StatusBadRequest
		err = errors.New(constant.EmailAlreadyRegistered)
		return
	}
	if req.Password != req.PasswordConfirmation {
		responseCode = http.StatusBadRequest
		err = errors.New(constant.SignInPasswordIncorrect)
		return
	}
	encryptedPassword, err := encrypt.GenerateFromPassword(req.Password)
	if err != nil {
		log.Printf("ERROR GenerateFromPassword %v : %v", req.Email, err)
		responseCode = http.StatusInternalServerError
		return
	}
	seqNumber, err := s.GeneralRepo.GetDataSequence(nil, "C")
	if err != nil {
		log.Printf("ERROR GetDataSeq customer %v : %v", req.Email, err)
		responseCode = http.StatusInternalServerError
		return
	}
	customerID := fmt.Sprintf("%s%04d", "C", seqNumber)
	newCustomer := model.Customer{
		ID:                customerID,
		Name:              req.Name,
		Email:             req.Email,
		Phone:             req.Phone,
		EncryptedPassword: encryptedPassword,
		CreatedAt:         utils.GetLocaltime(),
		CreatedBy:         req.Email,
	}
	err = s.Repository.Create(newCustomer)
	if err != nil {
		log.Printf("ERROR CreateAccount customer %v : %v", req.Email, err)
		responseCode = http.StatusInternalServerError
		return
	}
	responseCode = http.StatusOK
	return
}

func (s *customerService) Login(req model.LoginCustomerRequest) (resp model.LoginCustomerResponse, responseCode int, err error) {
	var Customer model.Customer
	req.Email = strings.ToLower(req.Email)

	Customer, err = s.Repository.Take([]string{"id", "email", "encrypted_password", "name"},
		&model.Customer{Email: req.Email})
	if err != nil {
		log.Printf("ERROR GetAccount customer %v : %v", req.Email, err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			responseCode = http.StatusBadRequest
			err = errors.New(constant.CustomerAccountNotFound)
			return
		}
		responseCode = http.StatusInternalServerError
		return
	}
	if err = encrypt.CompareHashAndPassword(&Customer.EncryptedPassword, &req.Password); err != nil {
		log.Printf("ERROR Login customer %v : %v", req.Email, err)
		responseCode = http.StatusBadRequest
		err = errors.New(constant.SignInPasswordIncorrect)
		return
	}
	account := accountModel.Account{
		Id:       Customer.ID,
		Name:     Customer.Name,
		Email:    Customer.Email,
		IsSeller: false,
	}
	token, err := encrypt.GenerateTokenLogin(account)
	if err != nil {
		log.Printf("ERROR Login customer %v : %v", req.Email, err)
		responseCode = http.StatusInternalServerError
		return
	}
	resp.Token = token
	responseCode = http.StatusOK
	return
}

func (s *customerService) GetListCustomer(req model.GetCustomerRequest) (resp []model.GetCustomerResponse, count int64, responseCode int, err error) {
	resp, count, err = s.Repository.GetListCustomer(req)
	if err != nil {
		log.Printf("ERROR GetListProduct: %v", err)
		responseCode = http.StatusInternalServerError
		return
	}
	responseCode = http.StatusOK
	return
}
