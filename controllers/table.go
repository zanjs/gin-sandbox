package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/zanjs/gin-sandbox/models"
)

type TableController struct {
	Controller
	db *gorm.DB
}

func NewTableController(db *gorm.DB) *TableController {
	return &TableController{db: db}
}

func (ctl TableController) CreateTable(c *gin.Context) {

	ctl.db.AutoMigrate(&models.User{}, &models.Article{}, &models.Tag{})

	ctl.SuccessResponse(c, gin.H{
		"users": "table",
	})
}
