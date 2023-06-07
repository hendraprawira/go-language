package router

import (
	auth "Remember-Golang/app/controller/auth"
	masteruser "Remember-Golang/app/controller/master-user"
	"Remember-Golang/app/db"
	"Remember-Golang/app/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Routes() *gin.Engine {
	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	db.ConnectDatabase()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "OPTIONS", "GET", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	r.MaxMultipartMemory = 8 << 20 // 8 MiB

	apiUri := r.Group("/api")

	masterUserRoute := apiUri.Group("").Use(middleware.DeserializeUser())
	{
		masterUserRoute.GET("/master-user", masteruser.GetAll)
		masterUserRoute.GET("/master-user-pagi", masteruser.GetAllPagi)
		masterUserRoute.GET("/master-user-by-id/:id", masteruser.GetByID)
		masterUserRoute.POST("/master-user-create", masteruser.Create)
		masterUserRoute.PUT("/master-user-update/:id", masteruser.Update)
		masterUserRoute.DELETE("/master-user-delete", masteruser.Delete)
	}

	authRoute := apiUri.Group("")
	{
		authRoute.POST("/login", auth.SignInUser)
		authRoute.GET("/logout", middleware.DeserializeUser(), auth.LogoutUser)
		authRoute.POST("/register", auth.SignUpUser)
		authRoute.GET("/verifyemail/:verificationCode", auth.VerifyEmail)
	}

	return r
}
