// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
)

// RSA RSA工具类，用于处理RSA加密解密操作
type RSA struct{}

// ReadKeyFromFile 从文件中读取RSA key
func (r *RSA) ReadKeyFromFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("无法打开文件 %s: %v", filename, err)
	}
	defer f.Close()

	fileInfo, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("无法获取文件信息: %v", err)
	}

	b := make([]byte, fileInfo.Size())
	_, err = f.Read(b)
	if err != nil {
		return nil, fmt.Errorf("无法读取文件内容: %v", err)
	}
	return b, nil
}

// Encrypt RSA加密
func (r *RSA) Encrypt(data, publicBytes []byte) ([]byte, error) {
	var res []byte
	// 解析公钥
	block, _ := pem.Decode(publicBytes)
	if block == nil {
		return res, fmt.Errorf("无法加密, 公钥可能不正确")
	}

	// 使用X509将解码之后的数据 解析出来
	// x509.MarshalPKCS1PublicKey(block):解析之后无法用，所以采用以下方法：ParsePKIXPublicKey
	keyInit, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return res, fmt.Errorf("无法加密, 公钥可能不正确, %v", err)
	}
	// 使用公钥加密数据
	pubKey := keyInit.(*rsa.PublicKey)
	res, err = rsa.EncryptPKCS1v15(rand.Reader, pubKey, data)
	if err != nil {
		return res, fmt.Errorf("无法加密, 公钥可能不正确, %v", err)
	}
	// 将数据加密为base64格式
	return []byte(r.EncodeStr2Base64(string(res))), nil
}

// Decrypt 对数据进行解密操作
func (r *RSA) Decrypt(base64Data, privateBytes []byte) ([]byte, error) {
	var res []byte
	// 将base64数据解析
	decodedStr, err := r.DecodeStrFromBase64(string(base64Data))
	if err != nil {
		return res, fmt.Errorf("base64解码失败: %v", err)
	}
	data := []byte(decodedStr)
	// 解析私钥
	block, _ := pem.Decode(privateBytes)
	if block == nil {
		return res, fmt.Errorf("无法解密, 私钥可能不正确")
	}
	// 还原数据
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return res, fmt.Errorf("无法解密, 私钥可能不正确, %v", err)
	}
	res, err = rsa.DecryptPKCS1v15(rand.Reader, privateKey, data)
	if err != nil {
		return res, fmt.Errorf("无法解密, 私钥可能不正确, %v", err)
	}
	return res, nil
}

// EncodeStr2Base64 加密base64字符串
func (r *RSA) EncodeStr2Base64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// DecodeStrFromBase64 解密base64字符串
func (r *RSA) DecodeStrFromBase64(str string) (string, error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", fmt.Errorf("base64解码失败: %v", err)
	}
	return string(decodeBytes), nil
}

// 全局实例
var RSAUtil = &RSA{}
