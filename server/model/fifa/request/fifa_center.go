package request

type AccountEvent struct {
	IDs   []uint `json:"i_ds"`
	Event string `json:"event" example:"start stop delete"`
}

type SetProxyReq struct {
	IDs     []uint `json:"i_ds"`
	ProxyID uint   `json:"proxy_id"`
}
