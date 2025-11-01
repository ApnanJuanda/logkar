package collection

import (
	"bsnack/app/controller/customer"
	"bsnack/app/controller/product"
	"bsnack/app/controller/redeem"
	"bsnack/app/controller/seller"
	"bsnack/app/controller/transaction"
	"bsnack/lib/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func ApiRouter(db *gorm.DB, redisClient *redis.Client, api *gin.RouterGroup) {
	sellerCtrl := seller.NewSellerController(db)
	sellerGroup := api.Group("seller")
	{
		sellerGroup.POST("/register", sellerCtrl.Register)
		sellerGroup.POST("/login", sellerCtrl.Login)
	}

	customerCtrl := customer.NewCustomerController(db)
	customerGroup := api.Group("customer")
	{
		customerGroup.POST("/register", customerCtrl.Register)
		customerGroup.POST("/login", customerCtrl.Login)
		customerGroup.GET("", middleware.WithAuh(), customerCtrl.GetAccount)
		customerGroup.GET("/all", middleware.WithAuh(), customerCtrl.GetAllAccount)
	}

	productCtrl := product.NewProductController(db)
	productGroup := api.Group("product")
	{
		productGroup.POST("", middleware.WithAuh(), productCtrl.InsertProduct)
		productGroup.GET("", productCtrl.GetListProduct)
		productGroup.POST("/type", middleware.WithAuh(), productCtrl.InsertProductType)
		productGroup.POST("/detail", middleware.WithAuh(), productCtrl.InsertProductDetail)
		productGroup.POST("/size", middleware.WithAuh(), productCtrl.InsertSize)
		productGroup.GET("/size", productCtrl.GetAllSize)
		productGroup.POST("/flavor", middleware.WithAuh(), productCtrl.InsertFlavor)
		productGroup.GET("/flavor", productCtrl.GetAllFlavor)
	}

	transactionCtrl := transaction.NewTransactionController(db)
	transactionGroup := api.Group("/transaction")
	{
		transactionGroup.POST("", middleware.WithAuh(), transactionCtrl.CreateTransaction)
		transactionGroup.GET("", middleware.WithAuh(), transactionCtrl.GetTransactionAccount)
		transactionGroup.GET("/all", middleware.WithAuh(), transactionCtrl.GetAllTransaction)
	}

	redeemCtrl := redeem.NewRedeemController(db)
	redeemGroup := api.Group("/redeem")
	{
		redeemGroup.POST("", middleware.WithAuh(), redeemCtrl.RedeemPoint)
	}
}
