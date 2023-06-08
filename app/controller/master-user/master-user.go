package masteruser

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"Remember-Golang/app/db"
	"Remember-Golang/app/models"
	"Remember-Golang/app/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAll(c *gin.Context) {

	var MasterUser []models.MasterUser
	var count int64

	// db.DB.Find(&MasterUser) <- for get All Data
	db.DB.Where("is_deleted = ?", false).Find(&MasterUser)
	db.DB.Model(&MasterUser).Where("is_deleted = ?", false).Count(&count)
	c.JSON(http.StatusOK, gin.H{"Total Data": count, "Data": MasterUser})

}

func GetAllPagi(c *gin.Context) {

	pagination := utils.GeneratePaginationFromRequest(c)
	var MasterUser []models.MasterUser

	offset := (pagination.Page - 1) * pagination.Limit
	queryBuider := db.DB.Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	result := queryBuider.Model(MasterUser).Where(MasterUser).Find(&MasterUser)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": result.Error,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": MasterUser,
	})
}

func GetByID(c *gin.Context) {
	var MasterUser models.MasterUser
	ids := c.Param("id")
	if err := db.DB.First(&MasterUser, ids).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Data tidak ditemukan"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"master_user": MasterUser})
}

func Create(c *gin.Context) {

	var MasterUser *models.MasterUser

	if err := c.ShouldBindJSON(&MasterUser); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(MasterUser.Password)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	now := time.Now()
	newUser := models.MasterUser{
		Username:  MasterUser.Username,
		Email:     strings.ToLower(MasterUser.Email),
		Password:  hashedPassword,
		Fullname:  MasterUser.Fullname,
		CreatedAt: now,
		UpdatedAt: now,
		IsDeleted: false,
	}

	db.DB.Create(&newUser)
	c.JSON(http.StatusOK, gin.H{"master_user": newUser})
}

func Update(c *gin.Context) {
	var MasterUser models.MasterUser
	id := c.Param("id")

	if err := c.ShouldBindJSON(&MasterUser); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if db.DB.Model(&MasterUser).Where("id = ?", id).Updates(&MasterUser).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "tidak dapat mengupdate product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil diperbarui"})

}

func Delete(c *gin.Context) {

	var MasterUser models.MasterUser

	var input struct {
		Id json.Number
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	id, _ := input.Id.Int64()

	if db.DB.Model(&MasterUser).Where("id = ?", id).Update("is_deleted", true).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Gagal Menghapus data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil dihapus"})
}
