package web

import (
	"deploy-server/app/tools"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

//responseObj 接口调用之后返回的结构
type responseObj struct {
	Code int `json:"code"`
	// Success bool        `json:"success"`
	TxID    string      `json:"txid,omitempty"`
	Payload interface{} `json:"datas,omitempty"`
	Error   string      `json:"message,omitempty"`
}

//用于列表返回专用结构
type responseListObj struct {
	Code      int         `json:"code"`
	TxID      string      `json:"txid,omitempty"`
	TotalPage int         `json:"totalpage"`
	Payload   interface{} `json:"datas,omitempty"`
	Error     string      `json:"message,omitempty"`
}

//outputFailed  组装接口返回失败值
func outputFailed(w http.ResponseWriter, msg string) {
	one := responseObj{
		Code:  400,
		Error: msg,
	}
	outputJSONObj(w, &one)
}

//outputJSONObj 响应组装header
func outputJSONObj(w http.ResponseWriter, obj interface{}) {
	buff, err := json.Marshal(obj)
	if err != nil {
		logrus.Debugln("OutputInvokeObj", "json failed", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(buff)
}

//OutputEnter 组装API接口返回的值
func OutputEnter(w http.ResponseWriter, id string, payload interface{}, msg error) {
	if msg != nil {
		outputFailed(w, msg.Error())
		return
	}
	outputSuccess(w, id, payload)
}

//outputSuccess 组装API接口成功返回的值
func outputSuccess(w http.ResponseWriter, id string, payload interface{}) {
	ret := responseObj{
		Code:  200,
		Error: "success",
	}
	if len(id) > 0 {
		ret.TxID = id
	}
	switch payload.(type) {
	case []byte:
		tmp := payload.([]byte)
		if len(tmp) > 0 {
			obj, err := tools.TryParseStringToObj(tmp)
			if err != nil {
				ret.Payload = string(tmp)
			} else {
				ret.Payload = obj
			}
		}
	default:
		if payload != nil {
			ret.Payload = payload
		}
	}
	outputJSONObj(w, &ret)
}

//GetParamsBody 获取传进来的参数
func GetParamsBody(r *http.Request) (map[string][]string, []byte) {
	r.ParseForm()
	ret := r.PostForm
	if len(r.Form) > 0 {
		for k, value := range r.Form {
			if len(value) == 0 {
				continue
			}
			for _, v := range value {
				ret.Add(k, v)
			}
		}
	}
	if r.Body == nil {
		return ret, nil
	}
	defer r.Body.Close()

	buff, err := ioutil.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		logrus.WithError(err).Debugln("read body all")
		return ret, nil
	}
	// buff, err := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	logrus.WithError(err).Debugln("read body all")
	// }
	return ret, buff
}

//OutputListEnter 组装API接口返回的值
func OutputListEnter(w http.ResponseWriter, id string, totalPage int, payload interface{}, msg error) {
	if msg != nil {
		outputFailed(w, msg.Error())
		return
	}
	outputListSuccess(w, id, totalPage, payload)
}

//outputSuccess 组装API接口成功返回的值
func outputListSuccess(w http.ResponseWriter, id string, totalPage int, payload interface{}) {
	ret := responseListObj{
		Code:      200,
		Error:     "success",
		TotalPage: totalPage,
	}
	if len(id) > 0 {
		ret.TxID = id
	}
	switch payload.(type) {
	case []byte:
		tmp := payload.([]byte)
		if len(tmp) > 0 {
			obj, err := tools.TryParseStringToObj(tmp)
			if err != nil {
				ret.Payload = string(tmp)
			} else {
				ret.Payload = obj
			}
		}
	default:
		if payload != nil {
			ret.Payload = payload
		}
	}
	outputJSONObj(w, &ret)
}
