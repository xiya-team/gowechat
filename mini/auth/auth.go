package auth

import (
	"github.com/xiya-team/gowechat/mini/commons"
	"github.com/xiya-team/gowechat/mini/base"
	"github.com/xiya-team/gowechat/util"
	"github.com/xiya-team/gowechat/wxcontext"
)

const (
	apiLogin          = "/sns/jscode2session"
	apiGetAccessToken = "/cgi-bin/token"
	apiGetPaidUnionID = "/wxa/getpaidunionid"
)

// LoginResponse 返回给用户的数据
type LoginResponse struct {
	commons.CommonError
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	// 用户在开放平台的唯一标识符
	// 只在满足一定条件的情况下返回
	UnionID string `json:"unionid"`
}

//Pay pay
type Auth struct {
	base.MiniBase
}

//NewOauth 实例化授权信息
func NewAuth(context *wxcontext.Context) *Auth {
	auth := new(Auth)
	auth.Context = context
	return auth
}

//获取小程序全局唯一后台接口调用凭据（access_token）。调用绝大多数后台接口时都需使用 access_token，开发者需要进行妥善保存。
// 文档 https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/access-token/auth.getAccessToken.html
// GET https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=APPID&secret=APPSECRET
func (c *Auth) GetAccessToken(appID, secret string) (res *LoginResponse, err error) {
	queries := util.RequestQueries{
		"appid":      appID,
		"secret":     secret,
		"grant_type": "client_credential",
	}

	api := util.BaseURL + apiGetAccessToken

	url, err := util.EncodeURL(api, queries)
	if err != nil {
		return nil, err
	}

	if err := util.GetJSON(url, res); err != nil {
		return nil, err
	}
	return
}

// 登录凭证校验
// 文档 https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/login/auth.code2Session.html
// GET https://api.weixin.qq.com/sns/jscode2session?appid=APPID&secret=SECRET&js_code=JSCODE&grant_type=authorization_code
// appID 小程序 appID
// secret 小程序的 app secret
// code 小程序登录时获取的 code
// Login 登录凭证校验。通过 wx.login 接口获得临时登录凭证 code 后传到开发者服务器调用此接口完成登录流程。
func (c *Auth) Code2Session(appID, secret, code string) (res *LoginResponse, err error) {
	queries := util.RequestQueries{
		"appid":      appID,
		"secret":     secret,
		"js_code":    code,
		"grant_type": "authorization_code",
	}
	api := util.BaseURL + apiLogin
	url, err := util.EncodeURL(api, queries)
	if err != nil {
		return nil, err
	}

	res = new(LoginResponse)
	if err := util.GetJSON(url, res); err != nil {
		return nil, err
	}
	return
}

// TokenResponse 获取 access_token 成功返回数据
type TokenResponse struct {
	commons.CommonError
	AccessToken string `json:"access_token"` // 获取到的凭证
	ExpiresIn   uint   `json:"expires_in"`   // 凭证有效时间，单位：秒。目前是7200秒之内的值。
}

// GetPaidUnionIDResponse response data
type GetPaidUnionIDResponse struct {
	commons.CommonError
	UnionID string `json:"unionid"`
}

// 用户支付完成后，获取该用户的 UnionId，无需用户授权。
// 文档 https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/user-info/auth.getPaidUnionId.html
// GET https://api.weixin.qq.com/wxa/getpaidunionid?access_token=ACCESS_TOKEN&openid=OPENID
// GetPaidUnionID 用户支付完成后，通过微信支付订单号（transaction_id）获取该用户的 UnionId，
func (c *Auth) GetPaidUnionId(accessToken, openID, transactionID string)(res *GetPaidUnionIDResponse, err error) {
	api := util.BaseURL + apiGetPaidUnionID
	queries := util.RequestQueries{
		"openid":         openID,
		"access_token":   accessToken,
		"transaction_id": transactionID,
	}

	return getPaidUnionIDRequest(api, queries)
}

// GetPaidUnionIDWithMCH 用户支付完成后，通过微信支付商户订单号和微信支付商户号（out_trade_no 及 mch_id）获取该用户的 UnionId，
func GetPaidUnionIDWithMCH(accessToken, openID, outTradeNo, mchID string) (*GetPaidUnionIDResponse, error) {
	api := util.BaseURL + apiGetPaidUnionID
	return getPaidUnionIDWithMCH(accessToken, openID, outTradeNo, mchID, api)
}

func getPaidUnionIDWithMCH(accessToken, openID, outTradeNo, mchID, api string) (*GetPaidUnionIDResponse, error) {
	queries := util.RequestQueries{
		"openid":       openID,
		"mch_id":       mchID,
		"out_trade_no": outTradeNo,
		"access_token": accessToken,
	}

	return getPaidUnionIDRequest(api, queries)
}

func getPaidUnionIDRequest(api string, queries util.RequestQueries) (res *GetPaidUnionIDResponse,err error) {
	url, err := util.EncodeURL(api, queries)
	if err != nil {
		return nil, err
	}

	res = new(GetPaidUnionIDResponse)
	if err := util.GetJSON(url, res); err != nil {
		return nil, err
	}
	return
}
