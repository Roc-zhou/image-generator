package handler

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"image-generator/internal/service"
	"image/color"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

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
	var width, height int
	if strings.Contains(size, "x") {
		parts := strings.Split(size, "x")
		if len(parts) != 2 {
			c.String(http.StatusBadRequest, "Invalid size format")
			return
		}
		w, err := strconv.Atoi(parts[0])
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid width")
			return
		}
		hgt, err := strconv.Atoi(parts[1])
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid height")
			return
		}
		width = w
		height = hgt
	} else {
		v, err := strconv.Atoi(size)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid size")
			return
		}
		width = v
		height = v
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

	// 将图片写入内存缓冲以便计算 ETag
	buf := &bytes.Buffer{}
	_, err = img.WriteTo(buf)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to write image")
		return
	}

	// 计算 ETag（使用 SHA1）
	sum := sha1.Sum(buf.Bytes())
	etag := "\"" + hex.EncodeToString(sum[:]) + "\""

	// 如果客户端有 If-None-Match，且匹配 ETag，则返回 304
	if inm := c.GetHeader("If-None-Match"); inm != "" {
		if inm == etag {
			maxAge := 60 * 60 * 24 * 30 // 30 days
			c.Header("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
			c.Status(http.StatusNotModified)
			return
		}
	}

	// 设置响应头并返回图片
	c.Header("ETag", etag)
	c.Header("Content-Type", fmt.Sprintf("image/%s", format))
	c.Header("Content-Length", strconv.Itoa(buf.Len()))
	c.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))

	_, _ = c.Writer.Write(buf.Bytes())
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
