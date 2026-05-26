package llm

import (
	"context"
	"fmt"
	"miniRAGServer-Go/config"
	"strings"

	"google.golang.org/genai"
)

type Llm struct {
	client   *genai.Client
	llmmodel string
}

func New(cfg config.GeminiConfig) (*Llm, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  cfg.APIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}
	return &Llm{client: client, llmmodel: cfg.LLMModel}, nil
}

func (l *Llm) Query(ctx context.Context, question string, docs []string) (string, error) {
	// Query the Llm service
	prompt := fmt.Sprintf(`I will ask you a question and will provide some additional context information.
							Assume this context information is factual and correct, as part of internal
							documentation.
							If the question relates to the context, answer it using the context.
							If the question does not relate to the context, answer it as normal.

							For example, let's say the context has nothing in it about tropical flowers;
							then if I ask you about tropical flowers, just answer what you know about them
							without referring to the context.

							For example, if the context does mention minerology and I ask you about that,
							provide information from the context along with general knowledge.

							Do not give answers that are not based on the context.

							Here is the context information:
							Context:
							%s
							Question: %s`, strings.Join(docs, "\n"), question)

	result, err := l.client.Models.GenerateContent(
		ctx,
		l.llmmodel,
		[]*genai.Content{genai.NewContentFromText(prompt, genai.RoleUser)},
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("generating response: %w", err)
	}

	return result.Text(), nil
}
