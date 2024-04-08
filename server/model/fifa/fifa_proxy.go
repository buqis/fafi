package fifa

import (
	"fmt"
	"net/url"
)

type FIFAProxy struct {
	Id       int    `json:"id"`
	UserID   int    `json:"user_id"`
	IP       string `json:"ip"`
	Port     string `json:"port"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Type     uint8  `json:"type"`
}

// 拼接出代理地址
func (f *FIFAProxy) ToURL() (*url.URL, error) {
	auth := ""
	if len(f.UserName) > 0 || len(f.Password) > 0 {
		auth = fmt.Sprintf("%s:%s", f.UserName, f.Password)
	}
	return url.Parse(fmt.Sprintf("socks5://%s%s:%s", auth, f.IP, f.Port))
}
