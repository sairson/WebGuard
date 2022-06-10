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
}

func New(cfg string, hotLoading bool, handler func(w http.ResponseWriter, r *http.Request)) *WebGuard {
	return &WebGuard{
		config:     cfg,
		handler:    handler,
		hotLoading: hotLoading,
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
		filter, _ := WebGuardFilter(cfg, r)
		// 规则满足
		if filter == true {
			g.handler(w, r)
		} else {
			// 规则不满足,我们看看是丢包还是进行重定向
			var (
				requestDrop = cfg.GetBool("proxy-handler.request-drop")
				redirect    = cfg.GetString("proxy-handler.redirect")
			)
			// 判断重定向还是直接丢包
			proxy, _ := NewProxy(redirect, requestDrop)
			proxy.ServeHTTP(w, r)
		}
	}
}
