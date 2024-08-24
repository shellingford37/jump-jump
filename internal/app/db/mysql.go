package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

var mysqlDB *sql.DB

func OpenMysql() error {
	var user = os.Getenv("MYSQL_USER")
	var host = os.Getenv("MYSQL_HOST")
	var dbName = os.Getenv("MYSQL_DBNAME")
	var pwd = os.Getenv("MYSQL_PWD")
	mysqlDB, _ = sql.Open("mysql", user+":"+pwd+"@tcp("+host+")/"+dbName+"?charset=utf8")
	mysqlDB.SetMaxOpenConns(2000)
	mysqlDB.SetMaxIdleConns(1000)
	return mysqlDB.Ping()
}

func GetMySql() (*sql.DB, error) {
	if mysqlDB == nil {
		err := OpenMysql()
		return nil, err
	}
	return mysqlDB, nil
}

//func Query(sql string, args ...any) (*sql.Rows, error) {
//	rows, err := mysqlDB.Query(sql, args)
//	if err != nil {
//		return rows, err
//	}
//	return rows, nil
//}

func ShowMysqlVersion() {
	var version string
	err2 := mysqlDB.QueryRow("SELECT VERSION()").Scan(&version)
	if err2 != nil {
		log.Fatal(err2)
	}
	fmt.Println(version)
}
