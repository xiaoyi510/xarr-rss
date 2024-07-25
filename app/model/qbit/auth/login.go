package auth

// ApiAuthLoginReq 登录请求
type ApiAuthLoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ApiAuthLoginRes 登录返回 字符串Ok.
type ApiAuthLoginRes string
