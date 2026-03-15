package services

import (
	"context"

	"github.com/jackpsaras/gorag/config"
	"github.com/sashabaranov/go-openai"
)

type Embedder struct {
	client *openai.Client
}

func NewEmbedder(cfg config.Config) *Embedder {
	return &Embedder{
		client: openai.NewClient(cfg.OpenAIKey),
	}
}

func (e *Embedder) Embed(text string) ([]float32, error) {
	resp, err := e.client.CreateEmbeddings(context.Background(), openai.EmbeddingRequest{
		Input: []string{text},
		Model: "text-embedding-3-small",
	})
	if err != nil {
		return nil, err
	}
	return resp.Data[0].Embedding, nil
}
