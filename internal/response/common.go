package response

type Response struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data,omitempty"`
	Msg   string      `json:"msg"`
	Error error       `json:"error,omitempty"`
}

func BuildResponse(data interface{}, code int, msg string, err error) Response {
	return Response{
		Code:  code,
		Data:  data,
		Msg:   msg,
		Error: err,
	}
}

// TrackedErrorResponse 有追踪信息的错误响应
type TrackedErrorResponse struct {
	Response
	TrackID string `json:"track_id"`
}

// 三位数错误编码为复用http原本含义
// 五位数错误编码为应用自定义错误
// 五开头的五位数错误编码为服务器端错误，比如数据库操作失败
// 四开头的五位数错误编码为客户端错误，有时候是客户端代码写错了，有时候是用户操作错误
const (
	PumpDecimals = 6
	SolDecimals  = 9
	// Redis key prefixes
	RedisKeyPrefixToken        = "auth:token:"
	RedisKeyPrefixMessage      = "auth:message:"
	RedisKeyPrefixMessageNonce = "auth:message:nonce:"
	RedisKeyPrefixPumpTrending = "token:pump:trending"
	RedisKeyHotTokens          = "tokens:hot:48h"
)

// 通用成功
const CodeSuccess = 200 // 请求成功

// 客户端错误 4xx
const (
	CodeUnauthorized = 401  // 未登录或认证失败
	CodeNoRightErr   = 403  // 权限不足，禁止访问
	CodeParamErr     = 4001 // 参数错误（自定义扩展）
	CodeNotFound     = 4041 // 资源未找到（扩展 404）
)

// 服务端错误 5xx
const (
	CodeDBError       = 5001 // 数据库操作失败
	CodeEncryptError  = 5002 // 加密失败
	CodeCacheError    = 5003 // 缓存操作失败
	CodeServerUnknown = 5000 // 未知的服务器错误
)

// CheckLogin 检查登录
func CheckLogin() Response {
	return Response{
		Code: CodeUnauthorized,
		Msg:  "未登录",
	}
}

// Err 通用错误处理
func Err(errCode int, msg string, err error) Response {
	res := Response{
		Code:  errCode,
		Msg:   msg,
		Error: err,
	}
	return res
}

// DBErr 数据库操作失败
func DBErr(msg string, err error) Response {
	if msg == "" {
		msg = "数据库操作失败"
	}
	return Err(CodeDBError, msg, err)
}

// ParamErr 各种参数错误
func ParamErr(msg string, err error) Response {
	if msg == "" {
		msg = "参数错误"
	}
	return Err(CodeParamErr, msg, err)
}
