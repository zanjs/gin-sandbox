package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/zanjs/gin-sandbox/models"
	"github.com/zanjs/gin-sandbox/shared/jwtauth"
	"github.com/zanjs/gin-sandbox/shared/passhash"
)

type AuthController struct {
	db *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{db}
}

type AuthJSON struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (controller AuthController) Auth(c *gin.Context) {
	var json AuthJSON
	if c.BindJSON(&json) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "パラメータが無効です",
		})
		return
	}

	pass, err := passhash.HashString(json.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "パラメータが無効です2",
		})
		return
	}

	var user models.User
	controller.db.Where(&models.User{
		Name:     json.Name,
		Password: pass,
	}).First(&user)

	if user.ID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "該当するユーザーが見つかりません",
		})
		return
	}

	claims, err := jwtauth.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "トークンの生成に失敗しました",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"user":   user,
		"jwt":    claims,
	})
}
