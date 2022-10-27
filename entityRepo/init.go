package entityRepo

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
)

var Db *sqlx.DB

func init() {
	var err error
	Db, err = sqlx.Open("mysql", "root:2002116yy@tcp(127.0.0.1:3306)/entity_manager2?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
}
