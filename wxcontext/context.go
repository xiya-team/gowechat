package wxcontext

import (
	"github.com/xiya-team/gowechat/util"
	"io/ioutil"
	"net/http"
	"sync"
)

// Context struct
type Context struct {
	*Config

	Writer  http.ResponseWriter
	Request *http.Request

	//accessTokenLock 读写锁 同一个AppID一个
	accessTokenLock *sync.RWMutex

	//jsAPITicket 读写锁 同一个AppID一个
	jsAPITicketLock *sync.RWMutex

	HTTPClient  *http.Client
	SHTTPClient *http.Client //SSL client
}

// Query returns the keyed url query value if it exists
func (ctx *Context) Query(key string) string {
	value, _ := ctx.GetQuery(key)
	return value
}

// GetQuery is like Query(), it returns the keyed url query value
func (ctx *Context) GetQuery(key string) (string, bool) {
	req := ctx.Request
	if values, ok := req.URL.Query()[key]; ok && len(values) > 0 {
		return values[0], true
	}
	return "", false
}

// SetJsAPITicketLock 设置jsAPITicket的lock
func (ctx *Context) SetJsAPITicketLock(lock *sync.RWMutex) {
	ctx.jsAPITicketLock = lock
}

// GetJsAPITicketLock 获取jsAPITicket 的lock
func (ctx *Context) GetJsAPITicketLock() *sync.RWMutex {
	return ctx.jsAPITicketLock
}

//InitHTTPClients Context中初始化 httpclient httpsclient
func (ctx *Context) InitHTTPClients() (err error) {
	//create http client
	if ctx.SslCertFilePath != "" && ctx.SslKeyFilePath != "" {
		SslCertContent, _ := ioutil.ReadFile(ctx.SslCertFilePath)
		SslKeyContent, _ := ioutil.ReadFile(ctx.SslKeyFilePath)

		if client, err := util.NewTLSHttpClientFromContent(string(SslCertContent), string(SslKeyContent)); err == nil {
			ctx.SHTTPClient = client
		} else {
			return err
		}
	}

	ctx.HTTPClient = http.DefaultClient
	return
}
