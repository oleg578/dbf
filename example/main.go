package main

import (
	"context"
	"database/sql"
	"dbf"
	"encoding/json"
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

	// create result slice
	outRows := make([]map[string]interface{}, 0, numbRows)

	startTime := time.Now()

	// get columns list
	columns, err := rs.Columns()
	if err != nil {
		log.Fatalf("fault columns: %v", err)
	}
	values := make([]interface{}, len(columns))
	valuePointers := make([]interface{}, len(values))
	for i := range values {
		valuePointers[i] = &values[i]
	}
	var counter int64

	for rs.Next() {
		counter++
		if err := rs.Scan(valuePointers...); err != nil {
			log.Fatalf("fault scan: %v", err)
		}
		rowMap, errMsg := dbf.Row2Map(columns, values)
		if errMsg != nil {
			log.Fatalf("fault row2json: %v", errMsg)
		}
		outRows = append(outRows, rowMap)
	}
	if err := rs.Err(); err != nil {
		log.Fatalf("rows row: %v", err)
	}

	endTime := time.Now()

	outData, errData := json.Marshal(outRows)
	if errData != nil {
		log.Fatalf("fault marshal: %v", errData)
	}

	_, _ = os.Stdout.Write(outData)

	_, _ = os.Stderr.WriteString(
		fmt.Sprintf("Elapsed time: %v ms, ", endTime.Sub(startTime).Milliseconds()))
	_, _ = os.Stderr.WriteString(printMemUsage())
	_, _ = os.Stderr.WriteString(
		fmt.Sprintf("%d records read\n", counter))

}
