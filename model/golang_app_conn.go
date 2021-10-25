package model

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewGolangAppConnection() (db *gorm.DB, err error) {
	golang_pass := os.Getenv("GOLANG_DB_PASS")
	conn_str := fmt.Sprintf("golang_conn:%s@/golang_app?&parseTime=True&loc=Local", golang_pass)
	db, err = gorm.Open(mysql.Open(conn_str), &gorm.Config{})
	return
}
