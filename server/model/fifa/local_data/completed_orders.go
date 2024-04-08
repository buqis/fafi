package local_data

// 完成的订单

type Orders struct {
	ID          int    `json:"id" gorm:"primary_key"`
	OrdersDate  string `json:"orders_date"`  //完成时间
	OrdersPrice string `json:"orders_price"` //完成价格
	ItemID      string `json:"item_id"`
}

func (Orders) TableName() string {
	return "completedorders"
}
