package api

import (
	"log"
	"net/http"
	"strconv"
	"vorto/backend/db"
	"vorto/backend/utils"

	"github.com/gin-gonic/gin"
)

func UploadFiles(c *gin.Context) {
	problemID, err := db.CreateProblem()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	files := form.File["files"]
	for _, file := range files {
		log.Printf("File Name: %s\n", file.Filename)
		log.Printf("File Size: %d\n", file.Size)
		log.Printf("File Header: %v\n", file.Header)

		if err := c.SaveUploadedFile(file, file.Filename); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		utils.ParseFile(file.Filename, problemID)
	}
	c.JSON(http.StatusOK, gin.H{"status": "files uploaded successfully"})
}

func ListProblems(c *gin.Context) {
	problems, err := db.ListProblems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"problems": problems})
}

func GetResults(c *gin.Context) {
	problemID, err := strconv.ParseInt(c.Query("problem_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid problem_id"})
		return
	}

	results, err := utils.CalculateResults(problemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
