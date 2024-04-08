package card_message

type CardMsg struct {
	BestBuyPrice   int       `json:"best_buy_price"`
	BestSellPrice  int       `json:"best_sell_price"`
	CompletedOrder []Orders  `json:"completed_orders"`
	Item           Item      `json:"item"`
	ListingName    string    `json:"listing_name"`
	PriceHistory   []History `json:"price_history"`
}
