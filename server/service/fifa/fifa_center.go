package fifa

import (
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/fifa/local_data"
	request2 "github.com/flipped-aurora/gin-vue-admin/server/model/fifa/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service/fifa/center"
	proxy2 "github.com/flipped-aurora/gin-vue-admin/server/service/fifa/proxy"
	"gorm.io/gorm"
)

type CenterService struct {
	center center.Center
}

// 获取账户信息列表
func (c *CenterService) GetAccountInfoList(uid uint, info request.PageInfo) (list interface{}, total int64, err error) {
	return c.center.GetAccountInfoList(uid, info)
}

// 账户活动
func (c *CenterService) AccountEvent(uid uint, ae request2.AccountEvent) error {
	switch ae.Event {
	case "start":
		//账户启动
		return c.center.StartAccount(uid, ae.IDs...)
	case "stop":
		//账户停止
		return c.center.StopAccount(uid, ae.IDs...)
	case "delete":
		//账户删除
		return c.center.DeleteAccount(uid, ae.IDs...)
	}
	return nil
}

// 根据账号id获取账号详细信息
func (c *CenterService) GetAccountById(uid, id uint) (account local_data.AccountReader, err error) {
	return c.center.GetAccountById(uid, id), nil
}

// 根据账号id修改账号
func (c *CenterService) UpdateAccount(uid uint, account local_data.Account) error {
	return c.center.UpdateAccount(uid, account)
}

// 设置代理
func (c *CenterService) SetProxy(uid uint, input request2.SetProxyReq) error {
	p, err := proxy2.GetProxy(uid, input.ProxyID, local_data.GoldAccount)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("代理%d不存在", input.ProxyID)
	}
	return c.center.SetAccountProxy(uid, *p, input.IDs)
}
