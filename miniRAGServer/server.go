package miniragserver

import (
	"context"
	config "miniRAGServer-Go/config"
	"miniRAGServer-Go/embedder"
	"miniRAGServer-Go/handler"
	"miniRAGServer-Go/llm"
	"miniRAGServer-Go/vectordb"
	"net/http"
)

type Server struct {
	Embedder *embedder.Embedder
	Vectordb *vectordb.VectorDB
	Llm      *llm.Llm
	Handler_ *handler.Handler
}

func ServerInit() *Server {
	cfg := config.Load()

	e, err := embedder.New(cfg.Gemini)
	if err != nil {
		panic(err)
	}

	v, err := vectordb.New(cfg.Weaviate)
	if err != nil {
		panic(err)
	}

	l, err := llm.New(cfg.Gemini)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	v.CreateCollection(ctx, cfg.Weaviate.Collection)

	h := handler.New(e, v, l)

	return &Server{
		Embedder: e,
		Vectordb: v,
		Llm:      l,
		Handler_: h,
	}
}

func (s *Server) Handler() http.Handler {

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./static")))
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	mux.HandleFunc("/adddocument", s.Handler_.AddDocumentHandler)
	mux.HandleFunc("/queryprompt", s.Handler_.QueryPromptHandler)
	return mux
}
