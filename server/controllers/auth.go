package controllers

import (
	"net/http"
	"recycle-trade/middleware"
	"recycle-trade/models"
	"recycle-trade/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
}

func (a *AuthController) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误")
		return
	}
	var user models.User
	if err := a.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		utils.Fail(c, 401, "用户名或密码错误")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.Fail(c, 401, "用户名或密码错误")
		return
	}
	token, err := middleware.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		utils.Fail(c, 500, "生成token失败")
		return
	}
	utils.OK(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
			"name":     user.Name,
			"phone":    user.Phone,
		},
	})
}

func (a *AuthController) Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误")
		return
	}
	var count int64
	a.DB.Model(&models.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		utils.Fail(c, 400, "用户名已存在")
		return
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := models.User{
		Username: req.Username,
		Password: string(hashed),
		Role:     "user",
		Name:     req.Name,
		Phone:    req.Phone,
	}
	if err := a.DB.Create(&user).Error; err != nil {
		utils.Fail(c, 500, "注册失败")
		return
	}
	utils.OK(c, gin.H{"id": user.ID})
}

func (a *AuthController) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var user models.User
	if err := a.DB.First(&user, userID).Error; err != nil {
		utils.Fail(c, 404, "用户不存在")
		return
	}
	utils.OK(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
		"name":     user.Name,
		"phone":    user.Phone,
	})
}

func (a *AuthController) Logout(c *gin.Context) {
	utils.OK(c, nil)
}

func (a *AuthController) ChangePassword(c *gin.Context) {
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误")
		return
	}
	userID, _ := c.Get("user_id")
	var user models.User
	a.DB.First(&user, userID)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		utils.Fail(c, 400, "原密码错误")
		return
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	a.DB.Model(&user).Update("password", string(hashed))
	utils.OK(c, nil)
}

var _ = http.StatusOK
