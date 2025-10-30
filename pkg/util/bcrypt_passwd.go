// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Password 密码工具类，用于处理密码加密和验证
type Password struct{}

// GenPasswd 密码加密 使用自适应hash算法, 不可逆
func (p *Password) GenPasswd(passwd string) (string, error) {
	hashPasswd, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("密码加密失败: %v", err)
	}
	return string(hashPasswd), nil
}

// ComparePasswd 通过比较两个字符串hash判断是否出自同一个明文
// hashPasswd 需要对比的密文
// passwd 明文
func (p *Password) ComparePasswd(hashPasswd string, passwd string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashPasswd), []byte(passwd)); err != nil {
		return err
	}
	return nil
}

// Encrypt 使用 bcrypt 加密纯文本
func (p *Password) Encrypt(source string) (string, error) {
	return p.GenPasswd(source)
}

// Compare 比较密文和明文是否相同
func (p *Password) Compare(hashedPassword, password string) error {
	return p.ComparePasswd(hashedPassword, password)
}

// 全局实例
var PasswordUtil = &Password{}
