package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func ShowMysqlVersion() {
	db, err := sql.Open("mysql", "root:mysql@tcp(o1.shellingford.cn:3306)/test?charset=utf8mb4")
	db.Ping()
	defer db.Close()
	if err != nil {
		fmt.Println("数据库连接失败！")
		log.Fatalln(err)
	}
	var version string
	err2 := db.QueryRow("SELECT VERSION()").Scan(&version)
	if err2 != nil {
		log.Fatal(err2)
	}
	fmt.Println(version)
}
