package middleware

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

func DeserializeUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var token string
		cookie, err := ctx.Cookie("token")
		errs := godotenv.Load(".env")
		if errs != nil {
			fmt.Println(err)
		}

		authorizationHeader := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			token = fields[1]
		} else if err == nil {
			token = cookie
		}

		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in"})
			return
		}

		TokenSecret := os.Getenv("TOKEN_SECRET")
		sub, err := utils.ValidateToken(token, TokenSecret)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail test", "message": err.Error()})
			return
		}

		var user models.MasterUser
		result := db.DB.First(&user, "id = ?", fmt.Sprint(sub))
		if result.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
			return
		}

		ctx.Set("currentUser", user)
		ctx.Set("id_user", user.ID)
		ctx.Next()
	}
}
