package local_data

type History struct {
	ID            int    `json:"id" gorm:"primary_key"`
	BestBuyPrice  int    `json:"best_buy_price"`
	BestSellPrice int    `json:"best_sell_price"`
	HistoryDate   string `json:"history_date"`
	ItemID        string `json:"item_ID"`
}

func (History) TableName() string {
	return "historys"
}
