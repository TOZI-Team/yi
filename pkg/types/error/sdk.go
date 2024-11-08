package yError

import "fmt"

type NotFoundSDKErr struct {
	path string
}

// NewNotFoundSDKErr
//
//	@Description: 无法找到对应SDK
//	@param p SDK位置
//	@return *NotFoundSDKErr
func NewNotFoundSDKErr(p string) *NotFoundSDKErr {
	return &NotFoundSDKErr{path: p}
}

func (e *NotFoundSDKErr) Error() string {
	return fmt.Sprintf("%s not found", e.path)
}
