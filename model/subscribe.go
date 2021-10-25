package model

type Subscribe struct {
	ID            uint `gorm:"primaryKey;AUTO_INCREMENT"`
	EndPoint      string
	UserPublicKey string
	UserAuthToken string
}
