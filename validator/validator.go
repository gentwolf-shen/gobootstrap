package validator

import (
	"reflect"
	"strings"

	"github.com/gentwolf-shen/gin-boost/binding"
	"github.com/gentwolf-shen/gohelper-v2/util"
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	trans ut.Translator
)

func UseLangZh() {
	loc := zh.New()
	uni := ut.New(loc, en.New())
	trans, _ = uni.GetTranslator("zh")

	v, _ := binding.Validator.Engine().(*validator.Validate)
	_ = zhTranslations.RegisterDefaultTranslations(v, trans)
}

func UseLangEn() {
	loc := en.New()
	uni := ut.New(loc, loc)
	trans, _ = uni.GetTranslator("en")
	v, _ := binding.Validator.Engine().(*validator.Validate)
	_ = enTranslations.RegisterDefaultTranslations(v, trans)
}

// TODO
func UseLang(loc locales.Translator) {

}

// TODO
func Register() {

}

func Translate(err error) string {
	if reflect.TypeOf(err).Name() != "ValidationErrors" {
		return err.Error()
	}

	errs := err.(validator.ValidationErrors)
	arr := make([]string, len(errs))
	rows := errs.Translate(trans)
	i := 0
	for _, value := range rows {
		arr[i] = util.ToLowFirst(value)
		i++
	}

	return strings.Join(arr, " ")
}
