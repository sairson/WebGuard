package WebGuard

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"net/http"
)

// WebGuard web卫兵配置
type WebGuard struct {
	config     string                                       // 配置文件路径
	handler    func(w http.ResponseWriter, r *http.Request) // http处理函数
	hotLoading bool                                         // 热加载
	logger     func(in interface{})
}

func New(cfg string, hotLoading bool, logger func(in interface{}), handler func(w http.ResponseWriter, r *http.Request)) *WebGuard {
	return &WebGuard{
		config:     cfg,
		handler:    handler,
		hotLoading: hotLoading,
		logger:     logger,
	}
}

func (g *WebGuard) RunGuard() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg := viper.New()
		cfg.SetConfigType("yaml")
		cfg.SetConfigFile(g.config)
		if err := cfg.ReadInConfig(); err != nil {
			// 配置文件没有找到,不影响规则,我们直接走规则控制函数
			g.handler(w, r)
			return
		}
		if g.hotLoading == true {
			// 配置热加载
			go func() {
				cfg.WatchConfig()
				cfg.OnConfigChange(func(in fsnotify.Event) {})
			}()
		}
		filter, err := WebGuardFilter(cfg, r)
		if err != nil {
			if g.logger != nil {
				g.logger(err) // 调用debug的logger函数,这个函数可以自定义,但是希望采用func(in interface{})格式输出
			}
		}
		// 规则满足
		if filter == true && err == nil {
			g.handler(w, r)
		} else {
			// 规则不满足,我们看看是丢包还是进行重定向
			var (
				requestDrop = cfg.GetBool("proxy-handler.request-drop")
				redirect    = cfg.GetString("proxy-handler.redirect")
			)
			// 判断重定向还是直接丢包
			proxy, _ := NewProxy(redirect, requestDrop)
			r.URL.Path = "/" // 重定向到根
			proxy.ServeHTTP(w, r)
			return
		}
	}
}
