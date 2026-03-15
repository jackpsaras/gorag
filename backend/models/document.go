package models

import (
	"github.com/pgvector/pgvector-go"
)

type DocumentChunk struct {
	ID        int
	Document  string
	Content   string
	Embedding pgvector.Vector `db:"embedding"`
}
