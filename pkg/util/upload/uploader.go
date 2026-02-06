package upload

import (
	"errors"
	"mime/multipart"
	"os"
	"strings"
)

// Uploader 定义上传接口
type Uploader interface {
	UploadFile(file *multipart.FileHeader) (UploadResource, error)
	DeleteFile(key string) error
}

type UploadResource struct {
	FileExt string
	Key     string
	Domain  string
}

// 支持的上传服务商
const (
	ProviderQiniu = "qiniu"
	ProviderOSS   = "oss"
	ProviderS3    = "s3"
)

// Service 提供上传和删除文件的功能
type Service struct {
	uploader Uploader
}

// NewService 根据 provider 创建对应的上传服务实例
// provider: qiniu(七牛) / oss(阿里云OSS) / s3(亚马逊S3，预留，当前返回未实现错误)
func NewService(provider, endpoint, accessKeyId, accessKeySecret, bucketName string) (*Service, error) {
	provider = strings.ToLower(strings.TrimSpace(provider))
	var uploader Uploader
	switch provider {
	case ProviderOSS:
		uploader = NewOSS(endpoint, accessKeyId, accessKeySecret, bucketName)
	case ProviderQiniu:
		uploader = NewQiniu(accessKeyId, accessKeySecret, bucketName)
	case ProviderS3:
		return nil, errors.New("s3 上传尚未实现，请先使用 qiniu 或 oss")
	default:
		return nil, errors.New("不支持的上传服务商: " + provider + "，可选: qiniu, oss, s3")
	}
	return &Service{uploader: uploader}, nil
}

// NewUploadFile 兼容旧用法：根据 enableOss 选择 oss 或 qiniu。新代码请使用 NewService(provider, ...)
func NewUploadFile(endpoint, accessKeyId, accessKeySecret, bucketName string, enableOss bool) (*Service, error) {
	provider := ProviderQiniu
	if enableOss {
		provider = ProviderOSS
	}
	return NewService(provider, endpoint, accessKeyId, accessKeySecret, bucketName)
}

// UploadFile 公用上传文件方法
func (s *Service) UploadFile(file *multipart.FileHeader) (UploadResource, error) {
	if s.uploader == nil {
		return UploadResource{}, os.ErrInvalid
	}
	result, err := s.uploader.UploadFile(file)
	if err != nil {
		return UploadResource{}, err
	}
	return result, nil
}

// DeleteFile 公用删除文件方法
func (s *Service) DeleteFile(key string) error {
	if s.uploader == nil {
		return os.ErrInvalid
	}
	err := s.uploader.DeleteFile(key)
	if err != nil {
		return err
	}
	return nil
}
