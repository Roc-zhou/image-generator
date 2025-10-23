package main

import (
	"image-generator/internal/handler"
	"image-generator/internal/service"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建图片服务
	imageService, err := service.NewImageService()
	if err != nil {
		log.Fatal("Failed to create image service:", err)
	}

	// 创建处理器
	imageHandler := handler.NewImageHandler(imageService)

	// 设置 gin 路由
	r := gin.Default()

	// 配置路由
	r.GET("/:size/:bg/:fg", imageHandler.GenerateImage)

	// 启动服务器
	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}