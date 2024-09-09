# Blockchain Book Checkout System

## Overview

This project is a simple yet robust blockchain implementation for managing book checkouts using a RESTful API. The application allows users to:

- **Retrieve** the current state of the blockchain.
- **Add new blocks** representing book checkout transactions.
- **Register new books** with unique identifiers.

![Blockchain Illustration](https://upload.wikimedia.org/wikipedia/commons/6/6f/Blockchain-illustration.png)

## Technical Design

The core of this system is a blockchain that ensures the integrity and immutability of book checkout records. The blockchain is built using fundamental cryptographic principles, such as hashing and block linking, to maintain a secure, verifiable ledger of transactions.

### Blockchain Components

1. **Blocks**: A block represents a single unit of data, containing:

   - `Position`: The block's sequence number in the chain.
   - `Timestamp`: The time when the block was created.
   - `Data`: The book checkout information.
   - `PrevHash`: The hash of the previous block in the chain.
   - `Hash`: The unique SHA-256 hash of the current block.

2. **Blockchain**: A chain of blocks where each block is linked to its predecessor by storing the previous block's hash. This linking forms a tamper-evident chain, ensuring any modification in one block invalidates all subsequent blocks.

3. **Genesis Block**: The initial block in the blockchain, created manually or automatically when the blockchain is instantiated. It has no predecessor, and its previous hash is typically set to zero.

### Core Cryptographic Concepts

- **SHA-256 Hashing**: Each block's data is hashed using SHA-256 to produce a unique, fixed-size output. This hash is used to verify data integrity and establish links between blocks.
- **MD5 Hashing for Book IDs**: To ensure unique identification of each book, an MD5 hash is generated from the book's ISBN and publish date.

![SHA-256 Hashing](https://upload.wikimedia.org/wikipedia/commons/2/29/SHA-2.svg)

## How It Works

### Blockchain Initialization

The blockchain is initialized with a genesis block that serves as the foundation. Any subsequent blocks will be appended to this genesis block.

### Adding Blocks

A new block is created by providing a `BookCheckout` data payload. The following steps are performed:

1. **Generate Hash**: The block's data, timestamp, and the previous block's hash are combined and hashed using SHA-256.
2. **Block Creation**: A new block is instantiated with the generated hash and other relevant details.
3. **Validation**: The new block is validated against the last block in the chain. This ensures:
   - Correct sequence (`Position`).
   - Hash integrity (`PrevHash` matches the previous block's hash).
   - Data integrity (`Hash` matches the computed hash).

If all validations pass, the new block is added to the blockchain.

![Blockchain Process](https://upload.wikimedia.org/wikipedia/commons/thumb/5/58/Blockchain_system.svg/800px-Blockchain_system.svg.png)

### Registering New Books

New books can be registered with unique IDs generated using an MD5 hash of their ISBN and publish date. This ensures that each book has a unique and identifiable key within the system.

## Data Structures

The following Go structs define the core entities:

- **`Book`**: Represents a book's metadata.
- **`BookCheckout`**: Represents a record of a book checkout, including user information and checkout date.
- **`Block`**: Represents an individual block in the blockchain.
- **`Blockchain`**: Represents the blockchain itself, a sequence of blocks.

### Functions

- **`GenerateHash`**: Computes the SHA-256 hash for a block, used for data integrity.
- **`CreateBlock`**: Constructs a new block using data from the previous block and the new checkout transaction.
- **`AddBlock`**: Adds a validated new block to the blockchain.
- **`GenesisBlock`**: Creates the initial block (genesis block) for the blockchain.
- **`NewBlockchain`**: Initializes a new blockchain instance.
- **`IsValidBlock`**: Validates a block against its predecessor.
- **`ValidateHash`**: Checks if a block's hash matches the expected computed hash.

### API Endpoints

#### `GET /`

Retrieves the current state of the blockchain.

- **Response:**
  - `200 OK`: JSON array representing the blockchain.

#### `POST /`

Adds a new book checkout record to the blockchain.

- **Request Body:**

  - JSON object representing `BookCheckout`.

- **Response:**
  - `200 OK`: JSON object of the added checkout record.

#### `POST /new`

Registers a new book.

- **Request Body:**

  - JSON object representing `Book`.

- **Response:**
  - `201 Created`: JSON object of the registered book.

## Running the Server

### Prerequisites

Ensure you have Go installed on your machine.

### Starting the Server

Run the following command to start the server:

````bash
go run main.go


- Adding a Block

```bash
curl -X POST http://localhost:3000 -H "Content-Type: application/json" -d '{"book_id": "123", "user": "John Doe", "checkout_date": "2024-09-08"}'
````

- Registering a New Book

```bash
curl -X POST http://localhost:3000/new -H "Content-Type: application/json" -d '{"title": "Go Programming", "author": "Jane Doe", "publish_date": "2024-09-08", "isbn": "1234567890"}'
```
