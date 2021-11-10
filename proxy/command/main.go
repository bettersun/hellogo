package main

import (
	"hellogo/proxy"
)

func main() {
	// 转发请求
	//proxy.VerOption.ApiHost = "http://www.baidu.com"
	proxy.VerOption.ApiHost = "http://127.0.0.1:8012"
	proxy.VerOption.Port = "9009"

	// 启动服务
	proxy.Proxy(proxy.VerOption)
}
