package local_data

type Account struct {
	ID       uint   `json:"id"`        //用户id
	Cookie   string `json:"cookie"`    //用户cookie
	ProxyID  uint   `json:"proxy_id"`  //代理id
	UserName string `json:"user_name"` //用户名
	Password string `json:"password"`  //密码
	UserID   uint   `json:"user_id"`
	Platform uint8  `json:"platform"`
}

// 账号运行状态
const (
	//离线
	CliStateOffline = 1
	//状态等待邮件码
	//CliStateWaitEmailCode = 2
	//登录状态
	CliStateLogin = 3
	//在线状态
	CliStateOnline = 4
	//国家等待贸易
	//CliStateWaitTrade = 5
	//国营贸易
	//CliStateTrading = 6
)

// 账号类型
const (
	GoldAccount     = 1
	CustomerAccount = 2
)

type AccountReader interface {
	GetID() uint
	GetCookie() string
	GetProxyID() uint
	GetUserName() string
	GetPassword() string
	GetUserID() uint
	GetAccountType() byte
	GetPlatform() uint8
}

func (f *Account) GetAccountType() byte {
	return GoldAccount
}
func (a *Account) GetID() uint {
	return a.ID
}
func (a *Account) GetCookie() string {
	return a.Cookie
}
func (a *Account) GetUserName() string {
	return a.UserName
}
func (a *Account) GetPassword() string {
	return a.Password
}
func (a *Account) GetProxyID() uint {
	return a.ProxyID
}
func (a *Account) GetUserID() uint {
	return a.UserID
}
func (a *Account) GetPlatform() uint8 {
	return a.Platform
}

type AccountWriter interface {
	SetID(id uint)
	SetCookie(cookie string)
	SetProxyID(proxyId uint)
	SetUserName(userName string)
	SetPassword(password string)
	SetUserID(userID uint)
	SetPlatform(platform uint8)
}

func (a *Account) SetUserName(userName string) {
	a.UserName = userName
}
func (a *Account) SetPassword(password string) {
	a.Password = password
}
func (a *Account) SetUserID(userID uint) {
	a.UserID = userID
}
func (a *Account) SetID(id uint) {
	a.ID = id
}
func (a *Account) SetCookie(cookie string) {
	a.Cookie = cookie
}
func (a *Account) SetProxyID(proxyId uint) {
	a.ProxyID = proxyId
}
func (a *Account) SetPlatform(platform uint8) {
	a.Platform = platform
}

type AccountWriterAndReader interface {
	AccountWriter
	AccountReader
}
