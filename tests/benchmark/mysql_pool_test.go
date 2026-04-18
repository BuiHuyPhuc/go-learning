package benchmark

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// setupDB initializes the database with a specific pool size
func setupDB(maxOpen int) *sql.DB {
	db, err := sql.Open("mysql", "root:root123@tcp(127.0.0.1:33306)/shopdevgo")
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxOpen)
	return db
}

func benchmarkDB(b *testing.B, conns int) {
	db := setupDB(conns)
	defer db.Close()

	b.ResetTimer() // Exclude setup time
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Perform a simple query
			var val int
			err := db.QueryRow("SELECT 1").Scan(&val)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// go test -bench=. -benchmem -run=^$

func BenchmarkDBConn1(b *testing.B)  { benchmarkDB(b, 1) }
func BenchmarkDBConn10(b *testing.B) { benchmarkDB(b, 10) }
