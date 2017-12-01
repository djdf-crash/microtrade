package db

type UserJSON struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Users struct {
	ID       uint
	Email    string
	Password string
	IsAdmin  bool
}
