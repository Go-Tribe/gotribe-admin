package util

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Money 金额工具类，用于处理分和元之间的转换
type Money struct{}

// YuanToCents 将元转换为分
// 例如: 1.23 -> 123, -1.23 -> -123
func (m *Money) YuanToCents(yuan float64) int64 {
	return int64(math.Round(yuan * 100))
}

// CentsToYuan 将分转换为元
// 例如: 123 -> 1.23, -123 -> -1.23
func (m *Money) CentsToYuan(cents int64) float64 {
	return float64(cents) / 100.0
}

// FormatYuan 格式化金额显示（保留2位小数）
// 例如: 1.23 -> "1.23", 1.0 -> "1.00"
func (m *Money) FormatYuan(yuan float64) string {
	return fmt.Sprintf("%.2f", yuan)
}

// FormatCents 格式化分显示为元
// 例如: 123 -> "1.23", -123 -> "-1.23"
func (m *Money) FormatCents(cents int64) string {
	return m.FormatYuan(m.CentsToYuan(cents))
}

// ParseYuan 解析字符串金额为分
// 例如: "1.23" -> 123, "-1.23" -> -123
func (m *Money) ParseYuan(yuanStr string) (int64, error) {
	// 移除可能的货币符号和空格
	yuanStr = strings.TrimSpace(yuanStr)
	yuanStr = strings.Trim(yuanStr, "¥$€£")

	yuan, err := strconv.ParseFloat(yuanStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid money format: %s", yuanStr)
	}

	return m.YuanToCents(yuan), nil
}

// Add 分相加
func (m *Money) Add(cents1, cents2 int64) int64 {
	return cents1 + cents2
}

// Subtract 分相减
func (m *Money) Subtract(cents1, cents2 int64) int64 {
	return cents1 - cents2
}

// Multiply 分乘以整数
func (m *Money) Multiply(cents int64, multiplier int64) int64 {
	return cents * multiplier
}

// Divide 分除以整数（四舍五入）
func (m *Money) Divide(cents int64, divisor int64) int64 {
	if divisor == 0 {
		return 0
	}
	return int64(math.Round(float64(cents) / float64(divisor)))
}

// IsPositive 检查分是否为正数
func (m *Money) IsPositive(cents int64) bool {
	return cents > 0
}

// IsNegative 检查分是否为负数
func (m *Money) IsNegative(cents int64) bool {
	return cents < 0
}

// IsZero 检查分是否为零
func (m *Money) IsZero(cents int64) bool {
	return cents == 0
}

// Abs 获取分的绝对值
func (m *Money) Abs(cents int64) int64 {
	if cents < 0 {
		return -cents
	}
	return cents
}

// Compare 比较两个分的大小
// 返回: -1 (cents1 < cents2), 0 (cents1 == cents2), 1 (cents1 > cents2)
func (m *Money) Compare(cents1, cents2 int64) int {
	if cents1 < cents2 {
		return -1
	} else if cents1 > cents2 {
		return 1
	}
	return 0
}

// 全局实例
var MoneyUtil = &Money{}
