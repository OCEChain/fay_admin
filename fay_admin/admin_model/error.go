package admin_model

import "github.com/go-errors/errors"

var (
	ParamEmpty = errors.New("参数不能为空")
	SystemFail = errors.New("服务器出错")
)
