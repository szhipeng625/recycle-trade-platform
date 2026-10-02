package routes

import (
	"recycle-trade/controllers"
	"recycle-trade/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(r *gin.Engine, db *gorm.DB) {
	auth := &controllers.AuthController{DB: db}
	goods := &controllers.GoodsController{DB: db}
	pricing := &controllers.PricingController{DB: db}
	trans := &controllers.TransactionController{DB: db}
	seller := &controllers.SellerController{DB: db}
	inventory := &controllers.InventoryController{DB: db}

	api := r.Group("/api")
	{
		// 公开接口
		api.POST("/auth/login", auth.Login)
		api.POST("/auth/register", auth.Register)

		// 需要登录
		authGroup := api.Group("")
		authGroup.Use(middleware.AuthMiddleware())
		{
			authGroup.GET("/auth/profile", auth.GetProfile)
			authGroup.POST("/auth/logout", auth.Logout)
			authGroup.POST("/auth/change-password", auth.ChangePassword)

			// 货物（登录可见，管理员可增删改）
			authGroup.GET("/goods", goods.List)
			authGroup.POST("/goods", middleware.AdminMiddleware(), goods.Create)
			authGroup.PUT("/goods/:id", middleware.AdminMiddleware(), goods.Update)
			authGroup.DELETE("/goods/:id", middleware.AdminMiddleware(), goods.Delete)

			// 定价
			authGroup.GET("/pricing", pricing.List)
			authGroup.GET("/pricing/today", pricing.Today)
			authGroup.GET("/pricing/trend", pricing.Trend)
			authGroup.POST("/pricing", middleware.AdminMiddleware(), pricing.Create)
			authGroup.DELETE("/pricing/:id", middleware.AdminMiddleware(), pricing.Delete)

			// 交易
			authGroup.GET("/transactions", trans.List)
			authGroup.GET("/transactions/stats", trans.Stats)
			authGroup.GET("/transactions/:id", trans.Detail)
			authGroup.POST("/transactions", middleware.AdminMiddleware(), trans.Create)
			authGroup.DELETE("/transactions/:id", middleware.AdminMiddleware(), trans.Delete)

			// 卖货人
			authGroup.GET("/sellers", seller.List)
			authGroup.GET("/sellers/:id/stats", seller.Stats)
			authGroup.POST("/sellers", middleware.AdminMiddleware(), seller.Create)
			authGroup.PUT("/sellers/:id", middleware.AdminMiddleware(), seller.Update)
			authGroup.DELETE("/sellers/:id", middleware.AdminMiddleware(), seller.Delete)

			// 库存
			authGroup.GET("/inventory", inventory.List)
			authGroup.GET("/inventory/records", inventory.Records)
			authGroup.POST("/inventory/in", middleware.AdminMiddleware(), inventory.StockIn)
			authGroup.POST("/inventory/out", middleware.AdminMiddleware(), inventory.StockOut)
		}
	}
}
