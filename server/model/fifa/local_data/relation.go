package local_data

type Item struct {
	ID                int              `json:"id" gorm:"primary_key"`
	AugmentEndDate    string           `json:"augment_end_date" `
	AugmentText       string           `json:"augment_text"`
	BakedImg          string           `json:"baked_img"`
	DisplayPosition   string           `json:"display_position"`
	Event             bool             `json:"event"`
	HasAugment        bool             `json:"has_augment"`
	HasMatchup        bool             `json:"has_matchup"`
	HasRankChange     bool             `json:"has_rank_change"`
	Img               string           `json:"img"`
	IsLiveSet         bool             `json:"is_live_set"`
	Name              string           `json:"name"`
	NewRank           int              `json:"new_rank"`
	Ovr               int              `json:"ovr"`
	Rarity            string           `json:"rarity"`
	ScBakedImg        string           `json:"sc_baked_img"`
	Series            string           `json:"series"`
	SeriesTextureName string           `json:"series_texture_name"`
	SeriesYear        int              `json:"series_year"`
	SetName           string           `json:"set_name"`
	Stars             string           `json:"stars"`
	Team              string           `json:"team"`
	TeamShortName     string           `json:"team_short_name"`
	Trend             string           `json:"trend"`
	Type              string           `json:"type"`
	UiAnimIndex       int              `json:"ui_anim_index"`
	Uuid              string           `json:"uuid"`
	History           `json:"history"` //卡片购买历史等于trends
	Data              `json:"data"`    //卡片购买价格等于current orders
	Orders            `json:"orders"`  //销售记录等于completed
}
