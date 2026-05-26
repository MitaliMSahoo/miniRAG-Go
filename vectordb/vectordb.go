package vectordb

import (
	"context"
	"fmt"
	"log"
	"miniRAGServer-Go/config"

	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/auth"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
)

type VectorDB struct {
	weaviateClient *weaviate.Client
}

func New(cfg config.WeaviateConfig) (*VectorDB, error) {
	connection, err := weaviate.NewClient(weaviate.Config{
		Host:       cfg.Host,
		Scheme:     cfg.Scheme,
		AuthConfig: auth.ApiKey{Value: cfg.APIKey},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to weaviate DB: %w", err)
	}

	return &VectorDB{
		weaviateClient: connection,
	}, nil
}

func (v *VectorDB) CreateCollection(ctx context.Context, documentName string) (string, error) {
	// 	create a collection/class in Weaviate to store your documents
	// think of it like creating a table in a regular DB
	// called once at startup before anything is stored
	exists, err := v.weaviateClient.Schema().ClassExistenceChecker().
		WithClassName(documentName).
		Do(ctx)

	if err != nil {
		return "", fmt.Errorf("checking collection existence: %w", err)
	}
	if exists {
		log.Println("collection '" + documentName + "' already exists, skipping creation")
		return "exists", nil
	}
	// Define the collection
	classObj := &models.Class{
		Class:      "Documents",
		Vectorizer: "none", // we supply our own vectors from embedder.go
		Properties: []*models.Property{
			{
				Name:     "text",
				DataType: []string{"text"},
			},
		},
	}

	err = v.weaviateClient.Schema().ClassCreator().WithClass(classObj).Do(ctx)
	if err != nil {
		return "", fmt.Errorf("creating collection: %w", err)
	}

	log.Println("collection 'Documents' created successfully")
	return "created", nil
}

func (v *VectorDB) AddDocument(ctx context.Context, text string, vector []float32) error {
	_, err := v.weaviateClient.Data().Creator().
		WithClassName("Documents").
		WithProperties(map[string]interface{}{
			"text": text,
		}).
		WithVector(vector).
		Do(ctx)

	if err != nil {
		return fmt.Errorf("storing document: %w", err)
	}
	return nil
}

func (v *VectorDB) Search(ctx context.Context, vector []float32, limit int) ([]string, error) {
	// Take vector
	// Search Weaviate for similar vectors
	// Return similar documents
	// Called during /query/ request
	nearVector := v.weaviateClient.GraphQL().NearVectorArgBuilder().WithVector(vector)

	result, err := v.weaviateClient.GraphQL().Get().
		WithClassName("Documents").
		WithFields(graphql.Field{Name: "text"}). // return the "text" property
		WithNearVector(nearVector).
		WithLimit(limit).
		Do(ctx)

	if err != nil {
		return nil, fmt.Errorf("searching documents: %w", err)
	}

	docs := []string{}
	data := result.Data["Get"].(map[string]interface{})

	objects := data["Documents"].([]interface{})
	for _, obj := range objects {
		text := obj.(map[string]interface{})["text"].(string)
		docs = append(docs, text)
	}

	return docs, nil
}
