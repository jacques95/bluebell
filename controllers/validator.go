package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTrs "github.com/go-playground/validator/v10/translations/en"
	zhTrs "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"strings"
)

var trans ut.Translator

func InitTrans(local string) (err error) {
	// 初始化翻译引擎，并将其注册到要绑定的校验器上
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		// 截去除 “json tag” 外的多余字段
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "_" {
				return ""
			}
			return name
		})

		zhT := zh.New()
		enT := en.New()
		uni := ut.New(enT, zhT, enT)

		var ok bool

		trans, ok = uni.GetTranslator(local)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", local)
		}

		switch local {
		case "en":
			err = enTrs.RegisterDefaultTranslations(v, trans)
		case "zh":
			err = zhTrs.RegisterDefaultTranslations(v, trans)
		default:
			err = enTrs.RegisterDefaultTranslations(v, trans)
		}
		return
	}
	return
}

// 截去提示信息中的结构体名称
func removeTopStruct(flds map[string]string) map[string]string {
	res := map[string]string{}

	for fld, err := range flds {
		res[fld[strings.Index(fld, ".")+1:]] = err
	}
	return res
}
