package controllers

import (
	"recycle-trade/models"
	"recycle-trade/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SellerController struct {
	DB *gorm.DB
}

func (s *SellerController) List(c *gin.Context) {
	keyword := c.Query("keyword")
	var list []models.Seller
	q := s.DB.Order("id desc")
	if keyword != "" {
		q = q.Where("name LIKE ? OR phone LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	q.Find(&list)
	utils.OK(c, list)
}

func (s *SellerController) Create(c *gin.Context) {
	var req struct {
		Name    string `json:"name" binding:"required"`
		Phone   string `json:"phone"`
		IdCard  string `json:"id_card"`
		Address string `json:"address"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误")
		return
	}
	seller := models.Seller{Name: req.Name, Phone: req.Phone, IdCard: req.IdCard, Address: req.Address}
	s.DB.Create(&seller)
	utils.OK(c, seller)
}

func (s *SellerController) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var seller models.Seller
	if err := s.DB.First(&seller, id).Error; err != nil {
		utils.Fail(c, 404, "卖货人不存在")
		return
	}
	var req struct {
		Name    string `json:"name"`
		Phone   string `json:"phone"`
		IdCard  string `json:"id_card"`
		Address string `json:"address"`
	}
	c.ShouldBindJSON(&req)
	if req.Name != "" {
		seller.Name = req.Name
	}
	seller.Phone = req.Phone
	seller.IdCard = req.IdCard
	seller.Address = req.Address
	s.DB.Save(&seller)
	utils.OK(c, seller)
}

func (s *SellerController) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	s.DB.Delete(&models.Seller{}, id)
	utils.OK(c, nil)
}

// 卖货人交易统计
func (s *SellerController) Stats(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var count int64
	var totalWeight, totalAmount float64
	s.DB.Model(&models.Transaction{}).Where("seller_id = ?", id).Count(&count)
	s.DB.Model(&models.Transaction{}).Where("seller_id = ?", id).
		Select("COALESCE(SUM(weight_in_kg),0)").Scan(&totalWeight)
	s.DB.Model(&models.Transaction{}).Where("seller_id = ?", id).
		Select("COALESCE(SUM(total_amount),0)").Scan(&totalAmount)
	utils.OK(c, gin.H{"count": count, "total_weight": totalWeight, "total_amount": totalAmount})
}
