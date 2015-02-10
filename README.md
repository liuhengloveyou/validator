validator.cn
================

一个有中国特色的数据格式验证包

安装
===

	go get github.com/liuhengloveyou/validator.cn

然后在自己的代码中引用:

	import (
		"github.com/liuhengloveyou/validator.cn"
	)

用法：
=====

一个简单的示例：

	type NewUserRequest struct {
		Username string `validate:"min=3,max=40,regexp=^[a-zA-Z]$"`
		Name string     `validate:"nonzero"`
		Age int         `validate:"min=21"`
		Password string `validate:"min=8"`
		Email string `validate:"email"`
		IdCard string `validate:"idcard"`
	}

	nur := NewUserRequest{Username: "something", Age: 20}
	if err := validator.Validate(nur); err != nil {
		// values not valid, deal with errors here
	}


内置数据验证格式：

	len
		For numeric numbers, max will simply make sure that the
		value is equal to the parameter given. For strings, it
		checks that the string length is exactly that number of
		characters. For slices,	arrays, and maps, validates the
		number of items. (Usage: len=10)
	
	max
		For numeric numbers, max will simply make sure that the
		value is lesser or equal to the parameter given. For strings,
		it checks that the string length is at most that number of
		characters. For slices,	arrays, and maps, validates the
		number of items. (Usage: max=10)
	
	min
		For numeric numbers, min will simply make sure that the value
		is greater or equal to the parameter given. For strings, it
		checks that the string length is at least that number of
		characters. For slices, arrays, and maps, validates the
		number of items. (Usage: min=10)
	
	nonzero
		This validates that the value is not zero. The appropriate
		zero value is given by the Go spec (e.g. for int it's 0, for
		string it's "", for pointers is nil, etc.) Usage: nonzero
	
	regexp
		Only valid for string types, it will validator that the
		value matches the regular expression provided as parameter.
		(Usage: regexp=^a.*b$)
	email
		验证是否为合法的email字符串格式
	phone
		验证是否为合法的电话号码字符串格式
	cellphone
		验证是否为合法的手机号码字符串格式
	idcard
		验证是否为合法的身份号码字符串格式
	unicn
		验证是否只包含中文Unicode字符

自定义格式验证

用户也可以用 SetValidationFunc 实现自定义的验证. 首先要实现自己的验证函数：

	// Very simple validator
	func notZZ(v interface{}, param string) error {
		st := reflect.ValueOf(v)
		if st.Kind() != reflect.String {
			return errors.New("notZZ only validates strings")
		}
		if st.String() == "ZZ" {
			return errors.New("value cannot be ZZ")
		}
		return nil
	}

然后添加到Validator.cn构架：

	validator.SetValidationFunc("notzz", notZZ)

就可以用"noZZ"标签做验证了：

	type T struct {
		A string  `validate:"nonzero,notzz"`
	}
	t := T{"ZZ"}
	if valid, errs := validator.Validate(t); !valid {
		fmt.Printf("Field A error: %s\n", errs["A"][0])
	}
