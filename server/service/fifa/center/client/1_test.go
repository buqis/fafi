package client

import (
	"context"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/fifa/local_data"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestFloor(t *testing.T) {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:3306)/fifa?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	global.GVA_DB = db
	account := local_data.Account{
		ID:       1,
		Cookie:   "npsso=nDJY1ZFdp7hypYEHewc3GMDnQDHnJ69uZmqbU5esc73CNvqx8iBKBDjXiEUnBnI0",
		ProxyID:  1,
		UserName: "chentao",
		Password: "123",
		UserID:   1,
	}
	accounts := account
	accounts.SetProxyID(1)
	accounts.SetID(1)
	accounts.SetUserName("chentao")
	accounts.SetCookie("npsso=nDJY1ZFdp7hypYEHewc3GMDnQDHnJ69uZmqbU5esc73CNvqx8iBKBDjXiEUnBnI0 ")
	accounts.SetPassword("123")
	accounts.SetUserID(1)
	ctx := context.Background()
	BaseClient := newBaseClient(ctx, &accounts)
	user, _ := BaseClient.Start()
	fmt.Println(user)
}
