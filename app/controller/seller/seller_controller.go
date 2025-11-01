package seller

import (
	generalRepository "bsnack/domain/api/general/repository"
	"bsnack/domain/api/seller"
	"bsnack/domain/api/seller/model"
	"bsnack/domain/api/seller/repository"
	"bsnack/lib/constant"
	"bsnack/lib/form"
	"bsnack/lib/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

type sellerController struct {
	SellerService seller.SellerServiceInterface
}

func NewSellerController(db *gorm.DB) *sellerController {
	return &sellerController{
		SellerService: seller.NewSellerService(repository.NewSellerRepository(db), generalRepository.NewGeneralRepository(db)),
	}
}

func (c *sellerController) Register(ctx *gin.Context) {
	var reqBody model.RegisterSellerRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if len(reqBody.Password) < 8 {
		response.Error(ctx, http.StatusBadRequest, constant.PasswordValidation)
		return
	}
	if len(reqBody.Password) > 50 { //validate max 50 char SN-738
		response.Error(ctx, http.StatusBadRequest, constant.PasswordValidation)
		return
	}
	checkNewPassword := form.ValidatePassword(reqBody.Password)
	if !checkNewPassword {
		response.Error(ctx, http.StatusBadRequest, constant.PasswordValidation)
		return
	}
	reqBody.Email = strings.ReplaceAll(reqBody.Email, " ", "") // trim space in email
	responseCode, err := c.SellerService.Register(reqBody)
	if err != nil {
		response.Error(ctx, responseCode, err.Error())
		return
	}
	response.Json(ctx, responseCode, nil)
}

func (c *sellerController) Login(ctx *gin.Context) {
	var reqBody model.LoginSellerRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	reqBody.Email = strings.ReplaceAll(reqBody.Email, " ", "")
	resBody, statusCode, err := c.SellerService.Login(reqBody)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Json(ctx, statusCode, resBody)
}
