package types

type User struct {
	ID       string `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
}

type UserRequest struct {
	Username string `json:"username"`
	Password string `password:"password"`
}
