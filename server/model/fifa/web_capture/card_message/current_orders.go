package card_message

// 当前的订单
type CurrentOrders struct {
	Data Data `json:"data"`
}

type Data struct {
	Info        Info          `json:"info"`          //信息
	ItemsToBuy  []ItemsToBuy  `json:"items_to_buy"`  //购买列表
	ItemsToSell []ItemsToSell `json:"items_to_sell"` //销售列表
}

type Info struct {
	Owned    string `json:"owend"`    //拥有
	SalesTax string `json:"salesTax"` //销售税
	Sellable string `json:"sellable"` //适于销售的
}

type ItemsToBuy struct {
	DisplayPrice    string `json:"display_price"`    //显示的价格
	DisplayQuantity string `json:"display_quantity"` //显示数量
	Price           string `json:"price"`            //价格
	Quantity        string `json:"quantity"`         //数量
}

type ItemsToSell struct {
	DisplayPrice    string `json:"display_price"`
	DisplayQuantity string `json:"display_quantity"`
	Price           string `json:"price"`
	Quantity        string `json:"quantity"`
}
