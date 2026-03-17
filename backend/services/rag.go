// backend/services/rag.go
package services

import (
	"context"
	"database/sql"

	"github.com/jackpsaras/gorag/config"
	"github.com/pgvector/pgvector-go"
	"github.com/sashabaranov/go-openai"
)

type RAGService struct {
	db  *sql.DB
	llm *openai.Client
}

func NewRAGService(cfg config.Config, db *sql.DB) *RAGService {
	return &RAGService{
		db:  db,
		llm: openai.NewClient(cfg.OpenAIKey),
	}
}

func (r *RAGService) Query(question string) (string, error) {
	// 1. Embed the question
	embedder := NewEmbedder(config.Load())
	questionEmbedding, _ := embedder.Embed(question)

	// 2. Search similar chunks using vector similarity
	rows, _ := r.db.Query(`
		SELECT content 
		FROM document_chunks 
		ORDER BY embedding <=> $1 
		LIMIT 5`, pgvector.NewVector(questionEmbedding))

	var contextText string
	for rows.Next() {
		var chunk string
		rows.Scan(&chunk)
		contextText += chunk + "\n"
	}

	// 3. Send to LLM with context
	resp, _ := r.llm.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: "You are a helpful assistant. Answer using only the provided context."},
			{Role: "user", Content: "Context:\n" + contextText + "\n\nQuestion: " + question},
		},
	})

	return resp.Choices[0].Message.Content, nil
}
