package client

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/fifa"
	"github.com/flipped-aurora/gin-vue-admin/server/model/fifa/local_data"
	"github.com/flipped-aurora/gin-vue-admin/server/model/fifa/web_capture/software_user"
	"github.com/flipped-aurora/gin-vue-admin/server/service/fifa/center/login"
	"github.com/flipped-aurora/gin-vue-admin/server/service/fifa/proxy"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type BaseClient struct {
	Account local_data.AccountWriterAndReader
	Proxy   *fifa.FIFAProxy
	Jar     *utils.Jar
	Task    func() error
	Client  *http.Client
	Once    sync.Once
}

func (c *BaseClient) start() error {
	if err := c.CheckProxyStatus(); err != nil {
		fmt.Println(err)
	}
	if err := c.InitHttpClient(); err != nil {
		fmt.Println(err)
	}
	if user, err := c.Login(); err != nil {
		fmt.Println(user)
		return err
	}
	return nil
}

// 检查代理状态
func (c *BaseClient) CheckProxyStatus() error {
	//根据用户的id，代理id，类型去获取代理
	temp, err := proxy.GetProxyOrRandom(c.Account.GetUserID(), c.Account.GetProxyID(), c.Account.GetAccountType())
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("获取代理失败")
		return err
	}
	//把用户的代理赋值给客户端
	c.Proxy = temp
	//如果客户端代理和用户原本的代理不同，则更新用户的代理
	if c.Proxy != nil && c.Proxy.Id != c.Account.GetProxyID() {
		c.Account.SetProxyID(c.Proxy.Id)
		err := global.GVA_DB.Model(&c.Account).Where("id = ?", c.Account.GetID()).Update("proxy_id = ?", c.Proxy.Id).Error
		if err != nil {
			fmt.Println("更新用户代理失败")
			return err
		}
	}
	//如果客户端代理为空，则无代理可用
	if c.Proxy == nil {
		fmt.Errorf("无代理可用")
	}
	return nil
}

func (c *BaseClient) InitHttpClient() error {
	funcProxy := func(req *http.Request) (*url.URL, error) {
		p := c.Proxy
		if p == nil {
			fmt.Errorf("无代理可用")
		}
		return p.ToURL()
	}
	if c.Jar == nil {
		c.Jar = new(utils.Jar)
		if len(c.Account.GetCookie()) > 0 {
			if err := c.Jar.Import(c.Account.GetCookie()); err != nil {
				return fmt.Errorf("导入Cookies失败 %w", err)
			}
		}
	}
	if c.Client == nil {
		c.Client = &http.Client{
			Transport: &http.Transport{
				Proxy:       funcProxy,
				DialContext: nil,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
			Jar:     c.Jar,
			Timeout: time.Second * 30,
		}
	}
	return nil
}

func (c *BaseClient) Login() (*software_user.User, error) {
	urls, err := login.GetUserMessage(c)
	if err != nil {
		fmt.Println("登录信息有误")
		return nil, err
	}
	u, err := url.Parse(urls)
	code := u.Query()
	co := code["code"]
	user, err := login.GetToken(co[0])
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}
	return nil, err
}
