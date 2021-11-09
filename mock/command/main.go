package main

import (
	"hellogo/mock"
	"net/http"
)

func main() {
	// 响应头
	header := make(map[string]string, 0)
	header["Access-Control-Allow-Origin"] = "*"
	header["Content-Encoding"] = "gzip"
	header["Content-Type"] = "application/json;charset=UTF-8"

	// 模拟服务信息切片
	var mockInfos []mock.MockInfo

	// 模拟服务信息1
	var mockInfo1 mock.MockInfo
	mockInfo1.Header = header
	mockInfo1.Method = http.MethodGet
	mockInfo1.URL = "/hello"
	mockInfo1.JsonFile = "/Users/sunjiashu/Documents/Develop/github.com/bettersun/hellogo/mock/command/hello.json"

	// 模拟服务信息2
	var mockInfo2 mock.MockInfo
	mockInfo2.Header = header
	mockInfo2.Method = http.MethodGet
	mockInfo2.URL = "/welcome"
	mockInfo2.JsonFile = "/Users/sunjiashu/Documents/Develop/github.com/bettersun/hellogo/mock/command/welcome.json"

	// 添加到模拟服务信息切片
	mockInfos = append(mockInfos, mockInfo1)
	mockInfos = append(mockInfos, mockInfo2)

	// 模拟服务选项
	mock.VerOption.UseMock = true
	mock.VerOption.Port = "9012"
	mock.VerOption.MockInfos = mockInfos

	// 启动模拟服务
	mock.Mock(mock.VerOption)
}
