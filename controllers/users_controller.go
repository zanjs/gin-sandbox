package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/zanjs/gin-sandbox/models"
	"github.com/zanjs/gin-sandbox/shared/passhash"
)

type UsersController struct {
	Controller
	db *gorm.DB
}

func NewUsersController(db *gorm.DB) *UsersController {
	return &UsersController{db: db}
}

func (ctl UsersController) GetAll(c *gin.Context) {
	users := []models.User{}
	ctl.db.Preload("Articles").Find(&users)

	ctl.SuccessResponse(c, gin.H{
		"users": users,
	})
}

func (ctl UsersController) Get(c *gin.Context) {
	id := c.Param("id")
	user := models.User{}
	ctl.db.Preload("Articles").First(&user, id)

	ctl.SuccessResponse(c, gin.H{
		"user": user,
	})
}

type CreateUserJSON struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (ctl UsersController) CreateUser(c *gin.Context) {
	var json CreateUserJSON
	if c.BindJSON(&json) != nil {
		ctl.ErrorResponse(c, http.StatusBadRequest, "参数无效")
		return
	}

	pass, err := passhash.HashString(json.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "パラメータが無効です",
		})
		return
	}

	user := models.User{}
	user.Name = json.Name
	user.Password = pass
	ctl.db.Create(&user)

	ctl.SuccessResponse(c, gin.H{
		"user": user,
	})
}
