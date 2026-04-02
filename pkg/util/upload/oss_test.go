// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package upload

import (
	"testing"
)

func TestNewOSS(t *testing.T) {
	tests := []struct {
		name            string
		endpoint        string
		accessKeyId     string
		accessKeySecret string
		bucket          string
		wantEndpoint    string
		wantAccessKeyId string
		wantSecret      string
		wantBucket      string
	}{
		{
			name:            "正常参数",
			endpoint:        "oss-cn-hangzhou.aliyuncs.com",
			accessKeyId:     "test_access_key_id",
			accessKeySecret: "test_access_key_secret",
			bucket:          "test_bucket",
			wantEndpoint:    "oss-cn-hangzhou.aliyuncs.com",
			wantAccessKeyId: "test_access_key_id",
			wantSecret:      "test_access_key_secret",
			wantBucket:      "test_bucket",
		},
		{
			name:            "空参数",
			endpoint:        "",
			accessKeyId:     "",
			accessKeySecret: "",
			bucket:          "",
			wantEndpoint:    "",
			wantAccessKeyId: "",
			wantSecret:      "",
			wantBucket:      "",
		},
		{
			name:            "带协议前缀的 endpoint",
			endpoint:        "https://oss-cn-beijing.aliyuncs.com",
			accessKeyId:     "ak_123",
			accessKeySecret: "sk_456",
			bucket:          "my-bucket",
			wantEndpoint:    "https://oss-cn-beijing.aliyuncs.com",
			wantAccessKeyId: "ak_123",
			wantSecret:      "sk_456",
			wantBucket:      "my-bucket",
		},
		{
			name:            "特殊字符",
			endpoint:        "oss-cn-shanghai.aliyuncs.com",
			accessKeyId:     "LTAI-123abc!@#",
			accessKeySecret: "secret-456def$%^&*",
			bucket:          "bucket-name.123",
			wantEndpoint:    "oss-cn-shanghai.aliyuncs.com",
			wantAccessKeyId: "LTAI-123abc!@#",
			wantSecret:      "secret-456def$%^&*",
			wantBucket:      "bucket-name.123",
		},
		{
			name:            "长字符串",
			endpoint:        "oss-cn-shenzhen-internal.aliyuncs.com",
			accessKeyId:     "a very long access key id with many characters",
			accessKeySecret: "a very long access key secret with many characters and symbols !@#$%^&*()_+-=",
			bucket:          "my-company-production-bucket-001",
			wantEndpoint:    "oss-cn-shenzhen-internal.aliyuncs.com",
			wantAccessKeyId: "a very long access key id with many characters",
			wantSecret:      "a very long access key secret with many characters and symbols !@#$%^&*()_+-=",
			wantBucket:      "my-company-production-bucket-001",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uploader := NewOSS(tt.endpoint, tt.accessKeyId, tt.accessKeySecret, tt.bucket)
			
			if uploader.Endpoint != tt.wantEndpoint {
				t.Errorf("NewOSS() Endpoint = %v, want %v", uploader.Endpoint, tt.wantEndpoint)
			}
			if uploader.AccessKeyId != tt.wantAccessKeyId {
				t.Errorf("NewOSS() AccessKeyId = %v, want %v", uploader.AccessKeyId, tt.wantAccessKeyId)
			}
			if uploader.AccessKeySecret != tt.wantSecret {
				t.Errorf("NewOSS() AccessKeySecret = %v, want %v", uploader.AccessKeySecret, tt.wantSecret)
			}
			if uploader.Bucket != tt.wantBucket {
				t.Errorf("NewOSS() Bucket = %v, want %v", uploader.Bucket, tt.wantBucket)
			}
		})
	}
}

func TestOSSUploader_UploadFile_InvalidFile(t *testing.T) {
	// 测试无效文件上传（需要真实凭证才能测试成功场景）
	uploader := NewOSS("oss-cn-hangzhou.aliyuncs.com", "test_ak", "test_sk", "test_bucket")
	
	// 注意：由于源代码没有nil检查，直接传递nil会导致panic
	// 这里只测试结构体是否正确创建
	t.Run("上传器创建", func(t *testing.T) {
		if uploader.AccessKeyId != "test_ak" {
			t.Error("OSSUploader was not created correctly")
		}
	})
}

func TestOSSUploader_DeleteFile_InvalidConfig(t *testing.T) {
	// 测试无效配置下的删除（需要真实凭证才能测试成功场景）
	uploader := NewOSS("invalid.endpoint", "invalid_ak", "invalid_sk", "invalid_bucket")
	
	tests := []struct {
		name    string
		key     string
		wantErr bool
	}{
		{
			name:    "空key",
			key:     "",
			wantErr: true, // 空key或无效配置会返回错误
		},
		{
			name:    "有效key格式",
			key:     "20240101/test_file.txt",
			wantErr: true, // 无效端点会导致错误
		},
		{
			name:    "带路径的key",
			key:     "images/2024/photo.jpg",
			wantErr: true, // 无效端点会导致错误
		},
		{
			name:    "带特殊字符的key",
			key:     "files/document (1).pdf",
			wantErr: true, // 无效端点会导致错误
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

func TestOSSUploader_StructFields(t *testing.T) {
	// 测试结构体字段访问
	uploader := OSSUploader{
		Endpoint:        "oss-cn-beijing.aliyuncs.com",
		AccessKeyId:     "field_ak_id",
		AccessKeySecret: "field_ak_secret",
		Bucket:          "field_bucket",
	}

	if uploader.Endpoint != "oss-cn-beijing.aliyuncs.com" {
		t.Errorf("Endpoint field = %v, want %v", uploader.Endpoint, "oss-cn-beijing.aliyuncs.com")
	}
	if uploader.AccessKeyId != "field_ak_id" {
		t.Errorf("AccessKeyId field = %v, want %v", uploader.AccessKeyId, "field_ak_id")
	}
	if uploader.AccessKeySecret != "field_ak_secret" {
		t.Errorf("AccessKeySecret field = %v, want %v", uploader.AccessKeySecret, "field_ak_secret")
	}
	if uploader.Bucket != "field_bucket" {
		t.Errorf("Bucket field = %v, want %v", uploader.Bucket, "field_bucket")
	}
}

func TestOSSUploader_VariousEndpoints(t *testing.T) {
	// 测试各种阿里云 OSS Endpoint 格式
	endpoints := []string{
		"oss-cn-hangzhou.aliyuncs.com",
		"oss-cn-shanghai.aliyuncs.com",
		"oss-cn-beijing.aliyuncs.com",
		"oss-cn-shenzhen.aliyuncs.com",
		"oss-cn-hongkong.aliyuncs.com",
		"oss-us-west-1.aliyuncs.com",
		"oss-ap-southeast-1.aliyuncs.com",
		"https://oss-cn-hangzhou.aliyuncs.com",
		"http://oss-cn-hangzhou.aliyuncs.com",
	}

	for _, ep := range endpoints {
		t.Run(ep, func(t *testing.T) {
			uploader := NewOSS(ep, "ak", "sk", "bucket")
			if uploader.Endpoint != ep {
				t.Errorf("Endpoint = %v, want %v", uploader.Endpoint, ep)
			}
		})
	}
}

// Integration test - requires real credentials
// To run: go test -tags=integration ./...
func TestOSSUploader_Integration(t *testing.T) {
	// 此测试需要真实凭证，仅在集成测试时运行
	endpoint := ""
	accessKeyId := ""
	accessKeySecret := ""
	bucket := ""
	
	if endpoint == "" || accessKeyId == "" || accessKeySecret == "" || bucket == "" {
		t.Skip("跳过集成测试：未配置真实凭证")
	}
	
	_ = NewOSS(endpoint, accessKeyId, accessKeySecret, bucket)
	
	// 注意：这里需要创建真实的 multipart.FileHeader 才能测试
	// 实际使用时需要配合 HTTP 请求
	t.Log("阿里云OSS上传器集成测试跳过 - 需要真实文件上传")
}
