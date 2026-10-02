package controllers

import (
	"fmt"
	"recycle-trade/models"
	"recycle-trade/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TransactionController struct {
	DB *gorm.DB
}

// 单位换算为公斤
func toKg(weight float64, unit string) float64 {
	switch unit {
	case "吨":
		return weight * 1000
	case "市斤":
		return weight * 0.5
	default: // 公斤
		return weight
	}
}

type CreateTransactionReq struct {
	SellerID uint    `json:"seller_id" binding:"required"`
	GoodsID  uint    `json:"goods_id" binding:"required"`
	Weight   float64 `json:"weight" binding:"required"`
	Unit     string  `json:"unit" binding:"required"`
	Remark   string  `json:"remark"`
}

func (t *TransactionController) Create(c *gin.Context) {
	var req CreateTransactionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误: "+err.Error())
		return
	}
	// 查当天定价
	today := time.Now().Format("2006-01-02")
	var pricing models.Pricing
	if err := t.DB.Where("goods_id = ? AND date = ?", req.GoodsID, today).First(&pricing).Error; err != nil {
		utils.Fail(c, 400, "该货物今日未定价，请先设置价格")
		return
	}
	weightInKg := toKg(req.Weight, req.Unit)
	total := weightInKg * pricing.Price

	operatorID, _ := c.Get("user_id")

	tx := t.DB.Begin()
	trans := models.Transaction{
		SellerID:        req.SellerID,
		GoodsID:         req.GoodsID,
		Weight:          req.Weight,
		Unit:            req.Unit,
		WeightInKg:      weightInKg,
		UnitPrice:       pricing.Price,
		TotalAmount:     total,
		PricingID:       pricing.ID,
		OperatorID:      operatorID.(uint),
		TransactionTime: time.Now(),
	}
	if err := tx.Create(&trans).Error; err != nil {
		tx.Rollback()
		utils.Fail(c, 500, "创建交易失败")
		return
	}
	// 自动入库
	var inv models.Inventory
	if err := tx.Where("goods_id = ?", req.GoodsID).First(&inv).Error; err != nil {
		inv = models.Inventory{GoodsID: req.GoodsID, Quantity: 0}
		tx.Create(&inv)
	}
	inv.Quantity += weightInKg
	tx.Save(&inv)
	// 库存流水
	record := models.InventoryRecord{
		GoodsID:       req.GoodsID,
		Type:          "in",
		Quantity:      weightInKg,
		Unit:          "公斤",
		ReferenceType: "transaction",
		ReferenceID:   trans.ID,
		Remark:        fmt.Sprintf("交易入库 #%d", trans.ID),
		OperatorID:    operatorID.(uint),
	}
	tx.Create(&record)
	tx.Commit()

	// 重新加载关联
	t.DB.Preload("Seller").Preload("Goods").Preload("Operator").First(&trans, trans.ID)
	utils.OK(c, trans)
}

func (t *TransactionController) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	sellerID := c.Query("seller_id")
	goodsID := c.Query("goods_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	var total int64
	q := t.DB.Model(&models.Transaction{})
	if sellerID != "" {
		q = q.Where("seller_id = ?", sellerID)
	}
	if goodsID != "" {
		q = q.Where("goods_id = ?", goodsID)
	}
	if startDate != "" {
		q = q.Where("transaction_time >= ?", startDate+" 00:00:00")
	}
	if endDate != "" {
		q = q.Where("transaction_time <= ?", endDate+" 23:59:59")
	}
	q.Count(&total)

	var list []models.Transaction
	offset := (page - 1) * size
	q.Preload("Seller").Preload("Goods").Preload("Operator").
		Order("transaction_time desc").Offset(offset).Limit(size).Find(&list)

	utils.OK(c, gin.H{"list": list, "total": total, "page": page, "size": size})
}

func (t *TransactionController) Detail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var trans models.Transaction
	if err := t.DB.Preload("Seller").Preload("Goods").Preload("Operator").First(&trans, id).Error; err != nil {
		utils.Fail(c, 404, "交易不存在")
		return
	}
	utils.OK(c, trans)
}

func (t *TransactionController) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var trans models.Transaction
	if err := t.DB.First(&trans, id).Error; err != nil {
		utils.Fail(c, 404, "交易不存在")
		return
	}
	tx := t.DB.Begin()
	// 回滚库存
	var inv models.Inventory
	if err := tx.Where("goods_id = ?", trans.GoodsID).First(&inv).Error; err == nil {
		inv.Quantity -= trans.WeightInKg
		if inv.Quantity < 0 {
			inv.Quantity = 0
		}
		tx.Save(&inv)
	}
	tx.Delete(&models.InventoryRecord{}, "reference_type = ? AND reference_id = ?", "transaction", trans.ID)
	tx.Delete(&trans)
	tx.Commit()
	utils.OK(c, nil)
}

// 统计：今日交易笔数、总重量、总金额
func (t *TransactionController) Stats(c *gin.Context) {
	today := time.Now().Format("2006-01-02")
	var count int64
	var totalWeight, totalAmount float64
	t.DB.Model(&models.Transaction{}).Where("DATE(transaction_time) = ?", today).Count(&count)
	t.DB.Model(&models.Transaction{}).Where("DATE(transaction_time) = ?", today).
		Select("COALESCE(SUM(weight_in_kg),0)").Scan(&totalWeight)
	t.DB.Model(&models.Transaction{}).Where("DATE(transaction_time) = ?", today).
		Select("COALESCE(SUM(total_amount),0)").Scan(&totalAmount)
	utils.OK(c, gin.H{
		"today_count":   count,
		"today_weight":  totalWeight,
		"today_amount":  totalAmount,
	})
}
