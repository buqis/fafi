package software_user

type User struct {
	AccountId    string `json:"account_id"`
	AccountToken string `json:"account_token"`
	AgeGroup     int    `json:"age_group"`
	Expiration   string `json:"expiration"`
	TsToken      string `json:"ts_token"`
	Username     string `json:"username"`
}

type Token struct {
	AccountId    string `json:"account_id"`
	AccountToken string `json:"account_token"`
	Uuid         string `json:"uuid"`
}
