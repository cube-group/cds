package errors

import (
    "fmt"
    "alex/utils"
)

//自营错误接口
type IMyError interface {
    String() string
    Error() string
    Msg() string
    Code() int
}

//自营错误基类
type MyError struct {
    code int
    msg  string
}

//新建错误对象，传入code和错误信息
func NewCodeErr(code int, msg ...interface{}) IMyError {
    return &MyError{code:code, msg:utils.StringJoin("\n", msg...)}
}

func (this *MyError)String() string {
    return fmt.Sprintf("[err]code: %v msg: %v", this.code, this.msg)
}

func (this *MyError)Error() string {
    return this.String()
}

func (this *MyError)Code() int {
    return this.code
}

func (this *MyError)Msg() string {
    return this.msg
}
