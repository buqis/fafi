package local_data

type Account struct {
	ID       int    `json:"id"`       //用户id
	Cookie   string `json:"cookie"`   //用户cookie
	ProxyID  int    `json:"proxy_id"` //代理id
	UserName string `json:"user_name"`
	Password string `json:"password"`
	UserID   int    `json:"user_id"`
}

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

type AccountWriter interface {
	SetID(id int)
	SetCookie(cookie string)
	SetProxyID(proxyId int)
	SetUserName(UserName string)
	SetPassword(Password string)
	SetUserID(UserID int) int
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

type AccountWriterAndReader struct {
	AccountWriter
	AccountReader
}
