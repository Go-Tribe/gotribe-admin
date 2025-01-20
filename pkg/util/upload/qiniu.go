package upload

import (
	"context"
	"mime/multipart"
	"path"
	"strconv"
	"time"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

// QiniuUploader 结构体
type QiniuUploader struct {
	AccessKey string
	SecretKey string
	Bucket    string
}

// NewQiniu 构造函数
func NewQiniu(ak, sk, bucket string) QiniuUploader {
	return QiniuUploader{
		AccessKey: ak,
		SecretKey: sk,
		Bucket:    bucket,
	}
}

// UploadFile 七牛上传文件
func (q QiniuUploader) UploadFile(file *multipart.FileHeader) (UploadResource, error) {
	src, err := file.Open()
	if err != nil {
		return UploadResource{}, err
	}
	defer src.Close()

	putPolicy := storage.PutPolicy{Scope: q.Bucket}
	mac := qbox.NewMac(q.AccessKey, q.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	// 配置参数
	cfg := storage.Config{
		Zone:          &storage.ZoneHuanan, // 华南区
		UseCdnDomains: false,
		UseHTTPS:      false, // 非https
	}

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}        // 上传后返回的结果
	putExtra := storage.PutExtra{} // 额外参数

	// 上传 自定义key，可以指定上传目录及文件名和后缀，
	currentTime := time.Now().Format("20060102")
	fileUnixName := strconv.FormatInt(time.Now().UnixNano(), 10)
	fileExt := path.Ext(file.Filename)
	key := currentTime + "/" + fileUnixName + fileExt // 上传路径，如果当前目录中已存在相同文件，则返回上传失败错误

	err = formUploader.Put(context.Background(), &ret, upToken, key, src, file.Size, &putExtra)
	if err != nil {
		return UploadResource{}, err
	}
	//增加返回文件后缀以及名字，通过结构体
	fileRet := UploadResource{
		FileExt: fileExt,
		Key:     key,
	}
	return fileRet, nil
}

// DeleteFile 删除文件
func (q QiniuUploader) DeleteFile(key string) error {
	mac := qbox.NewMac(q.AccessKey, q.SecretKey)
	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: false,
	}
	bucketManager := storage.NewBucketManager(mac, &cfg)
	err := bucketManager.Delete(q.Bucket, key)
	if err != nil {
		return err
	}
	return nil
}
