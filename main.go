package main

import (
	"database/sql"
	"fmt"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

type Book struct {
	isbn   string
	title  string
	author string
	price  float32
}

var db *sql.DB
var err error

func init() {
	db, err = sql.Open("postgres", "postgres://jonny:jonny_go@localhost:5432/bookstore?sslmode=disable")
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("You successfully connected to Postgres!")
}

func main() {
	r := httprouter.New()
	r.GET("/books", getAll)
	r.GET("/books/:isbn", get)
	r.POST("/books/", create)
	r.DELETE("/books/:id", del)
	log.Fatal(http.ListenAndServe(":8000", r))

}

func getAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	rows, err := db.Query("SELECT * FROM books;")
	if err != nil {
		log.Fatalln(err)
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	bks := make([]Book, 0)
	for rows.Next() {
		b := Book{}
		if err := rows.Scan(&b.isbn, &b.title, &b.author, &b.price); err != nil {
			if err != nil {
				log.Fatalln(err)
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				return
			}
		}
		bks = append(bks, b)
	}

	if err = rows.Err(); err != nil {
		if err != nil {
			log.Fatalln(err)
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
	}

	for _, b := range bks {
		fmt.Fprintf(w, "%s, %s, %s, $%.2f\n", b.isbn, b.title, b.author, b.price)
	}
}

func get(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	bIsbn := p.ByName("isbn")
	row := db.QueryRow("SELECT * FROM books WHERE isbn=$1", bIsbn)
	b := Book{}

	err := row.Scan(&b.isbn, &b.title, &b.author, &b.price)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		log.Fatalln(err)
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s, %s, %s, $%.2f\n", b.isbn, b.title, b.author, b.price)
}

func create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func del(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}
