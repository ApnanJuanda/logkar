package redeem

import (
	accountModel "bsnack/domain/api/account/model"
	customerRepository "bsnack/domain/api/customer/repository"
	generalRepository "bsnack/domain/api/general/repository"
	productRepository "bsnack/domain/api/product/repository"
	"bsnack/domain/api/redeem"
	"bsnack/domain/api/redeem/model"
	"bsnack/domain/api/redeem/repository"
	"bsnack/lib/response"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
)

type redeemController struct {
	RedeemService redeem.RedeemServiceInterface
}

func NewRedeemController(db *gorm.DB, redisClient *redis.Client) *redeemController {
	return &redeemController{
		RedeemService: redeem.NewRedeemService(
			repository.NewRedeemRepository(db),
			generalRepository.NewGeneralRepository(db),
			productRepository.NewProductRepository(db, redisClient),
			customerRepository.NewCustomerRepository(db, redisClient),
		),
	}
}

func (c *redeemController) RedeemPoint(ctx *gin.Context) {
	var reqBody model.RedeemPointRequest

	account := ctx.Request.Context().Value("auth_account").(accountModel.Account)
	if account.IsSeller {
		response.Error(ctx, http.StatusUnauthorized, "failed to create transaction")
	}
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	reqBody.CustomerID = account.Id

	responseCode, err := c.RedeemService.RedeemPoint(reqBody)
	if err != nil {
		response.Error(ctx, responseCode, err.Error())
		return
	}
	response.Json(ctx, responseCode, "")
}
