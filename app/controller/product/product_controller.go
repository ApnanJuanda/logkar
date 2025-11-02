package product

import (
	accountModel "bsnack/domain/api/account/model"
	generalRepository "bsnack/domain/api/general/repository"
	"bsnack/domain/api/product"
	"bsnack/domain/api/product/model"
	"bsnack/domain/api/product/repository"
	"bsnack/lib/form"
	"bsnack/lib/response"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type productController struct {
	ProductService product.ProductServiceInterface
}

func NewProductController(db *gorm.DB, redisClient *redis.Client) *productController {
	return &productController{
		ProductService: product.NewProductService(repository.NewProductRepository(db, redisClient), generalRepository.NewGeneralRepository(db)),
	}
}

func (c *productController) InsertSize(ctx *gin.Context) {
	var reqBody model.InsertSizeRequest

	account := ctx.Request.Context().Value("auth_account").(accountModel.Account)
	if !account.IsSeller {
		response.Error(ctx, http.StatusUnauthorized, "failed to insert size")
	}

	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	reqBody.SellerEmail = account.Email
	responseCode, err := c.ProductService.InsertSize(reqBody)
	if err != nil {
		response.Error(ctx, responseCode, err.Error())
		return
	}
	response.Json(ctx, responseCode, "")
}

func (c *productController) GetAllSize(ctx *gin.Context) {
	datas, responseCode, err := c.ProductService.GetAllSize()
	if err != nil {
		response.Error(ctx, responseCode, err.Error())
		return
	}
	response.Json(ctx, responseCode, datas)
}

func (c *productController) InsertFlavor(ctx *gin.Context) {
	var reqBody model.InsertFlavorRequest

	account := ctx.Request.Context().Value("auth_account").(accountModel.Account)
	if !account.IsSeller {
		response.Error(ctx, http.StatusUnauthorized, "failed to insert size")
	}

	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	reqBody.SellerEmail = account.Email
	responseCode, err := c.ProductService.InsertFlavor(reqBody)
	if err != nil {
		response.Error(ctx, responseCode, err.Error())
		return
	}
	response.Json(ctx, responseCode, "")
}

func (c *productController) GetAllFlavor(ctx *gin.Context) {
	datas, responseCode, err := c.ProductService.GetAllFlavor()
	if err != nil {
		response.Error(ctx, responseCode, err.Error())
		return
	}
	response.Json(ctx, responseCode, datas)
}

func (c *productController) InsertProductType(ctx *gin.Context) {
	var reqBody model.InsertProductTypeRequest

	account := ctx.Request.Context().Value("auth_account").(accountModel.Account)
	if !account.IsSeller {
		response.Error(ctx, http.StatusUnauthorized, "failed to insert product")
	}

	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	reqBody.SellerEmail = account.Email
	responseCode, err := c.ProductService.InsertProductType(reqBody)
	if err != nil {
		response.Error(ctx, responseCode, err.Error())
		return
	}
	response.Json(ctx, responseCode, "")
}

func (c *productController) GetProductType(ctx *gin.Context) {
	datas, responseCode, err := c.ProductService.GetProductType()
	if err != nil {
		response.Error(ctx, responseCode, err.Error())
		return
	}
	response.Json(ctx, responseCode, datas)
}

func (c *productController) InsertProduct(ctx *gin.Context) {
	var reqBody model.InsertProductRequest

	account := ctx.Request.Context().Value("auth_account").(accountModel.Account)
	if !account.IsSeller {
		response.Error(ctx, http.StatusUnauthorized, "failed to insert product")
	}

	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	reqBody.SellerID = account.Id
	reqBody.SellerEmail = account.Email
	responseCode, err := c.ProductService.InsertProduct(reqBody)
	if err != nil {
		response.Error(ctx, responseCode, err.Error())
		return
	}
	response.Json(ctx, responseCode, "")
}

func (c *productController) InsertProductDetail(ctx *gin.Context) {
	var reqBody model.InsertProductDetailRequest

	account := ctx.Request.Context().Value("auth_account").(accountModel.Account)
	if !account.IsSeller {
		response.Error(ctx, http.StatusUnauthorized, "failed to insert product")
	}

	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	reqBody.SellerID = account.Id
	reqBody.SellerEmail = account.Email
	responseCode, err := c.ProductService.InsertProductDetail(reqBody)
	if err != nil {
		response.Error(ctx, responseCode, err.Error())
		return
	}
	response.Json(ctx, responseCode, "")
}

func (c *productController) GetListProduct(ctx *gin.Context) {
	var err error
	page := form.SQLInjectorNumber(ctx.DefaultQuery("page", ""))
	limit := form.SQLInjectorNumber(ctx.DefaultQuery("limit", ""))

	startDate := form.SQLInjector(ctx.DefaultQuery("start_date", ""))
	endDate := form.SQLInjector(ctx.DefaultQuery("end_date", ""))
	request := model.GetProductRequest{
		StartDate: startDate,
		EndDate:   endDate,
	}

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

	datas, count, responseCode, err := c.ProductService.GetListProduct(request)
	if err != nil {
		response.Error(ctx, responseCode, err.Error())
		return
	}
	response.JsonPagination(ctx, responseCode, datas, request.Page, request.Limit, count)
}
