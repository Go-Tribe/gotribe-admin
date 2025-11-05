package util

import (
	"testing"
)

func TestPassword_GenPasswd(t *testing.T) {
	tests := []struct {
		name      string
		password  string
		expectErr bool
	}{
		{"正常密码", "123456", false},
		{"空密码", "", false},
		{"长密码", "this_is_a_very_long_password_123456789", false},
		{"特殊字符密码", "!@#$%^&*()", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := PasswordUtil.GenPasswd(tt.password)
			if (err != nil) != tt.expectErr {
				t.Errorf("GenPasswd() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && result == "" {
				t.Error("GenPasswd() returned empty string")
			}
		})
	}
}

func TestPassword_ComparePasswd(t *testing.T) {
	password := "123456"
	hashedPassword, err := PasswordUtil.GenPasswd(password)
	if err != nil {
		t.Fatalf("GenPasswd() error = %v", err)
	}

	tests := []struct {
		name           string
		hashedPassword string
		password       string
		expectErr      bool
	}{
		{"正确密码", hashedPassword, password, false},
		{"错误密码", hashedPassword, "wrong_password", true},
		{"空密码", hashedPassword, "", true},
		{"空哈希", "", password, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := PasswordUtil.ComparePasswd(tt.hashedPassword, tt.password)
			if (err != nil) != tt.expectErr {
				t.Errorf("ComparePasswd() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

func TestPassword_Encrypt(t *testing.T) {
	tests := []struct {
		name      string
		source    string
		expectErr bool
	}{
		{"正常字符串", "test_password", false},
		{"空字符串", "", false},
		{"特殊字符", "!@#$%^&*()", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := PasswordUtil.Encrypt(tt.source)
			if (err != nil) != tt.expectErr {
				t.Errorf("Encrypt() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && result == "" {
				t.Error("Encrypt() returned empty string")
			}
		})
	}
}

func TestPassword_Compare(t *testing.T) {
	source := "test_password"
	hashedPassword, err := PasswordUtil.Encrypt(source)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	tests := []struct {
		name           string
		hashedPassword string
		password       string
		expectErr      bool
	}{
		{"正确密码", hashedPassword, source, false},
		{"错误密码", hashedPassword, "wrong_password", true},
		{"空密码", hashedPassword, "", true},
		{"空哈希", "", source, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := PasswordUtil.Compare(tt.hashedPassword, tt.password)
			if (err != nil) != tt.expectErr {
				t.Errorf("Compare() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

func TestPassword_Consistency(t *testing.T) {
	password := "test_password"

	// 生成多个哈希值
	hashes := make([]string, 5)
	for i := 0; i < 5; i++ {
		hash, err := PasswordUtil.GenPasswd(password)
		if err != nil {
			t.Fatalf("GenPasswd() error = %v", err)
		}
		hashes[i] = hash
	}

	// 检查所有哈希值都不同（因为使用了随机盐）
	for i := 0; i < len(hashes); i++ {
		for j := i + 1; j < len(hashes); j++ {
			if hashes[i] == hashes[j] {
				t.Errorf("Generated identical hashes: %s", hashes[i])
			}
		}
	}

	// 检查所有哈希值都能验证原始密码
	for _, hash := range hashes {
		err := PasswordUtil.ComparePasswd(hash, password)
		if err != nil {
			t.Errorf("ComparePasswd() failed for hash %s: %v", hash, err)
		}
	}
}

func BenchmarkPassword_GenPasswd(b *testing.B) {
	password := "test_password"
	for i := 0; i < b.N; i++ {
		_, err := PasswordUtil.GenPasswd(password)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPassword_ComparePasswd(b *testing.B) {
	password := "test_password"
	hashedPassword, err := PasswordUtil.GenPasswd(password)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := PasswordUtil.ComparePasswd(hashedPassword, password)
		if err != nil {
			b.Fatal(err)
		}
	}
}
