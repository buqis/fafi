package card_message

type CardList struct {
	Listings   []Listings `json:"listings"`
	Page       int        `json:"page"`
	PerPage    int        `json:"per_page"`
	TotalPages int        `json:"total_pages"`
}

type Listings struct {
	BestBuyPrice  int    `json:"best_buy_price"`
	BestSellPrice int    `json:"best_sell_price"`
	Item          Item   `json:"item"`
	ListingName   string `json:"listing_name"`
}

type Item struct {
	AugmentEndDate    string `json:"augment_end_date" `
	AugmentText       string `json:"augment_text"`
	BakedImg          string `json:"baked_img"`
	DisplayPosition   string `json:"display_position"`
	Event             bool   `json:"event"`
	HasAugment        bool   `json:"has_augment"`
	HasMatchup        bool   `json:"has_matchup"`
	HasRankChange     bool   `json:"has_rank_change"`
	Img               string `json:"img"`
	IsLiveSet         bool   `json:"is_live_set"`
	Name              string `json:"name"`
	NewRank           int    `json:"new_rank"`
	Ovr               int    `json:"ovr"`
	Rarity            string `json:"rarity"`
	ScBakedImg        string `json:"sc_baked_img"`
	Series            string `json:"series"`
	SeriesTextureName string `json:"series_texture_name"`
	SeriesYear        int    `json:"series_year"`
	SetName           string `json:"set_name"`
	Stars             int    `json:"stars"`
	Team              string `json:"team"`
	TeamShortName     string `json:"team_short_name"`
	Trend             string `json:"trend"`
	Type              string `json:"type"`
	UiAnimIndex       int    `json:"ui_anim_index"`
}
