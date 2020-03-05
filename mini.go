package gowechat

import (
	"github.com/xiya-team/gowechat/mini/auth"
	"gowechat/mini/decrypt"
)

type MiniMgr struct {
	*Wechat
}

//GetAccessToken 获取access_token
func (wc *MiniMgr) GetNewAuth() *auth.Auth {
	return auth.NewAuth(wc.Context)
}

func (wc *MiniMgr) NewDecrypt() *decrypt.Decrypt {
	return decrypt.NewDecrypt(wc.Context)
}
