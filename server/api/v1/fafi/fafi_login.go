package fafi

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/model/fifa/web_capture/software_user"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

//// 获取刷新token
//func GetToken(urls string) (*software_user.User, error) {
//	//对传进来的url地址进行解析
//	ul, err := url.Parse(urls)
//	if err != nil {
//		global.GVA_LOG.Error(err.Error())
//	}
//	//获取url地址中的信息
//	code := ul.Query().Get("code")
//	cid := ul.Query().Get("cid")
//	//模拟登录的地址
//	u := "https://account.theshow.com/psn_sessions/app24_oauth.json?code=" + code + "&redirect_uri=" + cid
//	req, _ := http.NewRequest(http.MethodGet, u, nil)
//	//设置头消息
//	req.Header.Add("user-agent", "Dart/3.3 (dart:io)")
//	req.Header.Add("content-type", "application/json; charset=utf-8")
//	req.Header.Add("accept-encoding", "gzip")
//	req.Header.Add("host", "account.theshow.com")
//	//获取响应
//	response, err := http.DefaultClient.Do(req)
//	body, _ := ioutil.ReadAll(response.Body)
//	if err != nil {
//		fmt.Println(err)
//	}
//	//把响应中的数据转化为用户结构体
//	user := software_user.User{}
//	err = json.Unmarshal(body, &user)
//	if err != nil {
//		panic(err)
//	}
//	return &user, nil
//}

var AbortRedirect = errors.New("abort redirect")

// 获取刷新token的地址
func FafiLogin(client_id string) (string, error) {
	jar, _ := cookiejar.New(nil)
	hc := http.Client{
		Transport: &http.Transport{
			Proxy: func(*http.Request) (*url.URL, error) {
				return url.Parse("socks5://127.0.0.1:9090")
			},
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if req.Response.StatusCode == http.StatusFound && strings.HasPrefix(req.Response.Header.Get("Location"), "com.mlb.psn.app") {
				return AbortRedirect
			}
			return nil
		},
		Jar:     jar,
		Timeout: 0,
	}
	//模拟登录的地址
	u := "https://ca.account.sony.com/api/authz/v3/oauth/authorize?service_entity=urn:service-entity:psn&response_type=code" +
		"&client_id=" + client_id + "&redirect_uri=com.mlb.psn.app://redirect&scope=psn:s2s%20openid%20id_token:psn.basic_claims"
	req, _ := http.NewRequest(http.MethodGet, u, nil)
	//设置头消息
	req.Header.Add("Accept-Language", "zh-CN,zh-Hans;q=0.9")
	req.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 15_4_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.4 Mobile/15E148 Safari/604.1")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Add("Host", "ca.account.sony.com")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Connection", "keep-alive")
	//获取响应
	response, err := hc.Do(req)
	if err != nil && !errors.Is(err, AbortRedirect) {

	}
	body := response.Header.Get("Location")
	return body, nil
}

// 获取刷新token
func GetToken(urls string) (*software_user.User, error) {
	//对传进来的url地址进行解析
	req, _ := http.NewRequest(http.MethodGet, urls, nil)
	//设置头消息
	req.Header.Add("user-agent", "Dart/3.3 (dart:io)")
	req.Header.Add("content-type", "application/json; charset=utf-8")
	req.Header.Add("accept-encoding", "gzip")
	req.Header.Add("host", "account.theshow.com")
	//获取响应
	response, err := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}
	//把响应中的数据转化为用户结构体
	user := software_user.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		panic(err)
	}
	return &user, nil
}
