package gowechat

import (
	"github.com/xiya-team/gowechat/mini/auth"
)

type MiniMgr struct {
	*Wechat
}

//GetAccessToken 获取access_token
func (c *MiniMgr) GetAccessToken() *auth.Auth {
	return auth.NewAuth(c.Context)
}