package initialize

import (
	"fmt"
	"go-learning/global"
	"go-learning/internal/model"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

// docker exec -it c930a1205c22 bash
// mysql -uroot -proot123
// use shopdevgo;
// show tables;

/* Generate migration
docker exec -it mysql mysqldump -uroot -proot123 --databases shopdevgo --add-drop-databa
se --add-drop-table --add-drop-trigger --add-locks --no-data > migrations/shopdevgo.sql
*/

func checkErrorPanic(err error, errString string) {
	if err != nil {
		global.Logger.Error(errString, zap.Error(err))
		panic(err)
	}
}

func InitMysql() {
	// TODO: Init mysql from file or env
	m := global.Config.Mysql

	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	var s = fmt.Sprintf(dsn, m.UserName, m.Password, m.Host, m.Port, m.DBName)
	db, err := gorm.Open(mysql.Open(s), &gorm.Config{
		SkipDefaultTransaction: false,
	})
	checkErrorPanic(err, "InitMysql initialization error")
	global.Logger.Info("Initializing MySQL successfully")
	global.Mdb = db

	// set Pool
	SetPool()
	// genTableDAO()
	// migrateTables()
}

func SetPool() {
	m := global.Config.Mysql
	sqlDb, err := global.Mdb.DB()
	if err != nil {
		fmt.Printf("Pooling mysql error: %s::\n", err)
	}
	sqlDb.SetMaxIdleConns(m.MaxIdleConns)
	sqlDb.SetConnMaxIdleTime(time.Duration(m.ConnMaxIdleTime))
	sqlDb.SetMaxOpenConns(m.MaxOpenConns)
	sqlDb.SetConnMaxLifetime(time.Duration(m.ConnMaxLifetime))
}

func genTableDAO() {
	// Initiate the tables
	g := gen.NewGenerator(gen.Config{
		OutPath: "./internal/model",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	// gormdb, _ := gorm.Open(mysql.Open("root:@(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"))
	g.UseDB(global.Mdb) // reuse your gorm db
	// g.GenerateAllTable()
	g.GenerateModel("go_db_user")

	// // Generate basic type-safe DAO API for struct `model.User` following conventions
	// g.ApplyBasic(model.User{})

	// // Generate Type Safe API with Dynamic SQL defined on Querier interface for `model.User` and `model.Company`
	// g.ApplyInterface(func(Querier) {}, model.User{}, model.Company{})

	// Generate the code
	g.Execute()
}

func migrateTables() {
	err := global.Mdb.AutoMigrate(
		// &po.User{},
		// &po.Role{},
		&model.GoDbUserV2{},
	)
	if err != nil {
		fmt.Printf("Migrating tables error: %s::\n", err)
	}
}
