// @Title
// @Description $
// @Author  55
// @Date  2022/5/12
package httputil

import (
	"net"
	"net/http"
)

const UrlBase = "/"

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

// 获取IP
// clientIP 尽最大努力实现获取客户端 IP 的算法。
// 解析 X-Real-IP 和 X-Forwarded-For 以便于反向代理（nginx 或 haproxy）可以正常工作。
func GetIP(r *http.Request) string {
	remoteAddr := r.RemoteAddr
	if ip := r.Header.Get("XRealIP"); ip != "" {
		remoteAddr = ip
	} else if ip = r.Header.Get("XForwardedFor"); ip != "" {
		remoteAddr = ip
	} else if ip = r.Header.Get("X-Forwarded-For"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}

	return remoteAddr
}

// 路由结构体对象定义
type Router struct {
	Route map[string]HandlerFunc
}

// 路由表初始化
func (r *Router) HandleFunc(path string, f HandlerFunc) {
	if r.Route == nil {
		r.Route = make(map[string]HandlerFunc)
	}
	r.Route[path] = f
}

// http回调
func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// 正常路由
	if f, ok := r.Route[req.URL.Path]; ok {
		f(res, req)
	} else {
		// 找不到路由返回
		_ = Send(res, []byte("404"))
	}
}

func Send(w http.ResponseWriter, data []byte) (err error) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是数据流

	_, err = w.Write(data)
	return
}
