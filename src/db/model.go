package db

type UserJSON struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Users struct {
	ID       uint
	Username string
	Password string
	IsAdmin  bool
}
