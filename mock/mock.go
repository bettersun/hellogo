package mock

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// 全局变量：HTTP服务
var server http.Server

// VerOption 全局公共变量：模拟服务选项
var VerOption Option

// HTTP Header 常量 Content-Encoding
const headerContentEncoding = "Content-Encoding"

// Option 模拟服务选项
type Option struct {
	Port      string     `yaml:"port"`         // 端口
	UseMock   bool       `yaml:"useMock"`      // true: 使用模拟服务 false: 不使用
	MockInfos []MockInfo `yaml:"requestInfos"` // 模拟服务信息
}

// MockInfo 模拟服务信息
type MockInfo struct {
	Method   string            `yaml:"method"`   // 请求方法
	URL      string            `yaml:"url"`      // URL
	Header   map[string]string `yaml:"header"`   // 响应头
	JsonFile string            `yaml:"jsonFile"` // 响应内容的 json 文件
}

// Mock 模拟服务
//   option: 模拟服务选项
func Mock(option Option) {
	log.Println("mock")

	// 监听端口
	port := ":" + option.Port
	server = http.Server{
		Addr: port,
	}

	log.Println(fmt.Sprintf("服务运行中... 端口[%v]", option.Port))
	http.ListenAndServe(port, http.HandlerFunc(DoHandle))
}

// DoHandle 响应函数
func DoHandle(w http.ResponseWriter, r *http.Request) {
	if VerOption.UseMock {
		doMock(w, r, VerOption)
	} else {
		// 转发到真正的API
	}
}

// 模拟服务
func doMock(w http.ResponseWriter, r *http.Request, option Option) {
	statusCode := http.StatusOK
	var info MockInfo

	// 模拟服务通过 请求方法和 URL 来匹配
	// 当 请求 的 请求方法和URL 与 模拟服务选项 的 模拟服务信息 的 请求方法和URL 一致时，
	// 使用 模拟服务选项 的 模拟服务信息 的 Json 文件返回响应内容。
	isMatch := false
	reqMethodUrl := fmt.Sprintf("%s_%s", r.Method, r.URL.String())
	for _, item := range option.MockInfos {
		infoMethodUrl := fmt.Sprintf("%s_%s", item.Method, item.URL)
		if reqMethodUrl == infoMethodUrl {
			isMatch = true
			info = item
			break
		}
	}

	// 没有匹配的模拟服务，返回404
	if !isMatch {
		statusCode = http.StatusNotFound
		w.WriteHeader(statusCode)
	}

	// 响应头
	for k, v := range info.Header {
		w.Header().Set(k, v)
	}

	// 响应文件
	exist := false
	_, err := os.Stat(info.JsonFile)
	if err == nil || os.IsExist(err) {
		exist = true
	}

	// json 文件不存在
	if !exist {
		log.Printf("IsExist Error: %v\n", err)
		statusCode = http.StatusInternalServerError
	}

	// 读取 json 文件
	b, err := ioutil.ReadFile(info.JsonFile)
	if err != nil {
		log.Printf("ReadFile Error: %v\n", err)
		statusCode = http.StatusInternalServerError
	}

	// 响应状态码，必须放在w.Header().Set(k, v)之后
	w.WriteHeader(statusCode)

	isGzip := false
	contentEncoding, ok := info.Header[headerContentEncoding]
	if ok {
		if strings.Contains(contentEncoding, "gzip") {
			isGzip = true
		}
	}

	// 响应
	// gzip 压缩（不同的压缩需要不同的处理对应）
	if isGzip {
		//gzip压缩
		buffer := new(bytes.Buffer)
		gw := gzip.NewWriter(buffer)
		// 写入 json 文件的字节
		gw.Write(b)
		// 需要 Flush()
		gw.Flush()

		w.Write(buffer.Bytes())
	} else {
		w.Write(b)
	}
}
