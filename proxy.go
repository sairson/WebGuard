package WebGuard

import (
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewProxy(proxyURL string, dropType bool) (*httputil.ReverseProxy, error) {
	destinationURL, err := url.Parse(proxyURL)
	if err != nil {
		return nil, err
	}
	// 代理服务器转发地址
	proxy := httputil.NewSingleHostReverseProxy(destinationURL)
	// dropType检查对请求的响应是否改变
	proxy.ErrorLog = nil
	if dropType == true {
		proxy.ModifyResponse = modifyResponse() // 修改对请求的响应
	}
	return proxy, nil
}

func modifyResponse() func(*http.Response) error {
	return func(resp *http.Response) error {
		defer func(Body io.ReadCloser) {
			_ = Body.Close() // 直接关闭响应
		}(resp.Body)
		return nil
	}
}
