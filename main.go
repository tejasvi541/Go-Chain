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

var Blockchain *Blockchain

func writeBlock(w http.ResponseWriter, r *http.Request) {
	// Decode the incoming request body into a BookCheckout struct.
	var checkoutItem BookCheckout
	if err := json.NewDecoder(r.Body).Decode(&checkoutItem); err != nil {
		handleError(w, http.StatusInternalServerError, "Could not decode request body", err)
		return
	}

	// Add the new block to the blockchain.
	BlockChain.AddBlock(checkoutItem)

	// Marshal the checkout item into a JSON response.
	resp, err := json.MarshalIndent(checkoutItem, "", " ")
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Could not marshal response", err)
		return
	}

	// Write the response with a 200 OK status.
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}


func newBook(w http.ResponseWriter, r *http.Request) {
	var book Book

	// Decode the JSON request body into the Book struct.
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Could not decode the request body as Book", err)
		return
	}

	// Create a unique ID for the book using MD5 hash.
	h := md5.New()
	io.WriteString(h, book.ISBN+book.PublishDate)
	book.ID = fmt.Sprintf("%x", h.Sum(nil))

	// Marshal the book object into a JSON response.
	resp, err := json.MarshalIndent(book, "", " ")
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Could not marshal book object", err)
		return
	}

	// Write the successful response.
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(resp)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

// handleError is a helper function for handling errors consistently.
func handleError(w http.ResponseWriter, statusCode int, message string, err error) {
	w.WriteHeader(statusCode)
	log.Printf("%s: %v", message, err)
	_, writeErr := w.Write([]byte(message))
	if writeErr != nil {
		log.Printf("Error writing error message to response: %v", writeErr)
	}
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", getBlockchain).Methods("GET")
	r.HandleFunc("/", writeBlock).Methods("POST") 
	r.HandleFunc("/new", newBook).Methods("POST")

	log.Println("Listening on port 3000")
	log.Fatal((http.ListenAndServe(":3000", r)))
}