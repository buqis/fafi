package fifa

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/fifa/local_data"
	"strconv"
	"time"
)

// 查询对应球员卡片的售卖信息
func SelectOneDataByItemID(id int) (*local_data.Data, error) {
	var list local_data.Data
	now := time.Now()
	m, _ := time.ParseDuration("-30m")
	err := global.GVA_DB.Debug().Preload("ItemsToBuys").Preload("ItemsToSells").Preload("Infos").
		Where("datas.date > ?", now.Add(m).Format("2006-01-02 15:04:05")).
		Where(" datas.item_id = ?", id).
		Find(&list)
	if err != nil {
		panic(err)
	}
	return &list, nil
}

// 查找部分卡片,根据星级进行排序
func SelectAllItemMessageOrderByOvr(page, size int) *[]local_data.Item {
	var list []local_data.Item
	sql := "SELECT * FROM `items` ORDER BY ovr desc limit " + strconv.Itoa(page) + "," + strconv.Itoa(size)
	global.GVA_DB.Raw(sql).Find(&list)
	return &list
}

// 查询对应球员卡片的日历史销售记录
func SelectItemHistoryByItemID(id int) *[]local_data.History {
	var list []local_data.History
	global.GVA_DB.Where("item_id = ?", id).Find(&list)
	return &list
}

// 查询对应球员卡片的历史销售记录
func SelectItemSellHistoryByItemID(id int) *[]local_data.Orders {
	var list []local_data.Orders
	global.GVA_DB.Where("item_id = ?", id).Order("orders_date desc").Find(&list)
	return &list
}
