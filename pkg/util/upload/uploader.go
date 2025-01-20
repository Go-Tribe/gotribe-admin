package upload

import (
	"mime/multipart"
	"os"
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

// Service 提供上传和删除文件的功能
type Service struct {
	uploader Uploader
}

// NewService 创建一个新的 Service 实例
func NewUploadFile(endpoint, accessKeyId, accessKeySecret, bucketName string, enableOss bool) (*Service, error) {
	var uploader Uploader
	if enableOss {
		uploader = NewOSS(endpoint, accessKeyId, accessKeySecret, bucketName)
	} else {
		uploader = NewQiniu(accessKeyId, accessKeySecret, bucketName)
	}

	return &Service{
		uploader: uploader,
	}, nil
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
