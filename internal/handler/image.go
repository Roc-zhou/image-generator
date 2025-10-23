package handler

import (
	"fmt"
	"image-generator/internal/service"
	"image/color"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ImageHandler struct {
	imageService *service.ImageService
}

func NewImageHandler(imageService *service.ImageService) *ImageHandler {
	return &ImageHandler{
		imageService: imageService,
	}
}

func (h *ImageHandler) GenerateImage(c *gin.Context) {
	// 解析尺寸
	size := c.Param("size")
	dimensions := strings.Split(size, "x")
	if len(dimensions) != 2 {
		c.String(http.StatusBadRequest, "Invalid size format")
		return
	}

	width, err := strconv.Atoi(dimensions[0])
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid width")
		return
	}

	height, err := strconv.Atoi(dimensions[1])
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid height")
		return
	}

	// 解析背景色和前景色
	bgColorHex := c.Param("bg")
	fgColorHex := c.Param("fg")

	// 从前景色参数中移除文件扩展名
	if ext := path.Ext(fgColorHex); ext != "" {
		fgColorHex = strings.TrimSuffix(fgColorHex, ext)
	}

	bgColor, err := parseHexColor(bgColorHex)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid background color: "+bgColorHex)
		return
	}

	fgColor, err := parseHexColor(fgColorHex)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid foreground color: "+fgColorHex)
		return
	}

	// 解析格式
	format := strings.TrimPrefix(path.Ext(c.Request.URL.Path), ".")
	if format == "" {
		format = "png"
	}

	// 解析文本
	text := c.Query("text")

	// 生成图片
	img, err := h.imageService.GenerateImage(width, height, bgColor, fgColor, text, format)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to generate image")
		return
	}

	// 设置响应头
	c.Header("Content-Type", fmt.Sprintf("image/%s", format))
	c.Header("Cache-Control", "public, max-age=31536000")

	// 写入响应
	img.WriteTo(c.Writer)
}

func parseHexColor(hex string) (color.Color, error) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) == 3 {
		hex = string([]byte{hex[0], hex[0], hex[1], hex[1], hex[2], hex[2]})
	}
	if len(hex) != 6 {
		return nil, fmt.Errorf("invalid hex color length")
	}

	rgb, err := strconv.ParseUint(hex, 16, 24)
	if err != nil {
		return nil, err
	}

	return color.RGBA{
		R: uint8(rgb >> 16),
		G: uint8((rgb >> 8) & 0xFF),
		B: uint8(rgb & 0xFF),
		A: 255,
	}, nil
}