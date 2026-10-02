package main

import (
	"fmt"
	"log"
	"recycle-trade/config"
	"recycle-trade/models"
	"recycle-trade/routes"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func main() {
	cfg := config.Load()

	// 连接数据库
	db, err := gorm.Open(mysql.Open(cfg.MySQLDSN), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	fmt.Println("数据库连接成功")

	// 自动迁移
	db.AutoMigrate(
		&models.User{},
		&models.Goods{},
		&models.Pricing{},
		&models.Seller{},
		&models.Transaction{},
		&models.Inventory{},
		&models.InventoryRecord{},
	)
	fmt.Println("数据表迁移完成")

	// 种子数据
	seedData(db)

	// Gin
	r := gin.Default()
	r.Use(corsMiddleware())

	routes.Setup(r, db)

	fmt.Printf("服务器启动在 http://localhost:%s\n", cfg.Port)
	r.Run(":" + cfg.Port)
}

func seedData(db *gorm.DB) {
	// 默认管理员
	var adminCount int64
	db.Model(&models.User{}).Where("username = ?", "admin").Count(&adminCount)
	if adminCount == 0 {
		hashed, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		db.Create(&models.User{
			Username: "admin",
			Password: string(hashed),
			Role:     "admin",
			Name:     "管理员",
		})
		fmt.Println("创建默认管理员: admin / admin123")
	}

	// 默认货物
	defaultGoods := []models.Goods{
		{Name: "废铁", Category: "黑色金属", Unit: "公斤"},
		{Name: "废铜", Category: "有色金属", Unit: "公斤"},
		{Name: "废铝", Category: "有色金属", Unit: "公斤"},
		{Name: "不锈钢", Category: "特种金属", Unit: "公斤"},
		{Name: "废纸", Category: "其他", Unit: "公斤"},
	}
	for _, g := range defaultGoods {
		var count int64
		db.Model(&models.Goods{}).Where("name = ?", g.Name).Count(&count)
		if count == 0 {
			db.Create(&g)
			db.Create(&models.Inventory{GoodsID: g.ID, Quantity: 0})
		}
	}

	// 今日默认定价
	today := time.Now().Format("2006-01-02")
	prices := map[string]float64{
		"废铁":   1.20,
		"废铜":   45.00,
		"废铝":   13.50,
		"不锈钢": 8.00,
		"废纸":   0.80,
	}
	for name, price := range prices {
		var goods models.Goods
		db.Where("name = ?", name).First(&goods)
		if goods.ID > 0 {
			var pCount int64
			db.Model(&models.Pricing{}).Where("goods_id = ? AND date = ?", goods.ID, today).Count(&pCount)
			if pCount == 0 {
				db.Create(&models.Pricing{GoodsID: goods.ID, Price: price, Date: today})
			}
		}
	}
}
