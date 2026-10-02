package controllers

import (
	"fmt"
	"recycle-trade/models"
	"recycle-trade/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InventoryController struct {
	DB *gorm.DB
}

func (i *InventoryController) List(c *gin.Context) {
	var list []models.Inventory
	i.DB.Preload("Goods").Order("id asc").Find(&list)
	utils.OK(c, list)
}

type StockInReq struct {
	GoodsID  uint    `json:"goods_id" binding:"required"`
	Quantity float64 `json:"quantity" binding:"required"`
	Unit     string  `json:"unit" binding:"required"`
	Remark   string  `json:"remark"`
}

func (i *InventoryController) StockIn(c *gin.Context) {
	var req StockInReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误")
		return
	}
	qtyInKg := toKg(req.Quantity, req.Unit)
	operatorID, _ := c.Get("user_id")

	tx := i.DB.Begin()
	var inv models.Inventory
	if err := tx.Where("goods_id = ?", req.GoodsID).First(&inv).Error; err != nil {
		inv = models.Inventory{GoodsID: req.GoodsID, Quantity: 0}
		tx.Create(&inv)
	}
	inv.Quantity += qtyInKg
	tx.Save(&inv)
	record := models.InventoryRecord{
		GoodsID:       req.GoodsID,
		Type:          "in",
		Quantity:      qtyInKg,
		Unit:          "公斤",
		ReferenceType: "manual",
		Remark:        req.Remark,
		OperatorID:    operatorID.(uint),
	}
	tx.Create(&record)
	tx.Commit()
	utils.OK(c, inv)
}

type StockOutReq struct {
	GoodsID  uint    `json:"goods_id" binding:"required"`
	Quantity float64 `json:"quantity" binding:"required"`
	Unit     string  `json:"unit" binding:"required"`
	Remark   string  `json:"remark"`
}

func (i *InventoryController) StockOut(c *gin.Context) {
	var req StockOutReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误")
		return
	}
	qtyInKg := toKg(req.Quantity, req.Unit)
	operatorID, _ := c.Get("user_id")

	tx := i.DB.Begin()
	var inv models.Inventory
	if err := tx.Where("goods_id = ?", req.GoodsID).First(&inv).Error; err != nil {
		tx.Rollback()
		utils.Fail(c, 404, "库存不存在")
		return
	}
	if inv.Quantity < qtyInKg {
		tx.Rollback()
		utils.Fail(c, 400, fmt.Sprintf("库存不足，当前库存 %.2f 公斤", inv.Quantity))
		return
	}
	inv.Quantity -= qtyInKg
	tx.Save(&inv)
	record := models.InventoryRecord{
		GoodsID:       req.GoodsID,
		Type:          "out",
		Quantity:      qtyInKg,
		Unit:          "公斤",
		ReferenceType: "manual",
		Remark:        req.Remark,
		OperatorID:    operatorID.(uint),
	}
	tx.Create(&record)
	tx.Commit()
	utils.OK(c, inv)
}

func (i *InventoryController) Records(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	goodsID := c.Query("goods_id")
	typ := c.Query("type")

	var total int64
	q := i.DB.Model(&models.InventoryRecord{})
	if goodsID != "" {
		q = q.Where("goods_id = ?", goodsID)
	}
	if typ != "" {
		q = q.Where("type = ?", typ)
	}
	q.Count(&total)

	var list []models.InventoryRecord
	offset := (page - 1) * size
	q.Preload("Goods").Preload("Operator").
		Order("created_at desc").Offset(offset).Limit(size).Find(&list)
	utils.OK(c, gin.H{"list": list, "total": total, "page": page, "size": size})
}
