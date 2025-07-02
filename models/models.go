package models

type Message struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}
