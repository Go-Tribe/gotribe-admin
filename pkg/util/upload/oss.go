package upload

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"mime/multipart"
	"path"
	"strconv"
	"time"
)

// OSSUploader 结构体
type OSSUploader struct {
	Endpoint        string
	AccessKeyId     string
	AccessKeySecret string
	Bucket          string
}

// NewOSS 构造函数
func NewOSS(endpoint, accessKeyId, accessKeySecret, bucket string) OSSUploader {
	return OSSUploader{
		Endpoint:        endpoint,
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
		Bucket:          bucket,
	}
}

// UploadFile OSS上传文件
func (o OSSUploader) UploadFile(file *multipart.FileHeader) (UploadResource, error) {
	client, err := oss.New(o.Endpoint, o.AccessKeyId, o.AccessKeySecret)
	if err != nil {
		return UploadResource{}, err
	}

	bucket, err := client.Bucket(o.Bucket)
	if err != nil {
		return UploadResource{}, err
	}

	src, err := file.Open()
	if err != nil {
		return UploadResource{}, err
	}
	defer src.Close()

	currentTime := time.Now().Format("20060102")
	fileUnixName := strconv.FormatInt(time.Now().UnixNano(), 10)
	fileExt := path.Ext(file.Filename)
	objectName := currentTime + "/" + fileUnixName + fileExt

	err = bucket.PutObject(objectName, src)
	if err != nil {
		return UploadResource{}, err
	}

	//增加返回文件后缀以及名字，通过结构体
	fileRet := UploadResource{
		FileExt: fileExt,
		Key:     objectName,
	}
	return fileRet, nil
}

// DeleteFile 删除文件
func (o OSSUploader) DeleteFile(key string) error {
	client, err := oss.New(o.Endpoint, o.AccessKeyId, o.AccessKeySecret)
	if err != nil {
		return err
	}

	bucket, err := client.Bucket(o.Bucket)
	if err != nil {
		return err
	}

	err = bucket.DeleteObject(key)
	if err != nil {
		return err
	}

	return nil
}
