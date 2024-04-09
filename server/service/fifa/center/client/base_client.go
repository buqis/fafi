package client

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/fifa"
	"github.com/flipped-aurora/gin-vue-admin/server/model/fifa/local_data"
	"github.com/flipped-aurora/gin-vue-admin/server/model/fifa/web_capture/software_user"
	"github.com/flipped-aurora/gin-vue-admin/server/service/fifa/center/client/gameerror"
	"github.com/flipped-aurora/gin-vue-admin/server/service/fifa/proxy"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type BaseClient struct {
	Account      local_data.AccountWriterAndReader
	Proxy        *fifa.FIFAProxy
	Jar          *utils.Jar
	Task         func() (*software_user.User, error)
	Client       *http.Client
	Once         sync.Once
	StartLock    sync.Mutex         //用于判断客户端是否启动
	Comment      string             //返回最后的错误信息
	Background   sync.WaitGroup     //线程启动数量，让主线程等待
	IsDelete     bool               //是否删除
	Ctx          context.Context    //上下文对象
	Cancel       context.CancelFunc //避免上下文泄露
	ParentCtx    context.Context    //父级上下文对象
	CliState     int                //账号状态
	RestartCount int                //账号重启次数
}

type ClientI interface {
	GetID()
	start(interval time.Duration)
	NewAccount(acc local_data.AccountReader)
	Start()
	CheckProxyStatus()
	InitHttpClient()
	Login()
	SetCliState(cliState int)
}

func (b *BaseClient) GetID() int {
	return b.Account.GetID()
}

// 客户端启动方法
func (b *BaseClient) start(interval time.Duration) error {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				b.Comment = fmt.Sprintf("%v", err)
			}
		}()
		//尝试给客户端加锁，如果加锁失败，说明客户端已经启动，直接退出即可
		l := b.StartLock.TryLock()
		if !l {
			fmt.Println("账号已启动")
			return
		}
		defer b.StartLock.Unlock()

		//给主协程添加数量，告诉主协程还有多少协程没有关闭
		b.Background.Add(1)
		defer b.Background.Done()

		//判断客户端是否删除
		if b.IsDelete {
			return
		}

		//如果客户端的父上下文信息为空，生成一个空白上下文对象
		if b.ParentCtx == nil {
			b.ParentCtx = context.Background()
		}
		//
		b.Ctx, b.Cancel = context.WithCancel(b.ParentCtx)
		defer b.Cancel()

		b.Comment = fmt.Sprintf("账号启动等待%d秒", interval/time.Second)

		select {
		//如果函数执行结束，则会往通道里写入默认值，否则会等待
		case <-b.Ctx.Done():
			return
		case <-time.After(interval):
		}
		defer func() {
			if err := global.GVA_DB.Save(b.Account).Error; err != nil {
				b.Error("账号停止保存数据", zap.Error(err))
			}
		}()
		//设置为离线状态
		defer b.SetCliState(local_data.CliStateOffline)

	forFlag:
		for {
			b.Comment = ""
			//下线
			b.SetCliState(local_data.CliStateOffline)
			//被停止
			select {
			case <-b.Ctx.Done():
				break forFlag
			default:
			}
			if _, err := b.Task(); err != nil {
				//有错误
				b.Info("执行失败", zap.Error(err))
				//前端错误提示,客户端请求超时
				if errors.Is(err, context.Canceled) {
					//主动停止
					b.Comment = "主动停止"
				} else {
					b.Comment = err.Error()
				}
				var (
					restartErr *gameerror.RestartFlag
				)
				if errors.Is(err, context.Canceled) {
					break
				} else if errors.As(err, &restartErr) {
					if b.RestartCount > 3 {
						b.Info("账号重启次数过多")
						break
					}
					b.RestartCount++
					b.Info("账号重启")

				} else {
					break
				}
			} else {
				//正常停止
				break
			}
			b.Info("账号重启计数", zap.Int("restartCount", b.RestartCount))
			b.RestartCount++
			if b.RestartCount >= 5 {
				break
			}
			time.Sleep(10 * time.Second)
		}
	}()
	return nil
}

// 创建一个账户对象
func (b *BaseClient) NewAccount(acc local_data.AccountReader) error {
	had := false
	if len(acc.GetPassword()) > 0 {
		b.Account.SetPassword(acc.GetPassword())
		had = true
	}
	if acc.GetProxyID() > 0 {
		var tempProxy fifa.FIFAProxy
		err := global.GVA_DB.Model(&fifa.FIFAProxy{}).
			Where(&fifa.FIFAProxy{
				Id:     acc.GetProxyID(),
				UserID: acc.GetUserID(),
			}).First(&tempProxy).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("代理%d不存在", acc.GetProxyID())
		} else if err != nil {
			return fmt.Errorf("获取代理错误，原因：%w", err)
		}
		b.Account.SetProxyID(acc.GetProxyID())
		b.Proxy = &tempProxy
		had = true
	}
	if !had {
		return fmt.Errorf("账号无属性可更新")
	}

	return global.GVA_DB.Save(b.Account).Error
}

// 获得账户的token和uuid
func (c *BaseClient) Start() (*software_user.User, error) {
	c.Info("启动账号")
	defer c.Info("账号停止")
	if err := c.CheckProxyStatus(); err != nil {
		return nil, err
	}
	if err := c.InitHttpClient(); err != nil {
		return nil, err
	}
	user, err := c.Login()
	if err != nil {
		return nil, err
	}
	return user, nil
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
		return fmt.Errorf("无代理可用")
	}
	return nil
}

// 初始化客户端
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
		fmt.Println(c.Account.GetCookie())
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

// 登录
func (c *BaseClient) Login() (*software_user.User, error) {
	urls, err := GetUserMessage(c)
	if err != nil {
		fmt.Println("登录信息有误")
		return nil, err
	}
	u, err := url.Parse(urls)
	code := u.Query()
	co := code["code"]
	user, err := GetToken(co[0])
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}
	return nil, err
}

func (b *BaseClient) Debug(msg string, fields ...zap.Field) {
	fmt.Println(msg, fields)
}
func (b *BaseClient) Info(msg string, fields ...zap.Field) {
	fmt.Println(msg, fields)
}
func (b *BaseClient) Warn(msg string, fields ...zap.Field) {
	fmt.Println(msg, fields)
}
func (b *BaseClient) Error(msg string, fields ...zap.Field) {
	fmt.Println(msg, fields)
}

func (b *BaseClient) SetCliState(cliState int) {
	b.CliState = cliState
}
