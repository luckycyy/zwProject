package wx

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"zw/constant"
)

func GetWXAccessToken() string{
	resp, err := http.Get(constant.WXTokenURL)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	return string(body)
}