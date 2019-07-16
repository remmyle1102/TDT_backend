package models

// Account represent the account model
type Account struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
	RoleID   int    `json:"roleID"`
	Status   bool   `json:"status"`
}
