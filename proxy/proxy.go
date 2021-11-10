package proxy

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// 全局变量：HTTP服务
// 可不使用全局变量，如果需要手动停止服务时需要。
var server http.Server

// VerOption 全局公共变量：选项
var VerOption Option

// HTTP Header 常量 Content-Encoding
const headerContentEncoding = "Content-Encoding"
const encodingGzip = "gzip"

// Option 选项
type Option struct {
	Port    string `yaml:"port"`    // 端口
	ApiHost string `yaml:"apiHost"` // API主机（IP + Port）
}

// Proxy 代理
//   option: 选项
func Proxy(option Option) {
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
	// 转发请求
	doProxy(w, r, VerOption)
}

// 转发请求
func doProxy(w http.ResponseWriter, r *http.Request, option Option) {
	// 创建一个HttpClient用于转发请求
	cli := &http.Client{}

	// 读取请求的Body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("读取请求体发生错误")
		// 响应状态码
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// 转发的URL
	reqURL := option.ApiHost + r.URL.String()

	// 创建转发用的请求
	reqProxy, err := http.NewRequest(r.Method, reqURL, strings.NewReader(string(body)))
	if err != nil {
		log.Println("创建转发请求发生错误")
		// 响应状态码
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// 转发请求的 Header
	for k, v := range r.Header {
		reqProxy.Header.Set(k, v[0])
	}

	// 发起请求
	responseProxy, err := cli.Do(reqProxy)
	if err != nil {
		log.Println("转发请求发生错误")
		// 响应状态码
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	defer responseProxy.Body.Close()

	// 转发响应的 Header
	for k, v := range responseProxy.Header {
		w.Header().Set(k, v[0])
	}

	// 转发响应的Body数据
	var data []byte

	// 读取转发响应的Body
	data, err = ioutil.ReadAll(responseProxy.Body)
	if err != nil {
		log.Println("读取响应体发生错误")
		// 响应状态码
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// 转发响应的输出数据
	var dataOutput []byte
	// gzip压缩判断
	isGzipped := isGzipped(responseProxy.Header)
	// gzip压缩编码数据
	if isGzipped {
		// 读取后 r.Body 即关闭，无法再次读取
		// 若需要再次读取，需要用读取到的内容再次构建Reader
		resProxyGzippedBody := ioutil.NopCloser(bytes.NewBuffer(data))
		defer resProxyGzippedBody.Close() // 延时关闭

		// gzip Reader
		gr, err := gzip.NewReader(resProxyGzippedBody)
		if err != nil {
			log.Println("创建gzip读取器发生错误")
			// 响应状态码
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		defer gr.Close()

		// 读取gzip对象内容
		dataOutput, err = ioutil.ReadAll(gr)
		if err != nil {
			log.Println("读取gzip对象内容发生错误")
			// 响应状态码
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
	} else { // 非gzip压缩
		dataOutput = data
	}
	// 打印转发响应的Body数据，查看转发响应的响应数据时需要。
	log.Println(string(dataOutput))

	// response的Body不能多次读取，
	// 上面已经被读取过一次，需要重新生成可读取的Body数据。
	resProxyBody := ioutil.NopCloser(bytes.NewBuffer(data))
	defer resProxyBody.Close() // 延时关闭

	// 响应状态码
	w.WriteHeader(responseProxy.StatusCode)
	// 复制转发的响应Body到响应Body
	io.Copy(w, resProxyBody)
}

// gzip压缩判断
func isGzipped(header http.Header) bool {
	if header == nil {
		return false
	}

	contentEncoding := header.Get(headerContentEncoding)
	isGzipped := false
	if strings.Contains(contentEncoding, encodingGzip) {
		isGzipped = true
	}

	return isGzipped
}
