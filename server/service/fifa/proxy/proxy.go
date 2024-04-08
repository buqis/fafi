package proxy

import (
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/fifa"
	"gorm.io/gorm"
)

// 获取代理，如果没有则从数据库随机读取一个
func GetProxyOrRandom(id, uid int, typ uint8) (*fifa.FIFAProxy, error) {
	//使用旧的代理
	if id > 0 {
		p := &fifa.FIFAProxy{
			Id:     id,
			Type:   typ,
			UserID: uid,
		}
		//相当于select * from proxy where id = id and type = type and user_id = user_id
		err := global.GVA_DB.Model(p).First(p, p).Error
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("代理获取失败")
			return p, err
		}
	}
	//创建一个结构体数组，主要是为了判断有没有读到数据，结构体没有办法与nil作比较
	var proxy []fifa.FIFAProxy
	//从数据库随机读取第一个查询到的代理
	err := global.GVA_DB.Raw("select * from proxys where user_id = ? and type = ? order by rand() limit 0,1", uid, typ).Find(&proxy).Error
	if err != nil {
		return nil, err
	}
	if len(proxy) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	temp := proxy[0]
	return &temp, nil
}
