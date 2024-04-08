package card_message

type History struct {
	BestBuyPrice  int    `json:"best_buy_price"`
	BestSellPrice int    `json:"best_sell_price"`
	Date          string `json:"date"`
}
