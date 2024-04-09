package center

import (
	"context"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/fifa/local_data"
	"github.com/flipped-aurora/gin-vue-admin/server/service/fifa/center/client"
	"sync"
)

type Center struct {
	ctx     context.Context
	rw      sync.RWMutex
	list    []client.BaseClient
	once    sync.Once
	channel map[string]chan interface{}
}

func (c *Center) AddAccount(uid int, accounts ...local_data.Account) error {
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
		c.list = append(c.list, *tClient)
	}
	return nil
}

func (c *Center) getChannel(uid int, platform uint8) chan interface{} {
	key := fmt.Sprintf("%d-%d", uid, platform)
	_, ok := c.channel[key]
	if !ok {
		c.channel[key] = make(chan interface{})
	}
	return c.channel[key]
}

func (c *Center) lazy() {
	c.once.Do(func() {
		c.ctx = context.Background()
		c.channel = make(map[string]chan interface{})
		if err := c.loadAccountFromDB(); err != nil {

		}
	})
}
