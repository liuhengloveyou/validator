/*
 有中国特色的数据格式验证包
*/

package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

var defaultValidator *Validator

var validations map[string]ValidationFunc = map[string]ValidationFunc{
	"noneor":    noneor,      // 可为空, 当非空的时候...
	"nonone":    nonone,      // 非空
	"len":       length,      // 长度
	"min":       min,         //最小
	"max":       max,         // 最大
	"regexp":    regex,       // 正则验证
	"email":     email,       // 电子邮箱
	"phone":     phone,       // 电话号
	"cellphone": cellphone,   // 手机号
	"idcard":    idcardCheck, // 身份证
	"region":    regionCheck, // 行政区
	"unicn":     unicodecn,   // 中文字符
}

// 数据验证函数原型
type ValidationFunc func(v interface{}, param string) error

type Validator struct {
	tagName         string
	validationFuncs map[string]ValidationFunc
}

// 一条标签定义
type tagS struct {
	Name  string
	Param string
	Fun   ValidationFunc
}

func init() {
	defaultValidator = &Validator{
		tagName:         "validate",
		validationFuncs: validations,
	}
}

func (p *Validator) setValidationFunc(name string, vfun ValidationFunc) error {
	if name == "" {
		return errors.New("验证标签名不能为空")
	}

	if vfun == nil {
		delete(p.validationFuncs, name) // 删除一个验证功能
	} else {
		p.validationFuncs[name] = vfun
	}

	return nil
}

func (p *Validator) validate(v interface{}) error {
	sv := reflect.ValueOf(v)
	st := reflect.TypeOf(v)
	if sv.Kind() == reflect.Ptr && !sv.IsNil() {
		return p.validate(sv.Elem().Interface())
	}
	if sv.Kind() != reflect.Struct {
		return ErrUnsupported
	}

	for i := 0; i < sv.NumField(); i++ {
		f := sv.Field(i)
		for f.Kind() == reflect.Ptr && !f.IsNil() {
			f = f.Elem() // 处理指针
		}

		tagStr := st.Field(i).Tag.Get(p.tagName) // 取TAG
		if tagStr == "-" {
			continue // 不做验证
		}

		if tagStr == "" {
			if f.Kind() == reflect.Struct {
				return p.validate(f.Interface())
			}
			continue // 没有验证TAG
		}

		name := st.Field(i).Name
		if !unicode.IsUpper(rune(name[0])) {
			continue // 只有可导出变量可验证
		}

		err := p.validateVar(f.Interface(), tagStr)
		if err != nil {
			return fmt.Errorf("%s: %s", name, err.Error())
		}
	}

	return nil
}

func (p *Validator) validateVar(v interface{}, tagStr string) (err error) {
	var tags []tagS

	tags, err = p.parseTags(tagStr)
	if err != nil {
		return
	}

	for i := 0; i < len(tags); i++ {
		if err = tags[i].Fun(v, tags[i].Param); err != nil {
			if (tags[i].Name != "noneor") || (i == len(tags)-1) {
				return
			}
		}
	}

	return nil
}

// 解析出一个字段上的所有验证标签
func (p *Validator) parseTags(t string) (tags []tagS, err error) {
	ts := strings.Split(t, ",")
	tags = make([]tagS, 0, len(ts))

	for i := 0; i < len(ts); i++ {
		tg := tagS{}
		v := strings.SplitN(ts[i], "=", 2)
		tg.Name = strings.Trim(v[0], " ")
		if tg.Name == "" {
			return []tagS{}, ErrUnknownTag
		}
		if len(v) > 1 {
			tg.Param = strings.Trim(v[1], " ")
		}

		var found bool
		if tg.Fun, found = p.validationFuncs[tg.Name]; !found {
			return []tagS{}, ErrUnknownTag
		}

		tags = append(tags, tg)

	}

	return tags, nil
}

////////////////////////////////////////////////////////////////////////////////
// 公共接口
////////////////////////////////////////////////////////////////////////////////

// 设置新的验证函数
func SetValidationFunc(name string, vf ValidationFunc) error {
	return defaultValidator.setValidationFunc(name, vf)
}

// 执行验证
func Validate(v interface{}) error {
	return defaultValidator.validate(v)
}
