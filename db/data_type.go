package db

type UserInfo struct {
	Uid   int    `json:"uid"`
	Email string `json:"email" binding:"required"`
	Pwd   string `json:"pwd" binding:"required"`
}

type User struct {
	UID        int    `gorm:"primaryKey"`
	Email      string `gorm:"unique"`
	Verified   bool
	Authorized bool
	Salt       string
	PwdHash    string
	UUID       []byte
}
