package base

import (
	"github.com/xiya-team/gowechat/wxcontext"
)

//MchBase base mch
type MiniBase struct {
	*wxcontext.Context
}

//PostXML postXML
func (c *MiniBase) PostXML(url string, req map[string]string) (resp map[string]string, err error) {


	return
}
