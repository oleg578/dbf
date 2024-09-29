package main

import (
	"bufio"
	"context"
	"database/sql"
	"dbf"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: db2json <number-rows>")
		os.Exit(1)
	}
	numbRows, _ := strconv.ParseInt(os.Args[1], 10, 64)
	db, err := sql.Open("mysql", "root:admin@tcp(127.0.0.1:3307)/test")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	q := "SELECT * FROM `dummy` LIMIT ?"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	con, errCon := db.Conn(ctx)
	if errCon != nil {
		log.Fatal(errCon)
	}
	defer con.Close()

	rs, errRs := con.QueryContext(ctx, q, numbRows)
	if errRs != nil {
		log.Fatalf("query: %v", errRs)
	}
	defer rs.Close()

	// create new buffer
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	writer.WriteByte('[')
	startTime := time.Now()

	// get a columnTypes list
	columns, err := rs.ColumnTypes()
	if err != nil {
		log.Fatalf("fault get column types: %v", err)
	}
	values := make([]sql.RawBytes, len(columns))
	valuePointers := make([]interface{}, len(values))
	for i := range values {
		valuePointers[i] = &values[i]
	}
	var counter int64
	rs.Next()
	for {
		counter++
		if err := rs.Scan(valuePointers...); err != nil {
			log.Fatalf("fault scan: %v", err)
		}
		msg, errMsg := dbf.Row2Json(columns, values)
		if errMsg != nil {
			log.Fatalf("fault row2json: %v", errMsg)
		}
		if _, err := writer.WriteString(msg); err != nil {
			log.Fatalf("fault write row: %v", err)
		}
		if rs.Next() {
			writer.WriteByte(',')
		} else {
			break
		}
	}
	writer.WriteByte(']')
	if err := rs.Err(); err != nil {
		log.Fatalf("rows row: %v", err)
	}
	_, _ = os.Stderr.WriteString(
		fmt.Sprintf("Elapsed time: %v ms, ", time.Since(startTime).Milliseconds()))
	_, _ = os.Stderr.WriteString(printMemUsage())
	_, _ = os.Stderr.WriteString(
		fmt.Sprintf("%d records read\n", counter))
	if err := writer.Flush(); err != nil {
		log.Fatalf("flush out error: %v", err)
	}
}
