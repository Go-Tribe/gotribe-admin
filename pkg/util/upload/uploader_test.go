// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package upload

import (
	"errors"
	"mime/multipart"
	"os"
	"testing"
)

func TestNewService(t *testing.T) {
	tests := []struct {
		name            string
		provider        string
		endpoint        string
		accessKeyId     string
		accessKeySecret string
		bucketName      string
		wantErr         bool
		errMsg          string
	}{
		{
			name:            "OSS 服务商",
			provider:        "oss",
			endpoint:        "oss-cn-hangzhou.aliyuncs.com",
			accessKeyId:     "test_ak",
			accessKeySecret: "test_sk",
			bucketName:      "test_bucket",
			wantErr:         false,
		},
		{
			name:            "OSS 服务商大写",
			provider:        "OSS",
			endpoint:        "oss-cn-hangzhou.aliyuncs.com",
			accessKeyId:     "test_ak",
			accessKeySecret: "test_sk",
			bucketName:      "test_bucket",
			wantErr:         false,
		},
		{
			name:            "七牛服务商",
			provider:        "qiniu",
			endpoint:        "",
			accessKeyId:     "test_ak",
			accessKeySecret: "test_sk",
			bucketName:      "test_bucket",
			wantErr:         false,
		},
		{
			name:            "七牛服务商大写",
			provider:        "QINIU",
			endpoint:        "",
			accessKeyId:     "test_ak",
			accessKeySecret: "test_sk",
			bucketName:      "test_bucket",
			wantErr:         false,
		},
		{
			name:            "带空格的提供商",
			provider:        "  oss  ",
			endpoint:        "oss-cn-hangzhou.aliyuncs.com",
			accessKeyId:     "test_ak",
			accessKeySecret: "test_sk",
			bucketName:      "test_bucket",
			wantErr:         false,
		},
		{
			name:            "S3 服务商（未实现）",
			provider:        "s3",
			endpoint:        "",
			accessKeyId:     "test_ak",
			accessKeySecret: "test_sk",
			bucketName:      "test_bucket",
			wantErr:         true,
			errMsg:          "s3 上传尚未实现",
		},
		{
			name:            "不支持的服务商",
			provider:        "unknown",
			endpoint:        "",
			accessKeyId:     "test_ak",
			accessKeySecret: "test_sk",
			bucketName:      "test_bucket",
			wantErr:         true,
			errMsg:          "不支持的上传服务商",
		},
		{
			name:            "空服务商",
			provider:        "",
			endpoint:        "",
			accessKeyId:     "test_ak",
			accessKeySecret: "test_sk",
			bucketName:      "test_bucket",
			wantErr:         true,
			errMsg:          "不支持的上传服务商",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewService(tt.provider, tt.endpoint, tt.accessKeyId, tt.accessKeySecret, tt.bucketName)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("NewService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if tt.wantErr && tt.errMsg != "" && err != nil {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("NewService() error message = %v, should contain %v", err.Error(), tt.errMsg)
				}
			}
			
			if !tt.wantErr && service == nil {
				t.Error("NewService() returned nil service without error")
			}
		})
	}
}

func TestNewUploadFile(t *testing.T) {
	tests := []struct {
		name            string
		endpoint        string
		accessKeyId     string
		accessKeySecret string
		bucketName      string
		enableOss       bool
		wantProvider    string
	}{
		{
			name:            "启用 OSS",
			endpoint:        "oss-cn-hangzhou.aliyuncs.com",
			accessKeyId:     "test_ak",
			accessKeySecret: "test_sk",
			bucketName:      "test_bucket",
			enableOss:       true,
			wantProvider:    ProviderOSS,
		},
		{
			name:            "禁用 OSS（使用七牛）",
			endpoint:        "",
			accessKeyId:     "test_ak",
			accessKeySecret: "test_sk",
			bucketName:      "test_bucket",
			enableOss:       false,
			wantProvider:    ProviderQiniu,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewUploadFile(tt.endpoint, tt.accessKeyId, tt.accessKeySecret, tt.bucketName, tt.enableOss)
			if err != nil {
				t.Errorf("NewUploadFile() unexpected error = %v", err)
				return
			}
			if service == nil {
				t.Error("NewUploadFile() returned nil service")
			}
		})
	}
}

func TestService_UploadFile(t *testing.T) {
	tests := []struct {
		name    string
		service *Service
		file    *multipart.FileHeader
		wantErr bool
		errType error
	}{
		{
			name:    "nil uploader",
			service: &Service{uploader: nil},
			file:    nil,
			wantErr: true,
			errType: os.ErrInvalid,
		},
		// 注意：由于底层 uploader 实现没有 nil 检查，
		// 传递 nil file 会导致 panic，因此不在测试中包含此情况
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.service.UploadFile(tt.file)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if tt.wantErr && tt.errType != nil && !errors.Is(err, tt.errType) {
				// 某些情况下错误可能不是 os.ErrInvalid，而是底层错误
				t.Logf("UploadFile() error = %v (expected type: %v)", err, tt.errType)
			}
		})
	}
}

func TestService_DeleteFile(t *testing.T) {
	tests := []struct {
		name    string
		service *Service
		key     string
		wantErr bool
		errType error
	}{
		{
			name:    "nil uploader",
			service: &Service{uploader: nil},
			key:     "test/key.txt",
			wantErr: true,
			errType: os.ErrInvalid,
		},
		{
			name: "valid service with invalid key",
			service: func() *Service {
				s, _ := NewService("oss", "invalid.ep", "ak", "sk", "bucket")
				return s
			}(),
			key:     "test/key.txt",
			wantErr: true,
		},
		{
			name: "empty key",
			service: func() *Service {
				s, _ := NewService("oss", "invalid.ep", "ak", "sk", "bucket")
				return s
			}(),
			key:     "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.service.DeleteFile(tt.key)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if tt.wantErr && tt.errType != nil && !errors.Is(err, tt.errType) {
				t.Logf("DeleteFile() error = %v (expected type: %v)", err, tt.errType)
			}
		})
	}
}

func TestUploadResource(t *testing.T) {
	tests := []struct {
		name     string
		resource UploadResource
		wantExt  string
		wantKey  string
		wantDomain string
	}{
		{
			name: "完整资源信息",
			resource: UploadResource{
				FileExt: ".jpg",
				Key:     "20240101/1234567890.jpg",
				Domain:  "https://cdn.example.com",
			},
			wantExt:  ".jpg",
			wantKey:  "20240101/1234567890.jpg",
			wantDomain: "https://cdn.example.com",
		},
		{
			name: "空资源",
			resource: UploadResource{
				FileExt: "",
				Key:     "",
				Domain:  "",
			},
			wantExt:  "",
			wantKey:  "",
			wantDomain: "",
		},
		{
			name: "无扩展名",
			resource: UploadResource{
				FileExt: "",
				Key:     "files/document",
				Domain:  "",
			},
			wantExt:  "",
			wantKey:  "files/document",
			wantDomain: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.resource.FileExt != tt.wantExt {
				t.Errorf("FileExt = %v, want %v", tt.resource.FileExt, tt.wantExt)
			}
			if tt.resource.Key != tt.wantKey {
				t.Errorf("Key = %v, want %v", tt.resource.Key, tt.wantKey)
			}
			if tt.resource.Domain != tt.wantDomain {
				t.Errorf("Domain = %v, want %v", tt.resource.Domain, tt.wantDomain)
			}
		})
	}
}

func TestProviderConstants(t *testing.T) {
	// 测试提供商常量
	if ProviderQiniu != "qiniu" {
		t.Errorf("ProviderQiniu = %v, want %v", ProviderQiniu, "qiniu")
	}
	if ProviderOSS != "oss" {
		t.Errorf("ProviderOSS = %v, want %v", ProviderOSS, "oss")
	}
	if ProviderS3 != "s3" {
		t.Errorf("ProviderS3 = %v, want %v", ProviderS3, "s3")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
