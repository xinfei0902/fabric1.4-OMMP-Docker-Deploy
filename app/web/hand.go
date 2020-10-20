package web

import (
	"deploy-server/app/derrors"
	"net/http"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

var globalHandles map[string]http.Handler

var globalHandleFuncs map[string]http.HandlerFunc

func initHandle() {
	globalHandles = make(map[string]http.Handler)
	globalHandleFuncs = make(map[string]http.HandlerFunc)
}

//middleLevelHead 进入接口 运行服务日志显示 请求地址 请求方式 时间（时间戳）请求方地址
func middleLevelHead(now time.Time, r *http.Request) {
	r.ParseForm()
	one := logrus.WithField("uri", r.RequestURI)
	one = one.WithField("method", r.Method)
	one = one.WithField("time", now)
	one = one.WithField("remote", r.RemoteAddr)

	one.Debugln("enter")
}

//middleLevelTail 离开接口 运行服务日志显示 请求地址 请求方式 时间（时间戳）请求方地址
func middleLevelTail(now time.Time, r *http.Request) {
	one := logrus.WithField("uri", r.RequestURI)
	one = one.WithField("method", r.Method)
	one = one.WithField("time", time.Now().Unix())
	one = one.WithField("enter", now)
	one = one.WithField("remote", r.RemoteAddr)

	one.Debugln("leave")
}

// PushHandle into http default router
func PushHandle(path string, h http.Handler) error {
	if len(path) == 0 || h == nil {
		return derrors.ErrorEmptyValue
	}

	path = filepath.ToSlash(path)
	_, ok := globalHandles[path]
	if ok {
		return derrors.ErrorSameKeyExist
	}

	globalHandles[path] = h
	return nil
}

// PushHandleFunc into http default router
func PushHandleFunc(path string, hf http.HandlerFunc) error {
	if len(path) == 0 || hf == nil {
		return derrors.ErrorEmptyValue
	}
	//表示进入接口 退出接口 表示执行完毕
	one := func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		middleLevelHead(now, r)
		hf(w, r)
		middleLevelTail(now, r)
	}
	return pushHandleCore(path, one)
}

//pushHandleCore 路径与函数句柄对应
func pushHandleCore(path string, hf http.HandlerFunc) error {
	if len(path) == 0 || hf == nil {
		return derrors.ErrorEmptyValue
	}
	path = filepath.ToSlash(path)
	_, ok := globalHandleFuncs[path]
	if ok {
		return derrors.ErrorSameKeyExist
	}
	globalHandleFuncs[path] = hf
	return nil
}

//baseSignUPHandle 清理多路复用器  重新导入请求地址
func baseSignUPHandle() *http.ServeMux {
	// clear http handles
	ret := http.NewServeMux()

	for k, v := range globalHandleFuncs {
		ret.HandleFunc(k, v)
	}

	for k, v := range globalHandles {
		ret.Handle(k, v)
	}

	return ret
}
