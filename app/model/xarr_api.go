package model

type ApiPushGroupListRes struct {
	Data    []string `json:"data"`
	Message string   `json:"message"`
	Code    int      `json:"code"`
	Author  string   `json:"author"`
}
