package handlers

type Register struct {
	Email           string `json:"email" binding:"required,emailValidator"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type ChangePasswordReq struct {
	Password    string `json:"password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type ResetPasswordReq struct {
	Email string `json:"email" binding:"required,emailValidator"`
}

type ResponseMessage struct {
	Error Message
}

type Message struct {
	Code    int
	Message string
}
