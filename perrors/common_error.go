package perrors

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ServerError   = NewError(500, "服务内部错误")
	InvalidParams = NewError(400, "请求参数错误")
)

type Error struct {
	err  error
	code int
}

func (e *Error) Error() string {
	return e.err.Error()
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) Is(err error) bool {
	if _, ok := err.(*Error); ok {
		return true
	}

	return false
}

// Wrap 包装错误信息
func (e *Error) Wrap(msg string) *Error {
	return &Error{
		err:  errors.New(e.Error() + ":" + msg),
		code: e.code,
	}
}

// WithMessage 替换message
func (e *Error) WithMessage(msg string) *Error {
	return &Error{
		err:  errors.New(msg),
		code: e.code,
	}
}

var codes = map[int]string{}

// NewError returns instance an Error object.
func NewError(code int, msg string) *Error {
	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("错误码 %d 已经存在，请更换一个", code))
	}
	codes[code] = msg
	return &Error{
		err:  errors.New(msg),
		code: code,
	}
}

// GenError returns instance an Error object.
func GenError(code int, msg string) *Error {
	return &Error{
		err:  errors.New(msg),
		code: code,
	}
}

// Trace outputs a stack info if there has error.
func Trace(err interface{}) string {
	if err != nil {
		return ""
	}
	if e, ok := err.(Error); ok {
		return fmt.Sprintf("Err:%+v", e.err)
	}
	if e, ok := err.(error); ok {
		return fmt.Sprintf("err:%+v", e)
	}

	return "unknown error"
}

// AssertErrorNil asserts error is nil.
func AssertErrorNil(err error, throw error) {
	Assert(err == nil, throw)
}

// AssertNotEmpty asserts object is not empty.
func AssertNotEmpty(object interface{}, throw error) {
	Assert(isEmpty(object) == false, throw)
}

// AssertTrue asserts actual is true
func AssertTrue(actual bool, throw error) {
	Assert(actual == true, throw)
}

// AssertFalse asserts actual is false
func AssertFalse(actual bool, throw error) {
	Assert(actual == false, throw)
}

// Assert asserts condition is expected
func Assert(condition bool, throw error) {
	if !condition {
		panic(throw)
	}
}

// isEmpty checks object is empty
func isEmpty(object interface{}) bool {
	if object == nil {
		return true
	}

	objV := reflect.ValueOf(object)

	switch objV.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		return objV.Len() == 0
	case reflect.Ptr:
		if objV.IsNil() {
			return true
		}
		inter := objV.Elem().Interface()
		return isEmpty(inter)
	default:
		zero := reflect.Zero(objV.Type())
		return reflect.DeepEqual(object, zero.Interface())
	}
}
