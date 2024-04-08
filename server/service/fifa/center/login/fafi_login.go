package login

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/model/fifa/web_capture/software_user"
	"github.com/flipped-aurora/gin-vue-admin/server/service/fifa/center/client"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var AbortRedirect = errors.New("abort redirect")

func GetUserMessage(c *client.BaseClient) (string, error) {
	hc := http.Client{
		Transport: &http.Transport{
			Proxy: func(*http.Request) (*url.URL, error) {
				p := c.Proxy
				return p.ToURL()
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
		Jar:     c.Jar,
		Timeout: 0,
	}
	//模拟登录的地址
	u := "https://ca.account.sony.com/api/authz/v3/oauth/authorize?service_entity=urn:service-entity:psn&response_type=code" +
		"&client_id=cd5a424d-cfdc-49cb-b642-66e5f4c9bfd2&redirect_uri=com.mlb.psn.app://redirect&scope=psn:s2s%20openid%20id_token:psn.basic_claims"
	req, _ := http.NewRequest(http.MethodGet, u, nil)
	//设置头消息
	req.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 15_4_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.4 Mobile/15E148 Safari/604.1")
	req.Header.Add("Host", "ca.account.sony.com")
	//获取响应
	response, err := hc.Do(req)
	if err != nil && !errors.Is(err, AbortRedirect) {

	}
	body := response.Header.Get("Location")
	return body, nil
}

// 获取刷新token
func GetToken(code string) (*software_user.User, error) {
	urls := "https://account.theshow.com/psn_sessions/app24_oauth.json?code=" + code + "&redirect_uri=com.mlb.psn.app://redirect"
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
