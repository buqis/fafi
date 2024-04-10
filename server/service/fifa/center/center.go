package center

import (
	"context"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/fifa"
	"github.com/flipped-aurora/gin-vue-admin/server/model/fifa/local_data"
	"github.com/flipped-aurora/gin-vue-admin/server/service/fifa/center/client"
	"go.uber.org/zap"
	"sync"
	"time"
)

type Center struct {
	ctx     context.Context
	rw      sync.RWMutex
	list    []client.ClientI
	once    sync.Once
	channel map[string]chan interface{}
}

// 修改账户
func (c *Center) UpdateAccount(uid uint, account local_data.Account) error {
	c.lazy()
	c.rw.Lock()
	defer c.rw.Unlock()
	index := c.getIndexByID(uid, account.ID)
	if index == -1 {
		return fmt.Errorf("账号%d不存在", uid)
	}
	return c.list[index].SetAccount(&account)
}

// 设置账户代理
func (c *Center) SetAccountProxy(uid uint, p fifa.FIFAProxy, ids []uint) error {
	c.lazy()
	c.rw.RLock()
	defer c.rw.RUnlock()
	for _, id := range ids {
		obj := c.getClientByID(uid, id)
		if obj == nil {
			continue
		}
		if err := obj.SetProxy(p); err != nil {
			return err
		}
	}
	return nil
}

// 根据id获取客户端
func (c *Center) getClientByID(uid uint, id uint) client.ClientI {
	for _, v := range c.list {
		if v.GetID() == id {
			v := v
			return v
		}
	}
	return nil
}

// 根据用户id获取账户
func (c *Center) GetAccountById(uid, id uint) local_data.AccountReader {
	c.lazy()
	c.rw.RLock()
	defer c.rw.RUnlock()
	index := c.getIndexByID(uid, id)
	if index == -1 {
		return nil
	}
	return c.list[index].GetAccount()
}

// 添加账户
func (c *Center) AddAccount(uid uint, accounts ...local_data.Account) error {
	c.lazy()
	c.rw.Lock()
	defer c.rw.Unlock()
	return c.addAccount(uid, accounts...)
}

// 添加账户，把客户端添加到集合中，并且把账户添加到数据库
func (c *Center) addAccount(uid uint, accounts ...local_data.Account) error {
	for _, account := range accounts {
		account := account
		if uid > 0 {
			account.UserID = uid
		}
		//1.加入数据库
		if err := global.GVA_DB.Save(&account).Error; err != nil {
			return err
		}
		channel := c.getChannel(account.UserID, account.Platform)
		ctx := context.WithValue(c.ctx, "channel", channel)
		tClient := client.NewBaseClient(ctx, &account)
		c.list = append(c.list, tClient)
	}
	return nil
}

// 查找所有账户
func (c *Center) loadAccountFromDB() error {
	var accounts []local_data.Account
	if err := global.GVA_DB.Model(&local_data.Account{}).Find(&accounts).Error; err != nil {
		return err
	}
	return c.AddAccount(0, accounts...)
}

func (c *Center) getChannel(uid uint, platform uint8) chan interface{} {
	key := fmt.Sprintf("%d-%d", uid, platform)
	_, ok := c.channel[key]
	if !ok {
		c.channel[key] = make(chan interface{})
	}
	return c.channel[key]
}

// 获取账户信息列表
func (c *Center) GetAccountInfoList(uid uint, info request.PageInfo) (list interface{}, total int64, err error) {
	c.lazy()
	c.rw.RLock()
	defer c.rw.RUnlock()
	//根据用户id获取账户列表
	clientList := c.getAccountListByUserID(uid)
	//账户列表的长度
	length := len(clientList)
	total = int64(length)
	tempList := make([]map[string]interface{}, 0, info.PageSize)
	//
	for index := (info.Page - 1) * info.PageSize; index < info.Page*info.PageSize && index < length; index++ {
		item := clientList[index]
		tempList = append(tempList, item.ToMap())
	}
	list = tempList
	return
}

// 根据用户id获得账户列表
func (c *Center) getAccountListByUserID(uid uint) []client.ClientI {
	//创建一个客户端集合，长度是客户端的长度
	list := make([]client.ClientI, 0, len(c.list))
	//遍历客户端集合
	for _, v := range c.list {
		//获取每一个账户实例
		accInfo := v.GetAccount()
		//如果账户的用户id等于传进来的用户id，则添加到集合中
		if accInfo.GetUserID() == uid {
			v := v
			list = append(list, v)
		}
	}
	//返回账户集合
	return list
}

// 根据用户id获取下标
func (c *Center) getIndexByID(uid, id uint) int {
	for i, v := range c.list {
		if v.GetID() == id {
			return i
		}
	}
	return -1
}

// 懒加载
func (c *Center) lazy() {
	c.once.Do(func() {
		c.ctx = context.Background()
		c.channel = make(map[string]chan interface{})
		if err := c.loadAccountFromDB(); err != nil {
			global.GVA_LOG.Error("init center failed", zap.Error(err))
		} else {
			global.GVA_LOG.Info("init center success")
		}
	})
}

// 启动账号
func (c *Center) StartAccount(uid uint, ids ...uint) error {
	c.lazy()
	c.rw.RLock()
	defer c.rw.RUnlock()
	//遍历用户活动的账号
	for i, id := range ids {
		//通过用户id，获取对应的要进行启动的账号的下标，如果找不到返回-1
		index := c.getIndexByID(uid, id)
		if index == -1 {
			return fmt.Errorf("账号ID: %d不存在", id)
		}
		tmp := time.Duration(2*i) * time.Second
		//根据下标来启动对应的账号
		if err := c.list[index].Start(tmp); err != nil {
			return fmt.Errorf("启动账号: %d 错误: %w", id, err)
		}
	}
	return nil
}

// 停止账号
func (c *Center) StopAccount(uid uint, ids ...uint) error {
	c.lazy()
	c.rw.RLock()
	defer c.rw.RUnlock()

	for _, id := range ids {
		index := c.getIndexByID(uid, id)
		if index == -1 {
			return fmt.Errorf("账号ID: %d不存在", id)
		}
		if err := c.list[index].Stop(); err != nil {
			return fmt.Errorf("停止账号: %d 错误: %w", id, err)
		}
	}
	return nil
}

// 删除账户
func (c *Center) DeleteAccount(uid uint, ids ...uint) error {
	c.lazy()
	c.rw.RLock()
	defer c.rw.RUnlock()

	for _, id := range ids {
		index := c.getIndexByID(uid, id)
		if index == -1 {
			return fmt.Errorf("账号ID: %d 不存在", id)
		}
		if err := c.list[index].Delete(); err != nil {

		}
	}
	return nil
}
