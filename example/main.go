package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	const DSN = "root:admin@tcp(127.0.0.1:3307)/test"

	if len(os.Args) < 2 {
		fmt.Println("Usage: db2json <number-rows>")
		os.Exit(1)
	}
	numbRows, _ := strconv.ParseInt(os.Args[1], 10, 64)
	var result []Dummy
	db, err := sql.Open("mysql", DSN)
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

	startTime := time.Now()
	for rs.Next() {
		var dmy = Dummy{}
		if err := rs.Scan(
			&dmy.ID,
			&dmy.Product,
			&dmy.Description,
			&dmy.Price,
			&dmy.Qty,
			&dmy.Date); err != nil {
			panic(err)
		}
		result = append(result, dmy)
	}
	if err := rs.Err(); err != nil {
		panic(err)
	}
	msg, errRTJ := json.Marshal(result)
	if errRTJ != nil {
		panic(errRTJ)
	}

	endTime := time.Now()

	_, errOut := os.Stdout.Write(msg)
	if errOut != nil {
		panic(errOut)
	}
	_, _ = os.Stderr.WriteString(fmt.Sprintf("Elapsed time: %v ms, ", endTime.Sub(startTime).Milliseconds()))
	_, _ = os.Stderr.WriteString(printMemUsage())
	_, _ = os.Stderr.WriteString(fmt.Sprintf("%d records read\n", len(result)))

}
