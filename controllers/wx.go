package controllers

import (
	"github.com/astaxie/beego"
	"strings"
	"sort"
	"crypto/sha1"
	"io"
	"fmt"
	"net/http"
	"encoding/xml"
	"time"
	"io/ioutil"
	"log"
)

type WXController struct {
	beego.Controller
}

func (c *WXController) Get() {
	checkSignature(c.Ctx.ResponseWriter,c.Ctx.Request)
}

func (c *WXController) Post() {
	fmt.Println("2222")
	procRequest(c.Ctx.ResponseWriter,c.Ctx.Request)
}

func str2sha1(data string)string{
	t:=sha1.New()
	io.WriteString(t,data)
	return fmt.Sprintf("%x",t.Sum(nil))
}

func checkSignature(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	var token string="hello"
	var signature string=strings.Join(r.Form["signature"],"")
	var timestamp string=strings.Join(r.Form["timestamp"],"")
	var nonce string=strings.Join(r.Form["nonce"],"")
	var echostr string=strings.Join(r.Form["echostr"],"")
	tmps:=[]string{token,timestamp,nonce}
	sort.Strings(tmps)
	tmpStr:=tmps[0]+tmps[1]+tmps[2]
	tmp:=str2sha1(tmpStr)
	fmt.Println("rs:",tmp==signature)
	if tmp==signature{
		fmt.Fprintf(w,echostr)
	}
}






type TextRequestBody struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgType      string
	Content      string
	MsgId        int
}
func parseTextRequestBody(r *http.Request) *TextRequestBody {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	fmt.Println(string(body))
	requestBody := &TextRequestBody{}
	xml.Unmarshal(body, requestBody)
	return requestBody
}
type TextResponseBody struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   CDATAText
	FromUserName CDATAText
	CreateTime   time.Duration
	MsgType      CDATAText
	Content      CDATAText
}

type CDATAText struct {
	Text string `xml:",innerxml"`
}

func value2CDATA(v string) CDATAText {
	return CDATAText{""}
}

func makeTextResponseBody(fromUserName, toUserName, content string) ([]byte, error) {
	textResponseBody := &TextResponseBody{}
	textResponseBody.FromUserName = value2CDATA(fromUserName)
	textResponseBody.ToUserName = value2CDATA(toUserName)
	textResponseBody.MsgType = value2CDATA("text")
	textResponseBody.Content = value2CDATA(content)
	textResponseBody.CreateTime = time.Duration(time.Now().Unix())
	return xml.MarshalIndent(textResponseBody, " ", "  ")
}

func procRequest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
/*	if !validateUrl(w, r) {
		log.Println("Wechat Service: this http request is not from Wechat platform!")
		return
	}
*/
	fmt.Println("11111")
	if r.Method == "POST" {
		textRequestBody := parseTextRequestBody(r)
		if textRequestBody != nil {
			fmt.Printf("Wechat Service: Recv text msg [%s] from user [%s]!",
				textRequestBody.Content,
				textRequestBody.FromUserName)
			responseTextBody, err := makeTextResponseBody(textRequestBody.ToUserName,
				textRequestBody.FromUserName,
				"Hello, "+textRequestBody.FromUserName)
			if err != nil {
				log.Println("Wechat Service: makeTextResponseBody error: ", err)
				return
			}
			fmt.Fprintf(w, string(responseTextBody))
		}
	}
}