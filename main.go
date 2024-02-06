package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type Book struct {
	isbn   string
	title  string
	author string
	price  float32
}

var db *sql.DB
var err error

func main() {
	db, err = sql.Open("postgres", "postgres://jonny:jonny_go@localhost:5432/bookstore?sslmode=disable")
	if err != nil {
		panic(err)
	}

	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("You successfully connected to Postgres!")

	rows, err := db.Query(`SELECT * FROM books;`)
	if err != nil {
		log.Fatal(err)
	}
	bks := make([]Book, 0)
	for rows.Next() {
		b := Book{}
		if err := rows.Scan(&b.isbn, &b.title, &b.author, &b.price); err != nil {
			log.Fatal(err)
		}
		if err != nil {
			log.Fatal(err)
		}
		bks = append(bks, b)
	}

	if err = rows.Err(); err != nil {
		panic(err)
	}

	for _, b := range bks {
		fmt.Printf("%s, %s, %s, $%.2f\n", b.isbn, b.title, b.author, b.price)
	}

}
