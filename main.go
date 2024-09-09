package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)


type Book struct {
	ID 			string `json:"id"`
	Title 		string `json:"title"`
	Author 		string `json:"author"`
	PublishDate string `json:"publish_date"`
	ISBN 		string `json:"isbn"`
}

type Block struct {
	Pos 		int64 
	Data 		BookCheckout 
	PrevHash 	string 
	Hash 		string 
	Timestamp 	string 
}

type BookCheckout struct {
	BookID 			string `json:"book_id"`
	User 			string `json:"user"`
	CheckoutDate 	string `json:"checkout_date"`
	IsGenesis 		bool `json:"is_genesis"`
}

type Blockchain struct {
	blocks []*Block
}

var blockchain *Blockchain


func newBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Could not decode the request body as Book: %v", err)
		w.Write([]byte("Could not decode the request body as Book"))
		return
	}

	h := md5.New()

	io.WriteString(h, book.ISBN+book.PublishDate)
	book.ID = fmt.Sprintf("%x", h.Sum(nil))

	resp, err := json.MarshalIndent(book, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Could not marshal book object: %v", err)
		w.Write([]byte("Could not marshal book object"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(resp)

}

func main() {
	r := mux.NewRouter()

	r.Handle("/", getBlockchain).Methods("GET")
	r.Handle("/", writeBlock).Methods("POST") 
	r.Handle("/new", newBook).Methods("POST")

	log.Println("Listening on port 3000")
	log.Fatal((http.ListenAndServe(":3000", r)))
}