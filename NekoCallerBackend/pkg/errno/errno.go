package errno

import (
	"errors"
	"fmt"
)

// ========================== Error Codes ==========================

const (
	SuccessCode    = 0
	ServiceErrCode = iota + 10001 // 服务级错误
	ParamErrCode                  // 参数错误

	CreateClassErrCode          // 创建班级失败
	CreateClassOrStudentErrCode // 创建班级或学生失败
	StudentJoinClassErrCode     // 学生加入班级失败

	QueryStudentIDsErrCode // 查询学生ID失败

	StudentNotFoundErrCode // 学生未找到
)

// ========================== Error Messages ==========================

const (
	SuccessMsg    = "Success"
	ServiceErrMsg = "Service is unable to start successfully"
	ParamErrMsg   = "Wrong Parameter has been given"
)

// ErrNo defines a custom error type.
type ErrNo struct {
	ErrCode int32
	ErrMsg  string
}

// Error makes it compatible with the `error` interface.
func (e ErrNo) Error() string {
	return fmt.Sprintf("err_code=%d, err_msg=%s", e.ErrCode, e.ErrMsg)
}

// NewErrNo creates a new ErrNo.
func NewErrNo(code int32, msg string) ErrNo {
	return ErrNo{code, msg}
}

// WithMessage allows chaining to modify the error message.
func (e ErrNo) WithMessage(msg string) ErrNo {
	e.ErrMsg = msg
	return e
}

// ========================== Predefined Errors ==========================

var (
	// Common errors
	Success    = NewErrNo(SuccessCode, SuccessMsg)
	ServiceErr = NewErrNo(ServiceErrCode, ServiceErrMsg)
	ParamErr   = NewErrNo(ParamErrCode, ParamErrMsg)

	// Class related errors
	CreateClassErr = NewErrNo(CreateClassErrCode, "Create class failed")
	CreateClassOrStudentErr = NewErrNo(CreateClassOrStudentErrCode, "Create class or student failed")
	StudentJoinClassErr = NewErrNo(StudentJoinClassErrCode, "Student join class failed")

	// Enrollment related errors
	QueryStudentIDsErr = NewErrNo(QueryStudentIDsErrCode, "Query student IDs failed")

	// Roll call related errors
	StudentNotFoundErr = NewErrNo(StudentNotFoundErrCode, "Student not found")
)

func ConvertErr(err error) ErrNo {
	Err := ErrNo{}
	if errors.As(err, &Err) {
		return Err
	}

	s := ServiceErr
	s.ErrMsg = err.Error()
	return s
}