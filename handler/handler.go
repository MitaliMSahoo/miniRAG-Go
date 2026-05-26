package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"miniRAGServer-Go/embedder"
	"miniRAGServer-Go/llm"
	"miniRAGServer-Go/model"
	"miniRAGServer-Go/vectordb"
	"net/http"
	"path/filepath"
	"strings"
)

type Handler struct {
	embedder *embedder.Embedder
	vectordb *vectordb.VectorDB
	llm      *llm.Llm
}

func New(e *embedder.Embedder, v *vectordb.VectorDB, l *llm.Llm) *Handler {
	return &Handler{
		embedder: e,
		vectordb: v,
		llm:      l,
	}
}

func (h *Handler) AddDocumentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	contentType := r.Header.Get("Content-Type")
	var added int

	if strings.Contains(contentType, "multipart/form-data") {
		// ---- FILE UPLOAD PATH
		if err := r.ParseMultipartForm(50 << 20); err != nil {
			http.Error(w, "failed to parse form", http.StatusBadRequest)
			return
		}

		files := r.MultipartForm.File["files"]
		if len(files) == 0 {
			http.Error(w, "no files uploaded", http.StatusBadRequest)
			return
		}

		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				http.Error(w, "failed to open file: "+fileHeader.Filename, http.StatusInternalServerError)
				return
			}
			defer file.Close()

			ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
			log.Printf("processing file: %s", fileHeader.Filename)

			// extract text from file
			text, err := embedder.ExtractText(file, ext)

			if err != nil {
				http.Error(w, "extraction failed: "+err.Error(), http.StatusInternalServerError)
				return
			}

			// chunk text
			chunks, err := embedder.ChunkText(text)
			if err != nil {
				http.Error(w, "chunking failed: "+err.Error(), http.StatusInternalServerError)
				return
			}
			log.Printf("file %s split into %d chunks", fileHeader.Filename, len(chunks))

			// embed and store each chunk
			for _, chunk := range chunks {
				vector, err := h.embedder.Embed(ctx, chunk)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				log.Printf("Adding chunk into vector db..")
				if err := h.vectordb.AddDocument(ctx, chunk, vector); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				added++
			}
		}

	} else {
		// ---- JSON PATH
		var reqData model.AddRequest
		if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for _, doc := range reqData.Document {
			vector, err := h.embedder.Embed(ctx, doc.Text)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Printf("Adding Document into vector db..")
			if err := h.vectordb.AddDocument(ctx, doc.Text, vector); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			added++
		}
	}

	log.Printf("Added %d document(s)/chunk(s)", added)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(model.AddResponse{
		Added:   added,
		Message: fmt.Sprintf("%d document(s) added", added),
	})
}
func (h *Handler) QueryPromptHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var reqData model.QueryRequest

	log.Printf("Received query request: %s", reqData.Content)
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	vector, err := h.embedder.Embed(ctx, reqData.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Generated embedding for query: %s", reqData.Content)
	docs, err := h.vectordb.Search(ctx, vector, 3)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("calling LLM with query...")
	prompt, err := h.llm.Query(ctx, reqData.Content, docs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(model.QueryResponse{
		Answer: prompt,
	})

}
