package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// Book struct (Model)
type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

// Author struct
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

const (
	host     = "ec2-54-243-253-181.compute-1.amazonaws.com"
	port     = 5432
	user     = "hkuncdjazivhgt"
	password = "a6211ecc634a041ac919d9d0dc50a8a0ec9acae6448be9baee95b4a472c8a152"
	dbname   = "dboqmi8lc7fu3t"
)

// Init books var as a slice Book struct
var books []Book
var psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
	"password=%s dbname=%s sslmode=require",
	host, port, user, password, dbname)

// Get all books
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
	fmt.Println("Asked for books")
}

func getDBBooks(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("Successfully connected!")

	tsql := fmt.Sprintf("SELECT * FROM books;")
	rows, err := db.Query(tsql)
	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		//return -1, err
	}
	defer rows.Close()
	count := 0
	for rows.Next() {
		var Title, Id string

		err := rows.Scan(&Id, &Title)
		if err != nil {
			fmt.Println("Error reading rows: " + err.Error())
			//return -1, err
		}
		//fmt.Printf("%s \t %s \t %s \n", ReceiveTime, Message, DT)
		fmt.Printf("%s  %s \n", Id, Title)
		books[count].Title = Title
		count++
	}

	db.Close()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
	fmt.Println("Asked for books")

}

// Get single book
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Gets params
	// Loop through books and find one with the id from the params
	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}

// Add new book
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(100000000)) // Mock ID - not safe
	books = append(books, book)
	json.NewEncoder(w).Encode(book)
}

// Update book
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = params["id"]
			books = append(books, book)
			json.NewEncoder(w).Encode(book)
			return
		}
	}
}

// Delete book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}

func index(w http.ResponseWriter, r *http.Request) {
	s2 := template.Must(template.ParseFiles("about.html"))
	s2.Execute(w, nil)
	fmt.Println("about executed")
	//log.Println("about executed")
}

func connDB() {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("Successfully connected!")

	tsql := fmt.Sprintf("SELECT * FROM books;")
	rows, err := db.Query(tsql)
	if err != nil {
		fmt.Println("Error reading rows: " + err.Error())
		//return -1, err
	}
	defer rows.Close()
	count := 0
	for rows.Next() {
		var Title, Id string

		err := rows.Scan(&Id, &Title)
		if err != nil {
			fmt.Println("Error reading rows: " + err.Error())
			//return -1, err
		}
		//fmt.Printf("%s \t %s \t %s \n", ReceiveTime, Message, DT)
		fmt.Printf("%s  %s \n", Id, Title)
		count++
	}
	//return count, nil
}

// Main function
func main() {
	connDB()
	// Init router
	r := mux.NewRouter()

	// Hardcoded data - @todo: add database
	books = append(books, Book{ID: "1", Isbn: "438227", Title: "Book One", Author: &Author{Firstname: "John", Lastname: "Doe"}})
	books = append(books, Book{ID: "2", Isbn: "454555", Title: "Book Two", Author: &Author{Firstname: "Steve", Lastname: "Smith"}})

	// Route handles & endpoints

	r.HandleFunc("/", index)
	r.HandleFunc("/books", getBooks).Methods("GET")
	r.HandleFunc("/dbbooks", getDBBooks).Methods("GET")
	r.HandleFunc("/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/books", createBook).Methods("POST")
	r.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	port := os.Getenv("PORT")
	//port := "3000"

	fmt.Printf("Using port %s\n", port)
	// Start server
	log.Fatal(http.ListenAndServe(":"+port, r))
	fmt.Printf("Using port %s\n", port)
}
