package initialize

import (
	"database/sql"
	"fmt"
	"go-learning/global"
)

// docker exec -it c930a1205c22 bash
// mysql -uroot -proot123
// use shopdevgo;
// show tables;

// sqlc generate

func InitMysqlC() {
	// TODO: Init mysql from file or env
	m := global.Config.Mysql

	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	var s = fmt.Sprintf(dsn, m.UserName, m.Password, m.Host, m.Port, m.DBName)

	db, err := sql.Open("mysql", s)
	checkErrorPanic(err, "InitMysqlC initialization error")
	global.Logger.Info("Initializing MySQLC successfully")
	global.Mdbc = db
}
