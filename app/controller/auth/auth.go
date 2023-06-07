package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"Remember-Golang/app/db"
	"Remember-Golang/app/models"
	"Remember-Golang/app/utils"

	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
)

func SignInUser(c *gin.Context) {
	var payload *models.SignInInput

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var user models.MasterUser
	result := db.DB.First(&user, "username = ?", strings.ToLower(payload.Username))
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid username or Password"})
		return
	}

	if err := utils.VerifyPassword(user.Password, payload.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid username or Password"})
		return
	}
	TokenSecret := os.Getenv("TOKEN_SECRET")
	token, err := utils.GenerateToken(60, user.ID, TokenSecret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	c.SetCookie("token", token, 60*60, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"status": "success", "token": token})
}

func LogoutUser(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
