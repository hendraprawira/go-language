package utils

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Validate file extension
	ext := filepath.Ext(file.Filename)
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" && ext != ".pdf" && ext != ".csv" && ext != ".xls" && ext != ".xlsx" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid file type",
		})
		return
	}

	// Create directory if it doesn't exist
	uploadDir := "./uploaded"
	err = os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// get original file name
	filename := file.Filename

	// replace spaces with underscores
	filename = strings.ReplaceAll(filename, " ", "_")

	// generate timestamp
	t := time.Now().Format("2006-01-02_15-04-05")

	// Generate a unique filename
	newName := filename[:len(filename)-len(ext)] + "-" + t + ext

	// Save the file to disk
	err = c.SaveUploadedFile(file, filepath.Join(uploadDir, newName))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	host := os.Getenv("DIR_HOST")
	url := fmt.Sprintf("%suploaded/%s", host, newName)

	c.JSON(http.StatusOK, gin.H{
		"filename": newName,
		"url":      url,
		"message":  "sukses menambahkan file",
	})
}

func getFileExtension(fileName string) string {
	split := strings.Split(fileName, ".")
	if len(split) > 1 {
		return "." + split[len(split)-1]
	}
	return ""
}

func DeleteFile(c *gin.Context) {
	fileName := c.Query("filename")

	// get file extension from filename
	ext := getFileExtension(fileName)

	// check file extension and delete file from appropriate directory
	if ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".pdf" || ext == ".csv" || ext == ".xls" || ext == ".xlsx" {
		// file does not exist, return error message
		if _, err := os.Stat("./uploaded/" + filepath.Base(fileName)); os.IsNotExist(err) {
			c.JSON(404, gin.H{
				"message": "File not found",
			})
			return
		}
		err := os.Remove("./uploaded/" + filepath.Base(fileName))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "gagal menghapus file",
			})
			return
		}
	} else {
		// unsupported file extension, return error
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "gagal menghapus file",
		})
		return
	}

	// return success message
	c.JSON(http.StatusOK, gin.H{
		"message": "sukses menghapus file",
	})
}
