package user

type User struct {
	Id          string      `json:"id"`
	AccountType AccountType `json:"accountType"`
	Username    string      `json:"username"`
	IsVerified  bool        `json:"isVerified"`
}
