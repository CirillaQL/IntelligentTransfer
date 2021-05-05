/*
 自定义错误包，对项目中的错误进行统一管理
*/
package errorInfo

import "errors"

var (
	RegisterDriverParamsWrong   = errors.New("Register Driver with wrong params ")
	RegisterDriverInsertDBWrong = errors.New("Register Driver insert db wrong ")

	TableDoesNotExist = errors.New("Table is not exist ")

	GRPCServiceError = errors.New("gRPC Service Failed ")
)
