// Package util 提供通用工具函数
package util

import (
	"fmt"
	"math/rand"
)

// GenerateRandomColor 生成随机颜色（十六进制格式 #RRGGBB）
// 使用 HSL 模型：固定饱和度 70%、亮度 60%，随机色相 0-360
// 这样生成的色彩鲜艳且均匀分布，避免随机 RGB 产生的灰暗色
func GenerateRandomColor() string {
	// rand.Intn(360) 生成 0-359 的随机整数作为色相（Hue）
	hue := rand.Intn(360)
	// 饱和度 70%，亮度 60%
	return hslToHex(hue, 70, 60)
}

// hslToHex 将 HSL 颜色转换为十六进制字符串
// H: 0-360, S: 0-100, L: 0-100
// 算法参考 CSS Color Module Level 4
func hslToHex(h, s, l int) string {
	// 将 S, L 转换为 0-1 范围
	sat := float64(s) / 100.0
	light := float64(l) / 100.0

	// HSL -> RGB 转换
	c := (1 - abs(2*light-1)) * sat
	x := c * (1 - abs(float64(h)/60.0-1))
	m := light - c/2.0

	var r, g, b float64
	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}

	// 转换为 0-255 范围并格式化十六进制
	return fmt.Sprintf("#%02X%02X%02X",
		int((r+m)*255),
		int((g+m)*255),
		int((b+m)*255),
	)
}

// abs 返回浮点数的绝对值
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
