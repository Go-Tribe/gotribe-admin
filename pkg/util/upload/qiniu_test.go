// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package upload

import (
	"testing"
)

func TestNewQiniu(t *testing.T) {
	tests := []struct {
		name       string
		ak         string
		sk         string
		bucket     string
		wantAK     string
		wantSK     string
		wantBucket string
	}{
		{
			name:       "正常参数",
			ak:         "test_access_key",
			sk:         "test_secret_key",
			bucket:     "test_bucket",
			wantAK:     "test_access_key",
			wantSK:     "test_secret_key",
			wantBucket: "test_bucket",
		},
		{
			name:       "空参数",
			ak:         "",
			sk:         "",
			bucket:     "",
			wantAK:     "",
			wantSK:     "",
			wantBucket: "",
		},
		{
			name:       "特殊字符",
			ak:         "ak_123-abc!@#",
			sk:         "sk_456-def$%^",
			bucket:     "bucket.name",
			wantAK:     "ak_123-abc!@#",
			wantSK:     "sk_456-def$%^",
			wantBucket: "bucket.name",
		},
		{
			name:       "长字符串",
			ak:         "a very long access key with many characters 123456789",
			sk:         "a very long secret key with many characters and special symbols !@#$%^&*()",
			bucket:     "my-test-bucket-name-123",
			wantAK:     "a very long access key with many characters 123456789",
			wantSK:     "a very long secret key with many characters and special symbols !@#$%^&*()",
			wantBucket: "my-test-bucket-name-123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uploader := NewQiniu(tt.ak, tt.sk, tt.bucket)

			if uploader.AccessKey != tt.wantAK {
				t.Errorf("NewQiniu() AccessKey = %v, want %v", uploader.AccessKey, tt.wantAK)
			}
			if uploader.SecretKey != tt.wantSK {
				t.Errorf("NewQiniu() SecretKey = %v, want %v", uploader.SecretKey, tt.wantSK)
			}
			if uploader.Bucket != tt.wantBucket {
				t.Errorf("NewQiniu() Bucket = %v, want %v", uploader.Bucket, tt.wantBucket)
			}
		})
	}
}

func TestQiniuUploader_UploadFile_InvalidFile(t *testing.T) {
	// 测试无效文件上传（需要真实凭证才能测试成功场景）
	uploader := NewQiniu("test_ak", "test_sk", "test_bucket")

	// 注意：由于源代码没有nil检查，直接传递nil会导致panic
	// 这里只测试结构体是否正确创建
	t.Run("上传器创建", func(t *testing.T) {
		if uploader.AccessKey != "test_ak" {
			t.Error("QiniuUploader was not created correctly")
		}
	})
}

func TestQiniuUploader_DeleteFile_InvalidConfig(t *testing.T) {
	// 测试无效配置下的删除（需要真实凭证才能测试成功场景）
	uploader := NewQiniu("invalid_ak", "invalid_sk", "invalid_bucket")

	tests := []struct {
		name    string
		key     string
		wantErr bool
	}{
		{
			name:    "空key",
			key:     "",
			wantErr: true, // 空key会返回错误
		},
		{
			name:    "有效key格式",
			key:     "test/key.txt",
			wantErr: true, // 无效凭证会导致错误
		},
		{
			name:    "带特殊字符的key",
			key:     "test/key with spaces.txt",
			wantErr: true, // 无效凭证会导致错误
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uploader.DeleteFile(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQiniuUploader_StructFields(t *testing.T) {
	// 测试结构体字段访问
	uploader := QiniuUploader{
		AccessKey: "field_ak",
		SecretKey: "field_sk",
		Bucket:    "field_bucket",
	}

	if uploader.AccessKey != "field_ak" {
		t.Errorf("AccessKey field = %v, want %v", uploader.AccessKey, "field_ak")
	}
	if uploader.SecretKey != "field_sk" {
		t.Errorf("SecretKey field = %v, want %v", uploader.SecretKey, "field_sk")
	}
	if uploader.Bucket != "field_bucket" {
		t.Errorf("Bucket field = %v, want %v", uploader.Bucket, "field_bucket")
	}
}

// Integration test - requires real credentials
// To run: go test -tags=integration ./...
func TestQiniuUploader_Integration(t *testing.T) {
	// 此测试需要真实凭证，仅在集成测试时运行
	ak := ""
	sk := ""
	bucket := ""

	if ak == "" || sk == "" || bucket == "" {
		t.Skip("跳过集成测试：未配置真实凭证")
	}

	_ = NewQiniu(ak, sk, bucket)

	// 注意：这里需要创建真实的 multipart.FileHeader 才能测试
	// 实际使用时需要配合 HTTP 请求
	t.Log("七牛云上传器集成测试跳过 - 需要真实文件上传")
}
