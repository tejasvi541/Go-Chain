package main

import (
	"encoding/json"
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

var Blockchain *Blockchain


func newBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}

	defer r.Body.Close()

	newBlock := Blockchain.AddBlock(book)

	respondWithJSON(w, r, http.StatusCreated, newBlock)
}

func main() {
	r := mux.NewRouter()

	r.Handle("/", getBlockchain).Methods("GET")
	r.Handle("/", writeBlock).Methods("POST") 
	r.Handle("/new", newBook).Methods("POST")

	log.Println("Listening on port 3000")
	log.Fatal((http.ListenAndServe(":3000", r)))
}