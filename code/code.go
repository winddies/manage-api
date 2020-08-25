package code

type Code int

const (
	OK          Code = 200
	ResultError Code = 400
	ServerError Code = 500

	CreateFailed Code = 5000
	ErrInvalidId Code = 5001
)

var errCodeMap = map[Code]string{
	CreateFailed: "创建文章失败",
	ServerError:  "服务器错误",
}

func (code Code) Int() int {
	return int(code)
}

func (code Code) Error() string {
	return errCodeMap[code]
}
