package customer

import (
	accountModel "bsnack/domain/api/account/model"
	"bsnack/domain/api/customer"
	"bsnack/domain/api/customer/model"
	"bsnack/domain/api/customer/repository"
	generalRepository "bsnack/domain/api/general/repository"
	"bsnack/lib/constant"
	"bsnack/lib/form"
	"bsnack/lib/response"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

type customerController struct {
	CustomerService customer.CustomerServiceInterface
}

func NewCustomerController(db *gorm.DB, redisClient *redis.Client) *customerController {
	return &customerController{
		CustomerService: customer.NewCustomerService(repository.NewCustomerRepository(db, redisClient), generalRepository.NewGeneralRepository(db)),
	}
}

func (c *customerController) Register(ctx *gin.Context) {
	var reqBody model.RegisterCustomerRequest
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
	responseCode, err := c.CustomerService.Register(reqBody)
	if err != nil {
		response.Error(ctx, responseCode, err.Error())
		return
	}
	response.Json(ctx, responseCode, nil)
}

func (c *customerController) Login(ctx *gin.Context) {
	var reqBody model.LoginCustomerRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	reqBody.Email = strings.ReplaceAll(reqBody.Email, " ", "")
	resBody, statusCode, err := c.CustomerService.Login(reqBody)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Json(ctx, statusCode, resBody)
}

func (c *customerController) GetAccount(ctx *gin.Context) {
	account := ctx.Request.Context().Value("auth_account").(accountModel.Account)
	if account.IsSeller {
		response.Error(ctx, http.StatusUnauthorized, "failed to get account info")
	}

	var err error
	page := form.SQLInjectorNumber(ctx.DefaultQuery("page", ""))
	limit := form.SQLInjectorNumber(ctx.DefaultQuery("limit", ""))

	request := model.GetCustomerRequest{}
	if page != "" && limit != "" {
		request.Page, err = strconv.Atoi(page)
		if err != nil {
			response.Error(ctx, http.StatusBadRequest, err.Error())
			return
		}

		request.Limit, err = strconv.Atoi(limit)
		if err != nil {
			response.Error(ctx, http.StatusBadRequest, err.Error())
			return
		}
		request.Offset = (request.Page - 1) * request.Limit
	}

	request.CustomerID = account.Id
	datas, count, responseCode, err := c.CustomerService.GetListCustomer(request)
	if err != nil {
		response.Error(ctx, responseCode, err.Error())
		return
	}
	response.JsonPagination(ctx, responseCode, datas, request.Page, request.Limit, count)
}

func (c *customerController) GetAllAccount(ctx *gin.Context) {
	account := ctx.Request.Context().Value("auth_account").(accountModel.Account)
	if !account.IsSeller {
		response.Error(ctx, http.StatusUnauthorized, "failed to get all account info")
	}

	var err error
	page := form.SQLInjectorNumber(ctx.DefaultQuery("page", ""))
	limit := form.SQLInjectorNumber(ctx.DefaultQuery("limit", ""))

	request := model.GetCustomerRequest{}
	if page != "" && limit != "" {
		request.Page, err = strconv.Atoi(page)
		if err != nil {
			response.Error(ctx, http.StatusBadRequest, err.Error())
			return
		}

		request.Limit, err = strconv.Atoi(limit)
		if err != nil {
			response.Error(ctx, http.StatusBadRequest, err.Error())
			return
		}
		request.Offset = (request.Page - 1) * request.Limit
	}
	request.CustomerID = ""

	datas, count, responseCode, err := c.CustomerService.GetListCustomer(request)
	if err != nil {
		response.Error(ctx, responseCode, err.Error())
		return
	}
	response.JsonPagination(ctx, responseCode, datas, request.Page, request.Limit, count)
}
