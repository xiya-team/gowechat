package template_message

import "github.com/xiya-team/gowechat/mini/base"

type TemplateMessage struct {
	base.MiniBase
}

// https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/template-message/templateMessage.addTemplate.html
// POST https://api.weixin.qq.com/cgi-bin/wxopen/template/add?access_token=ACCESS_TOKEN
func (c *TemplateMessage) AddTemplate()  {
	
}

// https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/template-message/templateMessage.deleteTemplate.html
// POST https://api.weixin.qq.com/cgi-bin/wxopen/template/del?access_token=ACCESS_TOKEN
func (c *TemplateMessage) DeleteTemplate()  {

}