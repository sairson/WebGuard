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

- 满足用户定义规则

![image](https://user-images.githubusercontent.com/74412075/172836609-bece883b-f59b-4a0e-a64d-216cf72352ed.png)

- 不满足用户定义规则，跳转到用户定义的360.net网址

![image](https://user-images.githubusercontent.com/74412075/172836762-d292db01-08e0-491e-9ae4-f42670cb9131.png)

程序会根据用户在cfg.yml进行规则匹配，当http或https请求流量传入后,会走WebGuard的规则匹配，如果满足则放行继续走我们的handler控制函数，否则的话，会根据配置进行流量drop或者反代到指定域名。以扰乱溯源人员分析
