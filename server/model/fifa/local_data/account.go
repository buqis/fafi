package local_data

type Account struct {
	ID       int    `json:"id"`        //用户id
	Cookie   string `json:"cookie"`    //用户cookie
	ProxyID  int    `json:"proxy_id"`  //代理id
	UserName string `json:"user_name"` //用户名
	Password string `json:"password"`  //密码
	UserID   int    `json:"user_id"`
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
	GetID() int
	GetCookie() string
	GetProxyID() int
	GetUserName() string
	GetPassword() string
	GetUserID() int
	GetAccountType() byte
	GetPlatform() uint8
}

func (f *Account) GetAccountType() byte {
	return GoldAccount
}
func (a *Account) GetID() int {
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
func (a *Account) GetProxyID() int {
	return a.ProxyID
}
func (a *Account) GetUserID() int {
	return a.UserID
}
func (a *Account) GetPlatform() uint8 {
	return a.Platform
}

type AccountWriter interface {
	SetID(id int)
	SetCookie(cookie string)
	SetProxyID(proxyId int)
	SetUserName(userName string)
	SetPassword(password string)
	SetUserID(userID int)
	SetPlatform(platform uint8)
}

func (a *Account) SetUserName(userName string) {
	a.UserName = userName
}
func (a *Account) SetPassword(password string) {
	a.Password = password
}
func (a *Account) SetUserID(userID int) {
	a.UserID = userID
}
func (a *Account) SetID(id int) {
	a.ID = id
}
func (a *Account) SetCookie(cookie string) {
	a.Cookie = cookie
}
func (a *Account) SetProxyID(proxyId int) {
	a.ProxyID = proxyId
}
func (a *Account) SetPlatform(platform uint8) {
	a.Platform = platform
}

type AccountWriterAndReader interface {
	AccountWriter
	AccountReader
}
