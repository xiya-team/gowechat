package util

const (
	// baseURL 微信请求基础URL
	BaseURL = "https://api.weixin.qq.com"
)

// POST 参数
type RequestParams map[string]interface{}

// URL 参数
type RequestQueries map[string]string
