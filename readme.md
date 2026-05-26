# miniRAG Server - Go

A production-ready **Retrieval Augmented Generation (RAG)** server built in Go. It combines document embeddings, vector search, and large language models to answer questions based on your custom documents.

## 🎯 What It Does

miniRAG Server enables:
- **Document Ingestion** — Upload documents and automatically embed them using Google Gemini
- **Vector Search** — Store embeddings in Weaviate and find similar documents instantly
- **RAG Pipeline** — Generate accurate answers by combining retrieved context with Gemini LLM
- **REST API** — Easy-to-use HTTP endpoints for adding documents and querying
- **Web UI** — Interactive interface for document management and querying

**Example Flow:**
```
User asks: "What is Go?"
    ↓
Server embeds the question (Gemini)
    ↓
Search for similar documents (Weaviate)
    ↓
Send question + context to Gemini LLM
    ↓
Return AI-generated answer
```

---

## 🏗️ Architecture & Components

### Core Packages

| Package | Purpose |
|---------|---------|
| **`embedder/`** | Converts text to vector embeddings using Google Gemini API |
| **`vectordb/`** | Manages vector storage and similarity search with Weaviate |
| **`llm/`** | Generates responses using Google Gemini LLM with RAG context |
| **`handler/`** | HTTP request handlers for `/adddocument` and `/queryprompt` endpoints |
| **`config/`** | Loads configuration from environment variables |
| **`model/`** | Defines data structures (requests, responses, documents) |
| **`miniRAGServer/`** | Server initialization and route setup |
| **`static/`** | Web UI (HTML/CSS/JS) |

### External Services

| Service | Purpose |
|---------|---------|
| **Weaviate** | Open-source vector database for storing document embeddings |
| **Google Gemini API** | Provides embedding and LLM capabilities |

---

## 📋 Requirements

### Standalone
- Go 1.26+
- Docker (for Weaviate)
- Google Gemini API key

### Docker
- Docker & Docker Compose
- Google Gemini API key

---

## 🚀 Quick Start

### Option 1: Docker (Recommended)

**1. Clone and configure:**
```bash
cd miniRAGServer-Go
cp .env.example .env
# Edit .env and add your GEMINI_API_KEY
```

**2. Start all services:**
```bash
docker-compose up -d
```

**3. Access the application:**
- **Web UI:** http://localhost:80 (or http://localhost:8080 without nginx)
- **Weaviate Dashboard:** http://localhost:8081

**4. Stop services:**
```bash
docker-compose down
```

**Using Makefile (easier):**
```bash
make up          # Start all services
make logs        # View logs
make test        # Run test queries
make down        # Stop services
```

---

### Option 2: Standalone (Manual)

**1. Install Weaviate:**
```bash
docker run -d -p 8080:8080 \
  -e QUERY_DEFAULTS_LIMIT=25 \
  -e AUTHENTICATION_ANONYMOUS_ACCESS_ENABLED=true \
  semitechnologies/weaviate:latest
```

**2. Configure environment:**
```bash
cp .env.example .env
# Edit .env with your settings:
# - GEMINI_API_KEY (required)
# - WEAVIATE_HOST=localhost
# - WEAVIATE_PORT=8080
```

**3. Update main.go to load .env:**
```go
package main

import (
    "github.com/joho/godotenv"
    // ... other imports
)

func main() {
    godotenv.Load()  // Load .env file
    cfg := config.Load()
    s := miniragserver.ServerInit()
    // ...
}
```

**4. Run the server:**
```bash
go mod tidy
go run main.go
```

**5. Server is running:**
```
Server starting on 8080
```

---

## 🔑 Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | `8080` | Server listen port |
| `SERVER_IP` | `0.0.0.0` | Server bind address |
| `GEMINI_API_KEY` | (required) | Google Gemini API key |
| `GEMINI_EMBEDDING_MODEL` | `models/embedding-001` | Embedding model |
| `GEMINI_LLM_MODEL` | `gemini-2.0-flash` | LLM model |
| `WEAVIATE_HOST` | `localhost` | Weaviate server host |
| `WEAVIATE_PORT` | `8080` | Weaviate server port |
| `WEAVIATE_SCHEME` | `http` | Connection scheme (http/https) |
| `WEAVIATE_COLLECTION` | `Documents` | Weaviate collection name |

### Get Your Gemini API Key

1. Visit [Google AI Studio](https://aistudio.google.com/app/apikey)
2. Click "Create API Key"
3. Copy the key to your `.env` file

---

## 📡 API Endpoints

### 1. Add Documents
**Endpoint:** `POST /adddocument`

**Request:**
```json
{
  "documents": [
    {"text": "Go is a compiled language"},
    {"text": "Go has built-in concurrency with goroutines"}
  ]
}
```

**Response:**
```json
{"added": 2, "message": "2 document(s) added"}
```

**Example:**
```bash
curl -X POST http://localhost:8080/adddocument \
  -H "Content-Type: application/json" \
  -d '{
    "documents": [
      {"text": "Go is a compiled language"},
      {"text": "Go has built-in concurrency"}
    ]
  }'
```

---

### 2. Query with RAG
**Endpoint:** `POST /queryprompt`

**Request:**
```json
{"content": "What is Go?"}
```

**Response:**
```json
{"answer": "Go is a compiled, statically-typed programming language developed by Google..."}
```

**Example:**
```bash
curl -X POST http://localhost:8080/queryprompt \
  -H "Content-Type: application/json" \
  -d '{"content": "What is Go?"}'
```

---

### 3. Web UI
**URL:** `GET /`

Open http://localhost:8080 in your browser to use the interactive UI.

---

## 🔧 Weaviate Management

### Verify Weaviate is Running
```bash
curl http://localhost:8080/v1/.well-known/live
```

### Check Collections
```bash
curl http://localhost:8080/v1/schema
```

### Delete a Collection (⚠️ Deletes all documents)
```bash
curl -X DELETE http://localhost:8080/v1/schema/Documents
```

---

## 📁 Project Structure

```
miniRAGServer-Go/
├── main.go                  # Entry point, loads config & starts server
├── config/
│   └── config.go            # Environment config loading
├── embedder/
│   └── embedder.go          # Text → vectors (Gemini)
├── vectordb/
│   └── vectordb.go          # Vector storage & search (Weaviate)
├── llm/
│   └── llm.go               # Answer generation (Gemini)
├── handler/
│   └── handler.go           # HTTP handlers
├── miniRAGServer/
│   └── server.go            # Server setup & routing
├── model/
│   └── model.go             # Data structures
├── static/
│   └── index.html           # Web UI
├── Dockerfile               # Container image
├── docker-compose.yml       # Multi-container orchestration
├── .env.example             # Configuration template
└── go.mod / go.sum          # Dependency management
```

---

## 🛠️ Development

### Live Reload
Install `air` for automatic rebuild on file changes:

```bash
go install github.com/cosmtrek/air@latest
air
```

### Build Locally
```bash
go build -o miniragserver ./main.go
./miniragserver
```
