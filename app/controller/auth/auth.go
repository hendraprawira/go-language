package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"Remember-Golang/app/db"
	"Remember-Golang/app/models"
	"Remember-Golang/app/utils"

	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"

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

	if !user.IsVerified {
		c.JSON(http.StatusForbidden, gin.H{"status": "fail", "message": "Please verify your email"})
		return
	}

	if err := utils.VerifyPassword(user.Password, payload.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid username or Password"})
		return
	}
	TokenSecret := os.Getenv("TOKEN_SECRET")
	//set token expires 60minutes = 1 hour
	token, err := utils.GenerateToken(60, user.ID, TokenSecret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	//set cookies max age 60*60seconds = 1 hour
	c.SetCookie("token", token, 60*60, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"status": "success", "token": token})
}

func LogoutUser(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func SignUpUser(c *gin.Context) {
	var payload *models.SignUpInput

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if payload.Password != payload.PasswordConfirm {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	now := time.Now()
	newUser := models.MasterUser{
		Username:  payload.Username,
		Email:     strings.ToLower(payload.Email),
		Fullname:  payload.Fullname,
		Password:  hashedPassword,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := db.DB.Create(&newUser)

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		c.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User with that email already exists"})
		return
	} else if result.Error != nil {
		c.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Something bad happened"})
		return
	}

	errs := godotenv.Load(".env")
	if errs != nil {
		fmt.Println(err)
	}
	ClientOrigin := os.Getenv("CLIENT_ORIGIN")

	// Generate Verification Code
	code := randstr.String(20)

	verification_code := utils.Encode(code)

	// Update User in Database
	newUser.VerificationCode = verification_code
	db.DB.Save(newUser)

	var firstName = newUser.Username

	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	// ? Send Email
	emailData := utils.EmailData{
		URL:       ClientOrigin + "/api/verifyemail/" + code,
		FirstName: firstName,
		Subject:   "Your account verification code",
	}

	utils.SendEmail(&newUser, &emailData, "verificationCode.html")

	message := "We sent an email with a verification code to " + newUser.Email
	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": message})
}

func VerifyEmail(c *gin.Context) {

	code := c.Params.ByName("verificationCode")
	verification_code := utils.Encode(code)

	var updatedUser models.MasterUser
	result := db.DB.First(&updatedUser, "verification_code = ?", verification_code)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid verification code or user doesn't exists"})
		return
	}

	if updatedUser.IsVerified {
		c.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User already verified"})
		return
	}

	updatedUser.VerificationCode = ""
	updatedUser.IsVerified = true
	db.DB.Save(&updatedUser)

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Email verified successfully"})
}

func ForgotPassword(c *gin.Context) {
	var payload *models.ForgotPasswordInput

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var user models.MasterUser
	result := db.DB.First(&user, "email = ?", strings.ToLower(payload.Email))
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	if !user.IsVerified {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Account not verified"})
		return
	}

	errs := godotenv.Load(".env")
	if errs != nil {
		fmt.Println(errs)
	}
	ClientOrigin := os.Getenv("CLIENT_ORIGIN")

	// Generate Verification Code
	resetToken := randstr.String(20)

	passwordResetToken := utils.Encode(resetToken)
	user.PasswordResetToken = passwordResetToken
	user.PasswordResetAt = time.Now().Add(time.Minute * 15)
	db.DB.Save(&user)

	var firstName = user.Fullname

	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	// ? Send Email
	emailData := utils.EmailData{
		URL:       ClientOrigin + "api/resetpassword/" + resetToken,
		FirstName: firstName,
		Subject:   "Your password reset token (valid for 10min)",
	}

	utils.SendEmail(&user, &emailData, "resetPassword.html")

	message := "You will receive a reset email if user with that email exist"
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
}

func ResetPassword(c *gin.Context) {
	var payload *models.ResetPasswordInput
	resetToken := c.Params.ByName("resetToken")

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if payload.Password != payload.PasswordConfirm {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}

	hashedPassword, _ := utils.HashPassword(payload.Password)

	passwordResetToken := utils.Encode(resetToken)

	var updatedUser models.MasterUser
	result := db.DB.First(&updatedUser, "password_reset_token = ? AND password_reset_at > ?", passwordResetToken, time.Now())
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "The reset token is invalid or has expired"})
		return
	}

	updatedUser.Password = hashedPassword
	updatedUser.PasswordResetToken = ""
	db.DB.Save(&updatedUser)

	c.SetCookie("token", "", -1, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Password data updated successfully"})
}
