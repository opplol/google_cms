package controller

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type TestController struct{}

func (h TestController) Show(c *gin.Context) {
	golang_pass := os.Getenv("GOLANG_DB_PASS")
	conn_str := fmt.Sprintf("golang_conn:%s@/golang_app", golang_pass)
	db, err := sql.Open("mysql", conn_str)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	fmt.Printf("%#v\n", db)

	// insert test
	stmtIns, err := db.Prepare("insert into test_table values(?, ?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(1, "myname", "mymemo", time.Now(), time.Now())
	if err != nil {
		panic(err.Error())
	}

	// select test
	rows, err := db.Query("select * from test_table")
	if err != nil {
		panic(err.Error())
	}

	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			fmt.Println(columns[i], ": ", value)
		}
		fmt.Println("-----------------------------------")
	}
	if err = rows.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	fmt.Printf("%#v\n", values)

	c.JSON(http.StatusOK, gin.H{"info": db})
}

type TestTable struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Memo      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
type Tabler interface {
	TableName() string
}

// TableName overrides the table name used by User to `profiles`
func (TestTable) TableName() string {
	return "test_table"
}

func (h TestController) OrmTest(c *gin.Context) {
	golang_pass := os.Getenv("GOLANG_DB_PASS")
	conn_str := fmt.Sprintf("golang_conn:%s@/golang_app?&parseTime=True&loc=Local", golang_pass)
	db, err := gorm.Open(mysql.Open(conn_str), &gorm.Config{})

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	my_model := TestTable{ID: 2, Name: "orm test", Memo: "orm Memo"}
	result := db.Create(&my_model)
	fmt.Printf("%#v\n", my_model.ID)
	fmt.Printf("%#v\n", result.Error)
	fmt.Printf("%#v\n", result.RowsAffected)

	var all_record []TestTable
	db.Find(&all_record)
	c.JSON(http.StatusOK, gin.H{"info": all_record})
}
