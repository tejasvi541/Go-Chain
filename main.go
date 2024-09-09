package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Book represents a book with metadata.
type Book struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	PublishDate string `json:"publish_date"`
	ISBN        string `json:"isbn"`
}

// BookCheckout represents a book checkout record.
type BookCheckout struct {
	BookID       string `json:"book_id"`
	User         string `json:"user"`
	CheckoutDate string `json:"checkout_date"`
	IsGenesis    bool   `json:"is_genesis"`
}

// Block represents a single block in the blockchain.
type Block struct {
	Position   int64        `json:"position"`
	Timestamp  string       `json:"timestamp"`
	Data       BookCheckout `json:"data"`
	PrevHash   string       `json:"previous_hash"`
	Hash       string       `json:"hash"`
}

// Blockchain represents a chain of blocks.
type Blockchain struct {
	Blocks []*Block `json:"blocks"`
}

// Global variable for the blockchain instance.
var blockchain *Blockchain

// GenerateHash computes the SHA-256 hash of the block's data.
func (b *Block) GenerateHash() {
	dataBytes, _ := json.Marshal(b.Data)
	data := fmt.Sprintf("%d%s%s%s", b.Position, b.Timestamp, string(dataBytes), b.PrevHash)
	hash := sha256.New()
	hash.Write([]byte(data))
	b.Hash = hex.EncodeToString(hash.Sum(nil))
}

// CreateBlock constructs a new block from the previous block and checkout data.
func CreateBlock(prevBlock *Block, checkoutData BookCheckout) *Block {
	newBlock := &Block{
		Position:   prevBlock.Position + 1,
		Timestamp:  time.Now().Format(time.RFC3339),
		Data:       checkoutData,
		PrevHash:   prevBlock.Hash,
	}
	newBlock.GenerateHash()
	return newBlock
}

// AddBlock appends a new block to the blockchain if valid.
func (bc *Blockchain) AddBlock(data BookCheckout) {
	lastBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := CreateBlock(lastBlock, data)

	if IsValidBlock(newBlock, lastBlock) {
		bc.Blocks = append(bc.Blocks, newBlock)
	}
}

// GenesisBlock creates the initial block of the blockchain.
func GenesisBlock() *Block {
	return CreateBlock(&Block{}, BookCheckout{IsGenesis: true})
}

// NewBlockchain initializes a new blockchain with a genesis block.
func NewBlockchain() *Blockchain {
	return &Blockchain{Blocks: []*Block{GenesisBlock()}}
}

// IsValidBlock checks if a new block is valid compared to the previous block.
func IsValidBlock(newBlock, previousBlock *Block) bool {
	if previousBlock.Hash != newBlock.PrevHash {
		return false
	}
	if !newBlock.ValidateHash(newBlock.Hash) {
		return false
	}
	if previousBlock.Position+1 != newBlock.Position {
		return false
	}
	return true
}

// ValidateHash verifies if the block's hash is valid.
func (b *Block) ValidateHash(expectedHash string) bool {
	b.GenerateHash()
	return b.Hash == expectedHash
}

// GetBlockchainHandler returns the current blockchain as JSON.
func GetBlockchainHandler(w http.ResponseWriter, r *http.Request) {
	responseBytes, err := json.MarshalIndent(blockchain.Blocks, "", "  ")
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Failed to marshal blockchain", err)
		return
	}
	_, err = io.WriteString(w, string(responseBytes))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

// WriteBlockHandler adds a new block to the blockchain.
func WriteBlockHandler(w http.ResponseWriter, r *http.Request) {
	var checkoutItem BookCheckout
	if err := json.NewDecoder(r.Body).Decode(&checkoutItem); err != nil {
		handleError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	blockchain.AddBlock(checkoutItem)
	responseBytes, err := json.MarshalIndent(checkoutItem, "", "  ")
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Failed to marshal response", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}

// NewBookHandler adds a new book entry and returns it as JSON.
func NewBookHandler(w http.ResponseWriter, r *http.Request) {
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		handleError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	hash := md5.New()
	io.WriteString(hash, book.ISBN+book.PublishDate)
	book.ID = fmt.Sprintf("%x", hash.Sum(nil))
	responseBytes, err := json.MarshalIndent(book, "", "  ")
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Failed to marshal response", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(responseBytes)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

// handleError handles HTTP errors and logs them.
func handleError(w http.ResponseWriter, statusCode int, message string, err error) {
	w.WriteHeader(statusCode)
	log.Printf("%s: %v", message, err)
	_, writeErr := w.Write([]byte(message))
	if writeErr != nil {
		log.Printf("Error writing error message: %v", writeErr)
	}
}

// main sets up the HTTP server and routes.
func main() {
	blockchain = NewBlockchain()

	router := mux.NewRouter()
	router.HandleFunc("/", GetBlockchainHandler).Methods("GET")
	router.HandleFunc("/", WriteBlockHandler).Methods("POST")
	router.HandleFunc("/new", NewBookHandler).Methods("POST")

	log.Println("Server started on port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
