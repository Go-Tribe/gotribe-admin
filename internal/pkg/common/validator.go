// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package common

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	ch_translations "github.com/go-playground/validator/v10/translations/zh"
	"regexp"
)

// 全局Validate数据校验实列
var Validate *validator.Validate

// 全局翻译器（保持兼容：默认中文）
var Trans ut.Translator

// 双语翻译器
var TransZH ut.Translator
var TransEN ut.Translator

// 初始化Validator数据校验
func InitValidate() {
	zhLoc := zh.New()
	enLoc := en.New()
	// universal translator: first param is fallback, subsequent are supported
	uni := ut.New(enLoc, zhLoc, enLoc)

	// translators
	transZH, _ := uni.GetTranslator("zh")
	transEN, _ := uni.GetTranslator("en")
	TransZH = transZH
	TransEN = transEN
	// keep legacy default as Chinese
	Trans = TransZH

	Validate = validator.New()
	_ = ch_translations.RegisterDefaultTranslations(Validate, TransZH)
	_ = en_translations.RegisterDefaultTranslations(Validate, TransEN)
	_ = Validate.RegisterValidation("checkMobile", checkMobile)
	Log.Infof("初始化validator.v10数据校验器完成")
}

func checkMobile(fl validator.FieldLevel) bool {
	reg := `^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(fl.Field().String())
}

// GetTransByLang 返回指定语言的翻译器（默认中文）
func GetTransByLang(lang string) ut.Translator {
	if lang == "en" {
		return TransEN
	}
	return TransZH
}

// GetTransFromCtx 根据请求上下文返回翻译器（默认中文）
func GetTransFromCtx(c *gin.Context) ut.Translator {
	if c == nil {
		return TransZH
	}
	if v, ok := c.Get("lang"); ok {
		if s, sok := v.(string); sok && s == "en" {
			return TransEN
		}
	}
	return TransZH
}
