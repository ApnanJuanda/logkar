package seller

import (
	accountModel "bsnack/domain/api/account/model"
	generalRepository "bsnack/domain/api/general/repository"
	"bsnack/domain/api/seller/model"
	"bsnack/domain/api/seller/repository"
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

type SellerServiceInterface interface {
	Register(req model.RegisterSellerRequest) (responseCode int, err error)
	Login(req model.LoginSellerRequest) (resp model.LoginSellerResponse, responseCode int, err error)
}

type SellerService struct {
	Repository  repository.SellerRepositoryInterface
	GeneralRepo generalRepository.GeneralRepositoryInterface
}

func NewSellerService(repository repository.SellerRepositoryInterface, generalRepository generalRepository.GeneralRepositoryInterface) SellerServiceInterface {
	return &SellerService{
		Repository:  repository,
		GeneralRepo: generalRepository,
	}
}

func (s *SellerService) Register(req model.RegisterSellerRequest) (responseCode int, err error) {
	req.Email = strings.ToLower(req.Email)
	if _, err = s.Repository.Take([]string{"id"}, &model.Seller{Email: req.Email}); err == nil {
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
	seqNumber, err := s.GeneralRepo.GetDataSequence(nil, "SE")
	if err != nil {
		log.Printf("ERROR GetDataSeq seller %v : %v", req.Email, err)
		responseCode = http.StatusInternalServerError
		return
	}
	sellerID := fmt.Sprintf("%s%04d", "SE", seqNumber)
	newSeller := model.Seller{
		ID:                sellerID,
		Name:              req.Name,
		Email:             req.Email,
		EncryptedPassword: encryptedPassword,
		CreatedAt:         utils.GetLocaltime(),
		CreatedBy:         req.Email,
	}
	err = s.Repository.Create(newSeller)
	if err != nil {
		log.Printf("ERROR CreateAccount seller %v : %v", req.Email, err)
		responseCode = http.StatusInternalServerError
		return
	}
	responseCode = http.StatusOK
	return
}

func (s *SellerService) Login(req model.LoginSellerRequest) (resp model.LoginSellerResponse, responseCode int, err error) {
	var Seller model.Seller
	req.Email = strings.ToLower(req.Email)

	Seller, err = s.Repository.Take([]string{"id", "email", "encrypted_password", "name"},
		&model.Seller{Email: req.Email})
	if err != nil {
		log.Printf("ERROR GetAccount seller %v : %v", req.Email, err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			responseCode = http.StatusBadRequest
			err = errors.New(constant.SellerAccountNotFound)
			return
		}
		responseCode = http.StatusInternalServerError
		return
	}
	if err = encrypt.CompareHashAndPassword(&Seller.EncryptedPassword, &req.Password); err != nil {
		log.Printf("ERROR Login seller %v : %v", req.Email, err)
		responseCode = http.StatusBadRequest
		err = errors.New(constant.SignInPasswordIncorrect)
		return
	}
	account := accountModel.Account{
		Id:       Seller.ID,
		Name:     Seller.Name,
		Email:    Seller.Email,
		IsSeller: true,
	}
	token, err := encrypt.GenerateTokenLogin(account)
	if err != nil {
		log.Printf("ERROR Login seller %v : %v", req.Email, err)
		responseCode = http.StatusInternalServerError
		return
	}
	resp.Token = token
	responseCode = http.StatusOK
	return
}
