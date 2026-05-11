package main

import (
	"database/sql"
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/xuri/excelize/v2"
)

func main() {
	runtime.GC()

	start := time.Now()
	db, err := sql.Open("mysql", "root:root123@tcp(127.0.0.1:33306)/shopdevgo?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// db.SetMaxOpenConns(20)
	// db.SetMaxIdleConns(10)
	// db.SetConnMaxLifetime(time.Hour)

	// totalRecords := 990000
	// batchSize := 5000

	// err = seedUsers(db, totalRecords, batchSize)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("Insert success %d records in %v\n", totalRecords, time.Since(start))

	rows, err := db.Query("SELECT usr_id, usr_username, usr_email, usr_phone FROM pre_go_crm_user_c")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	elapsed := time.Since(start)
	log.Printf("Query took %s\n", elapsed)

	f := excelize.NewFile()
	streamWriter, err := f.NewStreamWriter("Sheet1")
	if err != nil {
		log.Fatal(err)
	}

	rowIndex := 1
	for rows.Next() {
		var usr_id int
		var usr_username string
		var usr_email string
		var usr_phone string

		err := rows.Scan(&usr_id, &usr_username, &usr_email, &usr_phone)
		if err != nil {
			log.Fatal(err)
		}

		cell, _ := excelize.CoordinatesToCellName(1, rowIndex)
		err = streamWriter.SetRow(cell, []interface{}{usr_id, usr_username, usr_email, usr_phone})
		if err != nil {
			log.Fatal(err)
		}

		rowIndex++
	}

	if err := streamWriter.Flush(); err != nil {
		log.Fatal(err)
	}

	if err := f.SaveAs("output.xlsx"); err != nil {
		log.Fatal(err)
	}

	endElapsed := time.Since(start).Seconds()
	log.Println("Finish... ", endElapsed)

	printBenchmarkInfo("Generate 4 columns * 1.000.000 rows: ", start)
}

func printBenchmarkInfo(info string, startTime time.Time) {
	var memStats runtime.MemStats
	// var rusage syscall.Rusage
	var bToMb = func(b uint64) uint64 {
		return b / 1024 / 1024
	}

	runtime.ReadMemStats(&memStats)

	// Hàm này chỉ hỗ trợ trên Unix/Linux/macOS
	// syscall.Getrusage(syscall.RUSAGE_SELF, &rusage)

	// log.Printf("%s\nRSS = %v\nAlloc = %v MB\nTotalAlloc = %v MB\nSys = %v MB\nNumGC = %v\nElapsed = %.2f seconds\n",
	// 	info, bToMb(uint64(rusage.Maxrss)), bToMb(memStats.Alloc), bToMb(memStats.TotalAlloc), bToMb(memStats.Sys), memStats.NumGC, time.Since(startTime).Seconds())

	// RSS = RAM thực tế process đang chiếm trong hệ điều hành
	// Alloc = memory đang dùng bởi heap Go
	// TotalAlloc = tổng memory đã cấp phát từ lúc start
	// Sys = memory Go xin từ OS
	// NumGC = số lần GC
	log.Printf("%s\nAlloc = %v MB\nTotalAlloc = %v MB\nSys = %v MB\nNumGC = %v\nElapsed = %.2f seconds\n",
		info, bToMb(memStats.Alloc), bToMb(memStats.TotalAlloc), bToMb(memStats.Sys), memStats.NumGC, time.Since(startTime).Seconds())
}

// Kiểm tra size cơ sở dữ liệu
// SELECT
//
//	table_schema AS database_name,
//	ROUND(SUM(data_length + index_length) / 1024 / 1024, 2) AS size_mb,
//	ROUND(SUM(data_length + index_length) / 1024 / 1024 / 1024, 2) AS size_gb
//
// FROM information_schema.tables
// WHERE table_schema = 'shopdevgo'
// GROUP BY table_schema;

// SELECT
//
//	table_schema AS db_name,
//	table_name,
//	table_rows,
//	ROUND(data_length / 1024 / 1024, 2) AS data_mb,
//	ROUND(index_length / 1024 / 1024, 2) AS index_mb,
//	ROUND((data_length + index_length) / 1024 / 1024, 2) AS total_mb
//
// FROM information_schema.tables
// WHERE table_schema = 'shopdevgo'
// ORDER BY total_mb DESC;
func seedUsers(db *sql.DB, total int, batchSize int) error {

	// Disable checks for faster insert
	_, err := db.Exec("SET autocommit = 0")
	if err != nil {
		return err
	}

	_, err = db.Exec("SET unique_checks = 0")
	if err != nil {
		return err
	}

	_, err = db.Exec("SET foreign_key_checks = 0")
	if err != nil {
		return err
	}

	now := time.Now().Unix()

	for i := 0; i < total; i += batchSize {

		tx, err := db.Begin()
		if err != nil {
			return err
		}

		valueStrings := make([]string, 0, batchSize)
		valueArgs := make([]interface{}, 0, batchSize*11)

		for j := 0; j < batchSize && (i+j) < total; j++ {

			index := i + j

			email := fmt.Sprintf("user%d@gmail.com", index)
			phone := fmt.Sprintf("090%08d", index)
			username := fmt.Sprintf("user_%d", index)

			valueStrings = append(valueStrings, "(?,?,?,?,?,?,?,?,?,?,?)")

			valueArgs = append(valueArgs,
				email,
				phone,
				username,
				"123456",
				now,
				now,
				"127.0.0.1",
				now,
				"127.0.0.1",
				0,
				1,
			)
		}

		query := fmt.Sprintf(`
			INSERT INTO pre_go_crm_user_c (
				usr_email,
				usr_phone,
				usr_username,
				usr_password,
				usr_created_at,
				usr_updated_at,
				usr_create_ip_at,
				usr_last_login_at,
				usr_last_login_ip_at,
				usr_login_times,
				usr_status
			)
			VALUES %s
		`, strings.Join(valueStrings, ","))

		_, err = tx.Exec(query, valueArgs...)
		if err != nil {
			tx.Rollback()
			return err
		}

		err = tx.Commit()
		if err != nil {
			return err
		}

		if i%batchSize == 0 {
			log.Printf(
				"Inserted %d/%d users (%.2f%%)",
				i,
				total,
				float64(i)*100/float64(total),
			)
		}
	}

	// Re-enable checks
	_, err = db.Exec("SET unique_checks = 1")
	if err != nil {
		return err
	}

	_, err = db.Exec("SET foreign_key_checks = 1")
	if err != nil {
		return err
	}

	return nil
}
