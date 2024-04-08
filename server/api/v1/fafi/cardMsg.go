package fafi

import (
	"bytes"
	"encoding/json"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/fifa/web_capture/card_message"
	"github.com/flipped-aurora/gin-vue-admin/server/model/fifa/web_capture/software_user"

	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strconv"
)

// 根据页码获取卡片列表
func GetCardList(page int) (*card_message.CardList, error) {
	url := ""
	if page != 0 {
		url = "https://mlb24.theshow.com/apis/listings.json?page=" + strconv.Itoa(page) + "&sort=rank&order=desc&type=mlb_card&min_rank=40&max_rank=99"
	} else {
		url = "https://mlb24.theshow.com/apis/listings.json?sort=rank&order=desc&type=mlb_card&min_rank=40&max_rank=99"
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		global.GVA_LOG.Error("发送请求获取卡片列表失败", zap.Errors("错误信息是", []error{err}))
		panic(err)
		return nil, err
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		global.GVA_LOG.Error("获取卡片列表响应失败", zap.Errors("错误信息是", []error{err}))
		panic(err)
		return nil, err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		global.GVA_LOG.Error("获取卡片列表响应解析失败", zap.Errors("错误信息是", []error{err}))
		panic(err)
		return nil, err
	}
	cardlist := card_message.CardList{}
	err = json.Unmarshal(body, &cardlist)
	if err != nil {
		global.GVA_LOG.Error("卡片列表解析格式失败", zap.Errors("错误信息是", []error{err}))
		panic(err)
		return nil, err
	}
	return &cardlist, nil
}

// 获取卡片详细信息
func GetOneCardMsg(uuid string) (*card_message.CardMsg, error) {
	url := "https://mlb24.theshow.com/apis/listing.json?uuid=" + uuid
	req, err := http.NewRequest(http.MethodGet, url, nil)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		global.GVA_LOG.Error("获取卡片详细信息响应失败", zap.Errors("错误信息是", []error{err}))
		panic(err)
		return nil, err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		global.GVA_LOG.Error("获取卡片详细信息响应解析失败", zap.Errors("错误信息是", []error{err}))
		panic(err)
		return nil, err
	}
	CardMsg := card_message.CardMsg{}
	err = json.Unmarshal(body, &CardMsg)
	if err != nil {
		global.GVA_LOG.Error("卡片详细信息解析格式失败", zap.Errors("错误信息是", []error{err}))
		panic(err)
		return nil, err
	}
	return &CardMsg, nil
}

// 获取正在售卖的卡片的价格
func GetBuyAndSellMsg(user software_user.User, uuid string) (*card_message.CurrentOrders, error) {
	url := "https://mlb24.theshow.com/apis/app/view_listing.json"
	token := software_user.Token{
		AccountId:    user.AccountId,
		AccountToken: user.AccountToken,
		Uuid:         uuid,
	}
	redate, err := json.Marshal(token)
	if err != nil {
		global.GVA_LOG.Error("转换json格式错误", zap.Errors("错误信息是", []error{err}))
		panic(err)
	}
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(redate))
	req.Header.Add("user-agent", "Dart/3.3 (dart:io)")
	req.Header.Add("content-type", "application/json; charset=utf-8")
	req.Header.Add("accept-encoding", "gzip")
	req.Header.Add("host", "mlb24.theshow.com")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		global.GVA_LOG.Error("获取卡片售卖信息响应错误", zap.Errors("错误信息是", []error{err}))
		panic(err)
		return &card_message.CurrentOrders{}, nil
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		global.GVA_LOG.Error("获取卡片售卖信息响应解析失败", zap.Errors("错误信息是", []error{err}))
		panic(err)
		return &card_message.CurrentOrders{}, nil
	}
	CurrentOrders := card_message.CurrentOrders{}
	err = json.Unmarshal(body, &CurrentOrders)
	if err != nil {
		global.GVA_LOG.Error("获取卡片售卖信息解析格式失败", zap.Errors("错误信息是", []error{err}))
		panic(err)
		return &card_message.CurrentOrders{}, nil
	}
	return &CurrentOrders, nil
}
