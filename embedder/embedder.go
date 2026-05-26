package embedder

import (
	"context"
	"fmt"
	"log"
	"miniRAGServer-Go/config"

	"google.golang.org/genai"
)

type Embedder struct {
	client *genai.Client
	model  string // ← stored here
}

func New(cfg config.GeminiConfig) (*Embedder, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: cfg.APIKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}
	return &Embedder{client: client, model: cfg.EmbeddingModel}, nil
}

func (e *Embedder) Embed(ctx context.Context, text string) ([]float32, error) {
	log.Printf("Embedding text: %s", text)
	contents := []*genai.Content{
		genai.NewContentFromText(text, genai.RoleUser),
	}
	result, err := e.client.Models.EmbedContent(ctx,
		e.model, // ← use stored model, no cfg needed
		contents,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("embedding failed: %w", err)
	}
	return result.Embeddings[0].Values, nil
}
