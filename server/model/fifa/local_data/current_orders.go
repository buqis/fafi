package local_data

import (
	"time"
)

type Data struct {
	ID           int           `json:"id" `
	Date         time.Time     `json:"date"`
	Infos        Info          `json:"info"  `          //信息
	ItemsToBuys  []ItemsToBuy  `json:"items_to_buy"  `  //购买列表
	ItemsToSells []ItemsToSell `json:"items_to_sell"  ` //销售列表
	ItemID       string        `json:"item_id" `
}

func (Data) TableName() string {
	return "datas"
}

type Info struct {
	ID       int    `json:"id"  `
	Owned    string `json:"owned"  `     //拥有
	SalesTax string `json:"sales_tax"  ` //销售税
	Sellable string `json:"sellable"  `  //适于销售的
	DataID   int    `json:"data_id" `
}

type ItemsToBuy struct {
	ID              int    `json:"id"  `
	DisplayPrice    string `json:"display_price"  `      //显示的价格
	DisplayQuantity string `json:"display_quantity"  `   //显示数量
	ItemsToBuyPrice string `json:"items_to_buy_price"  ` //价格
	Quantity        string `json:"quantity"  `           //数量
	DataID          int    `json:"data_id"`
}

type ItemsToSell struct {
	ID              int    `json:"id" `
	DisplayPrice    string `json:"display_price" `
	DisplayQuantity string `json:"display_quantity" `
	ItemsToBuyPrice string `json:"items_to_buy_price" `
	Quantity        string `json:"quantity" `
	DataID          int    `json:"data_id" `
}
