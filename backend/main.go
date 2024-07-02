package main

import (
	"vorto/backend/api"
	"vorto/backend/db"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	defer db.CloseDB()

	r := gin.Default()
	r.Use(cors.Default())

	r.POST("/upload", api.UploadFiles)
	r.GET("/results", api.GetResults)
	r.GET("/problems", api.ListProblems)
	r.Run(":8080")
}
