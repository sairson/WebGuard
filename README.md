# WebGuard
## WebGuard是根据风起师傅的RedGuard和mgeeky师傅的RedWarden结合出来的http请求过滤器go包，亦在帮助采用go编写C2 http监听器做流量过滤和规则匹配

- 例子:

```
package main

import (
	"fmt"
	"github.com/sairson/WebGuard"
	"net/http"
)

func main() {
	guard := WebGuard.New("cfg.yml", true, handler)
	http.HandleFunc("/", guard.RunGuard())
	http.ListenAndServe(":8555", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ok")
}
```

程序会根据用户在cfg.yml进行规则匹配，当http或https请求流量传入后,会走WebGuard的规则匹配，如果满足则放行继续走我们的handler控制函数，否则的话，会根据配置进行流量drop或者反代到指定域名
