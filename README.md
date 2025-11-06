# 动态图片生成服务

这是一个基于 Go 和 Gin 框架实现的动态图片生成服务。它可以根据 URL 参数动态生成不同尺寸、颜色和文字的图片，支持多种图片格式。

## 功能特点

- 支持自定义图片尺寸
- 支持自定义背景色和前景色
- 支持自定义文字内容
- 支持多种图片格式（PNG、JPEG、GIF）
- 自适应文字大小
- 响应速度快
- 支持图片缓存（Cache-Control）

## URL 格式

```
https://your-domain/{size}/{bgcolor}/{fgcolor}.{format}?text={text}
```

### 参数说明

- `size`: 图片尺寸，格式为 `宽x高`，如 `600x400`（必填，默认 100x100）
- `bgcolor`: 背景色，6位十六进制颜色值，如 `000`（必填，默认 000）
- `fgcolor`: 前景色（文字颜色），6位十六进制颜色值，如 `fff`（必填，默认 fff）
- `format`: 图片格式，支持 png、jpg/jpeg、gif（可选，默认 png）
- `text`: 文字内容（可选，默认显示图片尺寸）

### 示例

1. 基本用法：
```
http://localhost:8080/image/700x100/000/fff
http://localhost:8080/image/700x100/000/fff.png
http://localhost:8080/image/700x100/000/fff.jpg
```
生成一个 700x100 的黑底白字图片，显示尺寸文本 "700x100"

2. 自定义文字：
```
http://localhost:8080/image/700x100/000/fff.png?text=Hello
```
生成一个 700x100 的红底白字图片，显示文本 "Hello"

3. 不同格式：
```
http://localhost:8080/image/700x100/000/fff.png
http://localhost:8080/image/700x100/000/fff.jpg
```
生成不同格式的图片

## 安装和运行

1. 克隆项目：
```bash
git clone https://github.com/Roc-zhou/image-generator.git
```

2. 安装依赖：
```bash
cd image-generator
go mod download
```

3. 运行服务：
```bash
go run main.go
```

服务默认在 8080 端口启动。

## 依赖

- github.com/gin-gonic/gin
- github.com/golang/freetype
- golang.org/x/image