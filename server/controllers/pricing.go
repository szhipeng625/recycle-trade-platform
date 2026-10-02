package controllers

import (
	"recycle-trade/models"
	"recycle-trade/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PricingController struct {
	DB *gorm.DB
}

func (p *PricingController) List(c *gin.Context) {
	date := c.Query("date")
	goodsID := c.Query("goods_id")
	var list []models.Pricing
	q := p.DB.Preload("Goods").Order("date desc, id desc")
	if date != "" {
		q = q.Where("date = ?", date)
	}
	if goodsID != "" {
		q = q.Where("goods_id = ?", goodsID)
	}
	q.Find(&list)
	utils.OK(c, list)
}

func (p *PricingController) Today(c *gin.Context) {
	today := time.Now().Format("2006-01-02")
	var list []models.Pricing
	p.DB.Preload("Goods").Where("date = ?", today).Find(&list)
	utils.OK(c, list)
}

func (p *PricingController) Create(c *gin.Context) {
	var req struct {
		GoodsID uint    `json:"goods_id" binding:"required"`
		Price   float64 `json:"price" binding:"required"`
		Date    string  `json:"date"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误")
		return
	}
	date := req.Date
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	// 同一天同一货物只保留一条，更新即可
	var existing models.Pricing
	if err := p.DB.Where("goods_id = ? AND date = ?", req.GoodsID, date).First(&existing).Error; err == nil {
		existing.Price = req.Price
		p.DB.Save(&existing)
		utils.OK(c, existing)
		return
	}
	pricing := models.Pricing{GoodsID: req.GoodsID, Price: req.Price, Date: date}
	p.DB.Create(&pricing)
	utils.OK(c, pricing)
}

func (p *PricingController) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	p.DB.Delete(&models.Pricing{}, id)
	utils.OK(c, nil)
}

// 价格走势：按货物查历史价格
func (p *PricingController) Trend(c *gin.Context) {
	goodsID := c.Query("goods_id")
	days := c.DefaultQuery("days", "30")
	var list []models.Pricing
	q := p.DB.Where("date >= ?", time.Now().AddDate(0, 0, -func() int {
		n, _ := strconv.Atoi(days)
		return n
	}()).Format("2006-01-02"))
	if goodsID != "" {
		q = q.Where("goods_id = ?", goodsID)
	}
	q.Order("date asc").Find(&list)
	utils.OK(c, list)
}
