package fafi

import (
	"net/url"
)

// 带条件搜索
func ConditionalSearch(uv url.Values) string {
	url := "https://mlb24.theshow.com/apis/listings.json?sort=rank"
	//排序
	if uv.Get("order") == "" {
		uv.Set("order", "desc")
	}
	//卡片类型
	if uv.Get("type") == "" {
		uv.Set("type", "mlb_card")
	}
	//排名
	if uv.Get("min_rank") == "" {
		uv.Set("min_rank", "40")
	}
	if uv.Get("max_rank") == "" {
		uv.Set("max_rank", "99")
	}
	return url + uv.Encode()
}

//// 带条件搜索
//func ConditionalSearch(types, rarity, display_position, set_name, team, order string, min_best_buy_price, max_best_buy_price,
//	min_best_sell_price, max_best_sell_price, min_rank, series_id, max_rank, stars int, event, has_augment bool) {
//	url := "https://mlb24.theshow.com/apis/listings.json?sort=rank"
//	//排序
//	if order != "desc" {
//		url += "&order=" + order
//	} else {
//		url += "&order=desc"
//	}
//	//卡片类型
//	if types != "" {
//		url += "&type=" + types
//	} else {
//		url += "&type=mlb_card"
//	}
//	//品质
//	if rarity != "" {
//		url += "&rarity=" + rarity
//	}
//	//带价格筛选
//	if min_best_buy_price != 0 || max_best_buy_price != 0 || min_best_sell_price != 0 || max_best_sell_price != 0 {
//		if min_best_buy_price != 0 {
//			url += "&min_best_sell_price=" + strconv.Itoa(min_best_buy_price)
//		} else {
//			url += "&min_best_sell_price="
//		}
//		if max_best_buy_price != 0 {
//			url += "&max_best_sell_price=" + strconv.Itoa(max_best_buy_price)
//		} else {
//			url += "&max_best_sell_price="
//		}
//		if min_best_sell_price != 0 {
//			url += "&min_best_sell_price=" + strconv.Itoa(min_best_sell_price)
//		} else {
//			url += "&min_best_sell_price="
//		}
//		if max_best_sell_price != 0 {
//			url += "&max_best_sell_price=" + strconv.Itoa(max_best_sell_price)
//		} else {
//			url += "&max_best_sell_price="
//		}
//	}
//	//排名
//	if min_rank != 0 {
//		url += "&min_rank=" + strconv.Itoa(min_rank)
//	} else {
//		url += "&min_rank=40"
//	}
//	if max_rank != 0 {
//		url += "&max_rank=" + strconv.Itoa(max_rank)
//	} else {
//		url += "&max_rank=99"
//	}
//	//位置
//	if display_position != "" {
//		url += "&display_position=" + display_position
//	}
//	//队伍
//	if team != "" {
//		url += "&team=" + team
//	}
//	//集名称
//	if set_name != "" {
//		url += "&set_name=" + set_name
//	}
//	//系列
//	if series_id != 0 {
//		url += "&series_id=" + strconv.Itoa(series_id)
//	}
//	//事件
//	if event {
//		url += "&event=true"
//	}
//	//已经增加
//	if has_augment {
//		url += "&has_augment=true"
//	}
//	//星
//	if stars != 0 {
//		url += "&stars=" + strconv.Itoa(stars)
//	}
//}
