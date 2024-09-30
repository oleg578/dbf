package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

const (
	DSN = "root:admin@tcp(127.0.0.1:3307)/test"
)

func fillExampleTable() {
	db, err := sql.Open("mysql", DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	q := fmt.Sprint("INSERT INTO `dummy` VALUES(?,?,?,?,?,?)")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	con, err := db.Conn(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer con.Close()
	tx, errTx := con.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if errTx != nil {
		log.Fatal(errTx)
	}
	defer tx.Rollback()
	stmt, errStmt := tx.Prepare(q)
	if errStmt != nil {
		log.Fatal(errStmt)
	}
	defer stmt.Close()
	for i := 4; i <= 10_000_000; i++ {
		if ((i + 1) % 100_000) == 0 {
			fmt.Printf("%v rows inserted\n", i+1)
		}
		product := struct {
			ID    int
			Name  string
			Desc  interface{}
			Price float64
			Qty   int
			Date  string
		}{
			ID:    i,
			Name:  gofakeit.Product().Name,
			Desc:  nil, //we insert nil value forcibly
			Price: gofakeit.Product().Price,
			Qty:   gofakeit.IntRange(10, 1000),
			Date: gofakeit.DateRange(
				time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC)).
				Format("2006-01-02 15:04:05"),
		}
		if _, err := stmt.ExecContext(ctx,
			product.ID,
			product.Name,
			product.Desc,
			product.Price,
			product.Qty,
			product.Date); err != nil {
			log.Fatal(err)
		}
	}
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
}

func fillTestTable() {
	db, err := sql.Open("mysql", DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	q := "INSERT INTO `dummy` VALUES(?,?,?,?,?,?)"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	con, err := db.Conn(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer con.Close()
	tx, errTx := con.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if errTx != nil {
		log.Fatal(errTx)
	}
	defer tx.Rollback()
	stmt, errStmt := tx.Prepare(q)
	if errStmt != nil {
		log.Fatal(errStmt)
	}
	defer stmt.Close()
	for i := 1; i <= 3; i++ {
		product := struct {
			ID    int
			Name  string
			Desc  interface{}
			Price float64
			Qty   int
			Date  string
		}{
			ID:    i,
			Name:  fmt.Sprintf("product_%d", i),
			Desc:  nil, //we insert nil value forcibly
			Price: float64(i) + .23,
			Qty:   i * 10,
			Date:  "2021-01-01 00:00:00",
		}
		if _, err := stmt.ExecContext(ctx,
			product.ID,
			product.Name,
			product.Desc,
			product.Price,
			product.Qty,
			product.Date); err != nil {
			log.Fatal(err)
		}
	}
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	fillTestTable()
	fillExampleTable()
}
