package service

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

type ImageService struct {
	font *truetype.Font
}

func NewImageService() (*ImageService, error) {
	f, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}
	return &ImageService{font: f}, nil
}

func (s *ImageService) GenerateImage(width, height int, bgColor, fgColor color.Color, text string, format string) (io.WriterTo, error) {
	// 创建新图片
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 填充背景色
	draw.Draw(img, img.Bounds(), &image.Uniform{C: bgColor}, image.Point{}, draw.Src)

	// 添加文字
	if text == "" {
		text = fmt.Sprintf("%dx%d", width, height)
	}

	// 设置字体
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(s.font)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.NewUniform(fgColor))

	// 计算字体大小和位置
	// 根据图片的宽度和高度计算合适的字体大小
	minDimension := float64(width)
	if height < width {
		minDimension = float64(height)
	}

	// 根据文本长度动态调整字体大小
	textLength := float64(len(text))
	if textLength == 0 {
		textLength = 1
	}

	// 初始字体大小设置为最小边长的 1/3
	fontSize := minDimension / 3

	// 如果文本较长，进一步调整字体大小
	if textLength > 5 {
		fontSize = fontSize * 5 / textLength
	}

	// 确保字体大小不会太小
	if fontSize < 12 {
		fontSize = 12
	}

	c.SetFontSize(fontSize)

	// 估算文本宽度
	textWidth := fontSize * textLength * 0.6

	// 如果文本宽度超过图片宽度，调整字体大小
	if textWidth > float64(width)*0.9 {
		fontSize = fontSize * float64(width) * 0.9 / textWidth
		c.SetFontSize(fontSize)
		textWidth = fontSize * textLength * 0.6
	}

	// 居中显示文字
	pt := freetype.Pt(
		(width-int(textWidth))/2,
		(height+int(fontSize))/2,
	)
	_, err := c.DrawString(text, pt)
	if err != nil {
		return nil, err
	}

	// 根据格式返回不同的图片对象
	switch strings.ToLower(format) {
	case "jpg", "jpeg":
		return &jpegImage{img: img}, nil
	case "gif":
		return &gifImage{img: img}, nil
	default: // png
		return &pngImage{img: img}, nil
	}
}

type pngImage struct {
	img *image.RGBA
}

func (p *pngImage) WriteTo(w io.Writer) (int64, error) {
	err := png.Encode(w, p.img)
	return 0, err
}

type jpegImage struct {
	img *image.RGBA
}

func (j *jpegImage) WriteTo(w io.Writer) (int64, error) {
	err := jpeg.Encode(w, j.img, &jpeg.Options{Quality: 90})
	return 0, err
}

type gifImage struct {
	img *image.RGBA
}

func (g *gifImage) WriteTo(w io.Writer) (int64, error) {
	err := gif.Encode(w, g.img, &gif.Options{NumColors: 256})
	return 0, err
}
