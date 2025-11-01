package transaction

import (
	accountModel "bsnack/domain/api/account/model"
	customerRepository "bsnack/domain/api/customer/repository"
	generalRepository "bsnack/domain/api/general/repository"
	productRepository "bsnack/domain/api/product/repository"
	redeemRepository "bsnack/domain/api/redeem/repository"
	"bsnack/domain/api/transaction"
	"bsnack/domain/api/transaction/model"
	"bsnack/domain/api/transaction/repository"
	"bsnack/lib/form"
	"bsnack/lib/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type transactionController struct {
	TransactionService transaction.TransactionServiceInterface
}

func NewTransactionController(db *gorm.DB) *transactionController {
	return &transactionController{transaction.NewTransactionService(
		repository.NewTransactionRepository(db),
		productRepository.NewProductRepository(db),
		customerRepository.NewCustomerRepository(db),
		generalRepository.NewGeneralRepository(db),
		redeemRepository.NewRedeemRepository(db),
	)}
}

func (c *transactionController) CreateTransaction(ctx *gin.Context) {
	var reqBody model.TransactionRequest

	account := ctx.Request.Context().Value("auth_account").(accountModel.Account)
	if account.IsSeller {
		response.Error(ctx, http.StatusUnauthorized, "failed to create transaction")
	}
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	reqBody.CustomerID = account.Id
	reqBody.CustomerEmail = account.Email

	responseCode, err := c.TransactionService.CreateTransaction(reqBody)
	if err != nil {
		response.Error(ctx, responseCode, err.Error())
		return
	}
	response.Json(ctx, responseCode, "")
}

func (c *transactionController) GetTransactionAccount(ctx *gin.Context) {
	account := ctx.Request.Context().Value("auth_account").(accountModel.Account)
	if account.IsSeller {
		response.Error(ctx, http.StatusUnauthorized, "failed to get transaction")
	}

	var err error
	page := form.SQLInjectorNumber(ctx.DefaultQuery("page", "1"))
	limit := form.SQLInjectorNumber(ctx.DefaultQuery("limit", "10"))

	request := model.GetTransactionRequest{}
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
	request.CustomerID = account.Id

	datas, count, responseCode, err := c.TransactionService.GetListTransaction(request)
	if err != nil {
		response.Error(ctx, responseCode, err.Error())
		return
	}
	response.JsonPagination(ctx, responseCode, datas, request.Page, request.Limit, count)
}

func (c *transactionController) GetAllTransaction(ctx *gin.Context) {
	account := ctx.Request.Context().Value("auth_account").(accountModel.Account)
	if !account.IsSeller {
		response.Error(ctx, http.StatusUnauthorized, "failed to get all transaction")
	}

	var err error
	page := form.SQLInjectorNumber(ctx.DefaultQuery("page", "1"))
	limit := form.SQLInjectorNumber(ctx.DefaultQuery("limit", "10"))

	request := model.GetTransactionRequest{}
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
	request.CustomerID = ""

	datas, count, responseCode, err := c.TransactionService.GetListTransaction(request)
	if err != nil {
		response.Error(ctx, responseCode, err.Error())
		return
	}
	response.JsonPagination(ctx, responseCode, datas, request.Page, request.Limit, count)
}
