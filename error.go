package validator

import (
	"errors"
)

var (
	ErrUnknownTag   = errors.New("标签非法")
	ErrInvalid      = errors.New("值非法")
	ErrUnsupported  = errors.New("不支持的验证类型")
	ErrBadParameter = errors.New("错误的验证参数")
	ErrBadName      = errors.New("必须为可导出")

	ErrZeroValue = errors.New("必须非零值")
	ErrMin       = errors.New("小于最小值")
	ErrMax       = errors.New("大于最大值")
	ErrLen       = errors.New("长度不对")
	ErrRegexp    = errors.New("模式不匹配")
	ErrEmail     = errors.New("邮箱地址格式不正确")
	ErrPhone     = errors.New("电话号码格式不正确")
	ErrCellPhone = errors.New("手机号码格式不正确")
	ErrIdCard    = errors.New("身份证号格式不正确")
)
