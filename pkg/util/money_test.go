package util

import (
	"testing"
)

func TestMoney_YuanToCents(t *testing.T) {
	tests := []struct {
		name     string
		yuan     float64
		expected int64
	}{
		{"正数", 1.23, 123},
		{"负数", -1.23, -123},
		{"零", 0.0, 0},
		{"整数", 1.0, 100},
		{"小数", 0.01, 1},
		{"大数", 999.99, 99999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MoneyUtil.YuanToCents(tt.yuan)
			if result != tt.expected {
				t.Errorf("YuanToCents(%v) = %v, want %v", tt.yuan, result, tt.expected)
			}
		})
	}
}

func TestMoney_CentsToYuan(t *testing.T) {
	tests := []struct {
		name     string
		cents    int64
		expected float64
	}{
		{"正数", 123, 1.23},
		{"负数", -123, -1.23},
		{"零", 0, 0.0},
		{"整数", 100, 1.0},
		{"小数", 1, 0.01},
		{"大数", 99999, 999.99},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MoneyUtil.CentsToYuan(tt.cents)
			if result != tt.expected {
				t.Errorf("CentsToYuan(%v) = %v, want %v", tt.cents, result, tt.expected)
			}
		})
	}
}

func TestMoney_FormatYuan(t *testing.T) {
	tests := []struct {
		name     string
		yuan     float64
		expected string
	}{
		{"正数", 1.23, "1.23"},
		{"负数", -1.23, "-1.23"},
		{"零", 0.0, "0.00"},
		{"整数", 1.0, "1.00"},
		{"小数", 0.01, "0.01"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MoneyUtil.FormatYuan(tt.yuan)
			if result != tt.expected {
				t.Errorf("FormatYuan(%v) = %v, want %v", tt.yuan, result, tt.expected)
			}
		})
	}
}

func TestMoney_FormatCents(t *testing.T) {
	tests := []struct {
		name     string
		cents    int64
		expected string
	}{
		{"正数", 123, "1.23"},
		{"负数", -123, "-1.23"},
		{"零", 0, "0.00"},
		{"整数", 100, "1.00"},
		{"小数", 1, "0.01"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MoneyUtil.FormatCents(tt.cents)
			if result != tt.expected {
				t.Errorf("FormatCents(%v) = %v, want %v", tt.cents, result, tt.expected)
			}
		})
	}
}

func TestMoney_ParseYuan(t *testing.T) {
	tests := []struct {
		name      string
		yuanStr   string
		expected  int64
		expectErr bool
	}{
		{"正数", "1.23", 123, false},
		{"负数", "-1.23", -123, false},
		{"零", "0.00", 0, false},
		{"整数", "1", 100, false},
		{"带货币符号", "¥1.23", 123, false},
		{"带空格", " 1.23 ", 123, false},
		{"无效格式", "abc", 0, true},
		{"空字符串", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := MoneyUtil.ParseYuan(tt.yuanStr)
			if (err != nil) != tt.expectErr {
				t.Errorf("ParseYuan(%v) error = %v, expectErr %v", tt.yuanStr, err, tt.expectErr)
				return
			}
			if result != tt.expected {
				t.Errorf("ParseYuan(%v) = %v, want %v", tt.yuanStr, result, tt.expected)
			}
		})
	}
}

func TestMoney_Add(t *testing.T) {
	tests := []struct {
		name     string
		cents1   int64
		cents2   int64
		expected int64
	}{
		{"正数相加", 100, 200, 300},
		{"负数相加", -100, -200, -300},
		{"正负相加", 100, -50, 50},
		{"零相加", 100, 0, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MoneyUtil.Add(tt.cents1, tt.cents2)
			if result != tt.expected {
				t.Errorf("Add(%v, %v) = %v, want %v", tt.cents1, tt.cents2, result, tt.expected)
			}
		})
	}
}

func TestMoney_Subtract(t *testing.T) {
	tests := []struct {
		name     string
		cents1   int64
		cents2   int64
		expected int64
	}{
		{"正数相减", 300, 100, 200},
		{"负数相减", -300, -100, -200},
		{"正负相减", 100, -50, 150},
		{"零相减", 100, 0, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MoneyUtil.Subtract(tt.cents1, tt.cents2)
			if result != tt.expected {
				t.Errorf("Subtract(%v, %v) = %v, want %v", tt.cents1, tt.cents2, result, tt.expected)
			}
		})
	}
}

func TestMoney_Multiply(t *testing.T) {
	tests := []struct {
		name       string
		cents      int64
		multiplier int64
		expected   int64
	}{
		{"正数相乘", 100, 2, 200},
		{"负数相乘", -100, 2, -200},
		{"零相乘", 100, 0, 0},
		{"一相乘", 100, 1, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MoneyUtil.Multiply(tt.cents, tt.multiplier)
			if result != tt.expected {
				t.Errorf("Multiply(%v, %v) = %v, want %v", tt.cents, tt.multiplier, result, tt.expected)
			}
		})
	}
}

func TestMoney_Divide(t *testing.T) {
	tests := []struct {
		name     string
		cents    int64
		divisor  int64
		expected int64
	}{
		{"正数相除", 200, 2, 100},
		{"负数相除", -200, 2, -100},
		{"零相除", 100, 0, 0},
		{"一相除", 100, 1, 100},
		{"四舍五入", 150, 3, 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MoneyUtil.Divide(tt.cents, tt.divisor)
			if result != tt.expected {
				t.Errorf("Divide(%v, %v) = %v, want %v", tt.cents, tt.divisor, result, tt.expected)
			}
		})
	}
}

func TestMoney_IsPositive(t *testing.T) {
	tests := []struct {
		name     string
		cents    int64
		expected bool
	}{
		{"正数", 100, true},
		{"负数", -100, false},
		{"零", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MoneyUtil.IsPositive(tt.cents)
			if result != tt.expected {
				t.Errorf("IsPositive(%v) = %v, want %v", tt.cents, result, tt.expected)
			}
		})
	}
}

func TestMoney_IsNegative(t *testing.T) {
	tests := []struct {
		name     string
		cents    int64
		expected bool
	}{
		{"正数", 100, false},
		{"负数", -100, true},
		{"零", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MoneyUtil.IsNegative(tt.cents)
			if result != tt.expected {
				t.Errorf("IsNegative(%v) = %v, want %v", tt.cents, result, tt.expected)
			}
		})
	}
}

func TestMoney_IsZero(t *testing.T) {
	tests := []struct {
		name     string
		cents    int64
		expected bool
	}{
		{"正数", 100, false},
		{"负数", -100, false},
		{"零", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MoneyUtil.IsZero(tt.cents)
			if result != tt.expected {
				t.Errorf("IsZero(%v) = %v, want %v", tt.cents, result, tt.expected)
			}
		})
	}
}

func TestMoney_Abs(t *testing.T) {
	tests := []struct {
		name     string
		cents    int64
		expected int64
	}{
		{"正数", 100, 100},
		{"负数", -100, 100},
		{"零", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MoneyUtil.Abs(tt.cents)
			if result != tt.expected {
				t.Errorf("Abs(%v) = %v, want %v", tt.cents, result, tt.expected)
			}
		})
	}
}

func TestMoney_Compare(t *testing.T) {
	tests := []struct {
		name     string
		cents1   int64
		cents2   int64
		expected int
	}{
		{"相等", 100, 100, 0},
		{"大于", 200, 100, 1},
		{"小于", 100, 200, -1},
		{"负数比较", -100, -200, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MoneyUtil.Compare(tt.cents1, tt.cents2)
			if result != tt.expected {
				t.Errorf("Compare(%v, %v) = %v, want %v", tt.cents1, tt.cents2, result, tt.expected)
			}
		})
	}
}
