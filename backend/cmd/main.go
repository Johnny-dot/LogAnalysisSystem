package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"loganalysissystem/backend/pkg/processor"
)

func main() {
	r := gin.Default()
	r.POST("/api/upload", processor.UploadHandler)
	r.GET("/api/analyze", processor.AnalyzeHandler)
	log.Fatal(r.Run(":8080"))
}
