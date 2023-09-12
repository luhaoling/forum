package controllers

import (
	"fmt"
	"project/models"
	"reflect"
	"strings"

	enTranslations "github.com/go-playground/validator/v10/translations/en"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

// 定义一个全局翻译器 T
var trans ut.Translator

// InitTrans 初始化翻译器
func InitTrans(locale string) (err error) {
	// 修改 gin 框架中的 Validator 引擎属性，实现自定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册一个获取 json tag 的自定义方法(用于响应中的错误展示)
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		// 为 SignUpParam 注册自定义校验方法（用于响应中的错误提示信息）
		v.RegisterStructValidation(SignParamStructLevelValidation, models.ParamSignUp{})

		zhT := zh.New() // 中文翻译器
		enT := en.New() // 英文翻译器

		uni := ut.New(enT, zhT, enT)

		// locale 通常取决于 http 请求头的 'Accept_Language'
		var ok bool
		// 也可以使用 uni.FindTranslator(...) 传入多个 locale 进行查找
		trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", locale)
		}
		fmt.Println(locale)

		// 注册翻译器
		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		}
		return
	}
	return
}

// removeTopStruct 去除提示信息中的结构体名称
func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}

// SignParamStructLevelValidation 自定义 SignUpParam 结构体校验函数
func SignParamStructLevelValidation(sl validator.StructLevel) {
	su := sl.Current().Interface().(models.ParamSignUp)

	if su.Password != su.RePassword {
		// 输出错误提示信息，最后一个参数就是传递的 param
		sl.ReportError(su.RePassword, "re_password", "RePassword", "eqfield", "password")

	}
}
