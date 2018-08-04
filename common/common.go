package common

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
)

//从请求体中读取数据转化为对象
func ProcJsonRequest(resp http.ResponseWriter, req *http.Request,objPtr interface{}) {
	body, _ := ioutil.ReadAll(req.Body)
	json.Unmarshal(body, objPtr)
}
//func ProcJsonRespose(resp http.ResponseWriter, req *http.Request,objPtr interface{}) {
//	body, _ := ioutil.ReadAll(req.Body)
//	json.Marshal("")
//}