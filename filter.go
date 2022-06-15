package WebGuard

import (
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// WebGuardFilter 过滤器,判断规则是否可用
// req 请求接收 rules 指定的规则配置路径
func WebGuardFilter(cfg *viper.Viper, req *http.Request) (bool, error) {
	if req == nil {
		return false, fmt.Errorf("request is nil pointer")
	}
	// 配置读取完毕,我们接下来验证规则
	var (
		allowHost        = cfg.GetString("proxy-rules.allow-host")
		allowLocation    = cfg.GetString("proxy-rules.allow-location")
		allowTime        = cfg.GetString("proxy-rules.allow-time")
		allowPath        = cfg.GetString("proxy-rules.allow-path")
		allowIpList      = cfg.GetString("proxy-rules.allow-ip-list")
		allowTokenHeader = cfg.GetStringMapString("proxy-rules.allow-token-header")
		allowUserAgent   = cfg.GetStringSlice("proxy-rules.allow-user-agent")
		refuseIpList     = cfg.GetString("proxy-rules.refuse-ip-list")
		allowMaxLength   = cfg.GetInt("proxy-rules.allow-body-max-length")
		allowMinLength   = cfg.GetInt("proxy-rules.allow-body-min-length")
	)
	if allowHost != "" && allowHost != "*" {
		hosts := strings.Split(allowHost, ",")
		var isHost = false
		for _, h := range hosts {
			if req.Host == h {
				isHost = true
				break
			}
		}
		if isHost == false {
			return false, fmt.Errorf("allow-host rules trigger")
		}
	}
	address := strings.Split(req.RemoteAddr, ":")[0]
	if req.Header.Get("X-Forwarded-For") != "" {
		address = strings.Split(strings.TrimSpace(strings.Split(req.Header.Get("X-Forwarded-For"), ",")[0]), ":")[0]
	}
	// 判断是否在允许的allow-location中
	if allowLocation != "" && allowLocation != "*" {
		if !LoopUpLocation(strings.Split(allowLocation, ","), address) {
			return false, fmt.Errorf("allow-location rules trigger")
		}
	}
	// 判断是否在允许的时间段内
	if allowTime != "" && allowTime != "*" {
		num := strings.Split(allowTime, "-")
		afterTime, _ := time.Parse("15:04", strings.TrimSpace(num[0]))
		beforeTime, _ := time.Parse("15:04", strings.TrimSpace(num[1]))
		nowTime, _ := time.Parse("15:04", strings.TrimSpace(fmt.Sprintf("%d:%d", time.Now().Hour(), time.Now().Minute())))
		if nowTime.After(afterTime) && nowTime.Before(beforeTime) {
		} else {
			return false, fmt.Errorf("allow-time rules trigger")
		}
	}
	// 检查白名单ip列表
	if allowIpList != "" && allowIpList != "*" {
		ipList := strings.Split(allowIpList, ",")
		var newList []string
		for _, ip := range ipList {
			if strings.Contains(ip, "/") || strings.Contains(ip, "-") {
				if ipSlices, err := IPIntoSlices(ip); err != nil {
					return false, fmt.Errorf("allow-ip-list format is error")
				} else {
					newList = append(newList, ipSlices...)
				}
			} else {
				newList = append(newList, ip)
			}
		}
		if _, status := _find(newList, address); status == false {
			return false, fmt.Errorf("allow-ip-list rules trigger")
		}
	}
	// 检查header头中是否包含指定的头
	if len(allowTokenHeader) != 0 {
		for k, v := range allowTokenHeader {
			// 获取到值,但是值不为所设置到的值
			if req.Header.Get(k) != "" {
				if req.Header.Get(k) != v {
					return false, fmt.Errorf("not meeting allow-token-header")
				}
			} else {
				//没有找到这个key的话,我们也返回
				return false, fmt.Errorf("allow-token-header rules trigger")
			}
		}
	}
	// 检查是否包含指定的user-agent
	if !(len(allowUserAgent) == 1 && allowUserAgent[0] == "*") {
		var isUserAgent = false
		for _, ua := range allowUserAgent {
			if req.UserAgent() == ua {
				isUserAgent = true
				break
			}
		}
		if isUserAgent == false {
			return false, fmt.Errorf("allow-user-agent rules trigger")
		}
	}
	// 判断uri是否指定
	if allowPath != "" && allowPath != "*" {
		path := strings.Split(allowPath, ",")
		var isPath = false
		for _, p := range path {
			if strings.Contains(req.URL.Path, p) {
				isPath = true
				break
			}
		}
		if isPath == false {
			return false, fmt.Errorf("allow-path rules trigger")
		}
	}
	// 黑名单判断
	if refuseIpList != "" && refuseIpList != "-" {
		ipList := strings.Split(refuseIpList, ",") // 获取拒绝的
		var newList []string
		for _, ip := range ipList {
			if strings.Contains(ip, "/") || strings.Contains(ip, "-") {
				if ipSlices, err := IPIntoSlices(ip); err != nil {
					return false, fmt.Errorf("refuse-ip-list format is error")
				} else {
					newList = append(newList, ipSlices...)
				}
			} else {
				newList = append(newList, ip)
			}
		}
		if _, status := _find(newList, address); status == true {
			return false, fmt.Errorf("refuse-ip-list rules trigger")
		}
	}
	if allowMaxLength != 0 && allowMinLength != 0 {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return false, fmt.Errorf("parse request body length is error")
		}
		if len(body) < allowMinLength || len(body) > allowMaxLength {
			return false, fmt.Errorf("allow-body-length rules trigger")
		}
	}
	return true, nil
}

func _find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
