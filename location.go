package WebGuard

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var (
	_apiUrl = []string{
		// 中国用户的api接口
		"https://sp0.baidu.com/8aQDcjqpAAV3otqbppnN2DJv/api.php?query=%s&co=&resource_id=6006",
		// 其他地区
		"https://ipapi.co/%s/json/",
	}
)

// LoopUpLocation 查询地理位置是否合法
func LoopUpLocation(locations []string, ip string) bool {
	// 先检查我们的Location的第一个字符是不是英文,如果是英文的话,优先使用国外API
	var (
		isEnglish   = 0
		LocationTag string
		Location    string
		body        string
		allowStatus int
	)
	// 如果我们的地理位置切片长度为0,且ip地址不是空的话
	if len(locations) == 0 && ip != "" {
		return false
	}
	// 遍历全部地址
	for _, location := range locations {
		for _, url := range _apiUrl {
			if isEnglish != 1 {
				if regexp.MustCompile("[a-zA-Z]").MatchString(location[0:1]) {
					// 是英文的话，我们使用国外api
					url, isEnglish = _apiUrl[1], 1
				} else {
					url = _apiUrl[0]
				}
			}
			// 执行查询操作
			allowStatus, body = _httpRequest(fmt.Sprintf(url, ip))
			if allowStatus == 200 {
				// 如果请求的url是国内的
				if url == _apiUrl[0] {
					LocationTag = `data.#.location`
					break
				}
				LocationTag = `city`
				Location += gjson.Get(body, `region`).String()
				break
			}
		}
		// 检查json数据是否符合
		if gjson.Valid(body) {
			resp := gjson.Get(body, LocationTag)
			if resp.Exists() {
				for _, name := range resp.Array() {
					Location += name.String()
				}
				var prettyJSON bytes.Buffer
				// Format output JSON data
				_ = json.Indent(&prettyJSON, []byte(body), "", "\t")
				// Check whether the IP address is the same as the specified location
				if strings.Contains(strings.ToLower(Location), strings.ToLower(location)) {
					return true // The query result is true
				}
			}
		}
	}
	return false
}

func _httpRequest(url string) (int, string) {
	client := resty.New()
	// 请求超时
	client.SetTimeout(8 * time.Second)
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // 不检查tls
	// HTTP request header information
	client.Header = http.Header{
		"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36"},
		"Accept":     {"text/html, application/xhtml+xml, image/jxr, */*"},
		"RedGuard":   {"True"},
		"charset":    {"UTF-8"},
	}
	resp, err := client.R().
		EnableTrace(). // 所触发请求的rest客户端跟踪
		Get(url) // http get请求
	// 检查url请求是否成功
	if err != nil {
		return 0, ""
	}
	return resp.StatusCode(),
		strings.TrimSpace(mahonia.NewDecoder("gbk").ConvertString(string(resp.Body())))
}
