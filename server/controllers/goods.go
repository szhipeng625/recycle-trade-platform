package controllers

import (
	"recycle-trade/models"
	"recycle-trade/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GoodsController struct {
	DB *gorm.DB
}

func (g *GoodsController) List(c *gin.Context) {
	var list []models.Goods
	g.DB.Order("id asc").Find(&list)
	utils.OK(c, list)
}

func (g *GoodsController) Create(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Category string `json:"category"`
		Unit     string `json:"unit"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误")
		return
	}
	unit := req.Unit
	if unit == "" {
		unit = "公斤"
	}
	goods := models.Goods{Name: req.Name, Category: req.Category, Unit: unit}
	if err := g.DB.Create(&goods).Error; err != nil {
		utils.Fail(c, 400, "创建失败，名称可能已存在")
		return
	}
	// 自动创建库存记录
	g.DB.Create(&models.Inventory{GoodsID: goods.ID, Quantity: 0})
	utils.OK(c, goods)
}

func (g *GoodsController) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var goods models.Goods
	if err := g.DB.First(&goods, id).Error; err != nil {
		utils.Fail(c, 404, "货物不存在")
		return
	}
	var req struct {
		Name     string `json:"name"`
		Category string `json:"category"`
		Unit     string `json:"unit"`
	}
	c.ShouldBindJSON(&req)
	if req.Name != "" {
		goods.Name = req.Name
	}
	if req.Category != "" {
		goods.Category = req.Category
	}
	if req.Unit != "" {
		goods.Unit = req.Unit
	}
	g.DB.Save(&goods)
	utils.OK(c, goods)
}

func (g *GoodsController) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	g.DB.Delete(&models.Goods{}, id)
	utils.OK(c, nil)
}
