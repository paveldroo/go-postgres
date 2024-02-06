package main

import (
	"database/sql"
	"fmt"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type Book struct {
	Isbn   string
	Title  string
	Author string
	Price  float32
}

var db *sql.DB
var err error
var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))

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
	r.POST("/books/", insert)
	r.PUT("/books/:id", update)
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
		if err := rows.Scan(&b.Isbn, &b.Title, &b.Author, &b.Price); err != nil {
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

	tpl.ExecuteTemplate(w, "books.gohtml", bks)
}

func get(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	bIsbn := p.ByName("isbn")
	row := db.QueryRow("SELECT * FROM books WHERE Isbn=$1", bIsbn)
	b := Book{}

	err := row.Scan(&b.Isbn, &b.Title, &b.Author, &b.Price)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		log.Fatalln(err)
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	tpl.ExecuteTemplate(w, "show.gohtml", b)
}

func insert(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	b := Book{}

	b.Isbn = r.FormValue("isbn")
	b.Title = r.FormValue("title")
	b.Author = r.FormValue("author")
	p64, err := strconv.ParseFloat(r.FormValue("price"), 32)
	if err != nil {
		http.Error(w, "Wrong Price format", http.StatusUnprocessableEntity)
		return
	}
	b.Price = float32(p64)

	if b.Isbn == "" || b.Title == "" || b.Author == "" {
		http.Error(w, "Wrong book data", http.StatusUnprocessableEntity)
		return
	}

	_, err = db.Exec("INSERT INTO books (isbn, title, author, price) VALUES ($1, $2, $3, $4)", b.Isbn, b.Title, b.Author, b.Price)
	if err != nil {
		log.Fatalln(err)
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
	tpl.ExecuteTemplate(w, "created.gohtml", b)

}

func update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	b := Book{}

	b.Isbn = r.FormValue("isbn")
	b.Title = r.FormValue("title")
	b.Author = r.FormValue("author")
	p64, err := strconv.ParseFloat(r.FormValue("price"), 32)
	if err != nil {
		http.Error(w, "Wrong Price format", http.StatusUnprocessableEntity)
		return
	}
	b.Price = float32(p64)

	if b.Isbn == "" || b.Title == "" || b.Author == "" {
		http.Error(w, "Wrong book data", http.StatusUnprocessableEntity)
		return
	}

	_, err = db.Exec("UPDATE books VALUES ($1, $2, $3, $4)", b.Isbn, b.Title, b.Author, b.Price)
	fmt.Fprintf(w, "Book %s, %s, %s, $%.2f successfully updated!", b.Isbn, b.Title, b.Author, b.Price)
}

func del(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	bIsbn := p.ByName("isbn")
	_, err := db.Exec("DELETE FROM books WHERE Isbn=$1", bIsbn)
	if err != nil {
		http.NotFound(w, r)
	}

	fmt.Fprintf(w, "Book %s successfully deleted!", bIsbn)
}
