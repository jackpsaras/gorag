// backend/handlers/upload.go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackpsaras/gorag/services"
	"github.com/ledongthuc/pdf"
)

type UploadHandler struct {
	embedder *services.Embedder
}

func NewUploadHandler(embedder *services.Embedder) *UploadHandler {
	return &UploadHandler{embedder: embedder}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("document")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Only accept PDFs
	if file.Header.Get("Content-Type") != "application/pdf" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only PDF files are supported"})
		return
	}

	// Open the uploaded PDF
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer f.Close()

	// Extract text from PDF using ledongthuc/pdf
	reader, err := pdf.NewReader(f, file.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read PDF"})
		return
	}

	var fullText string
	for i := 1; i <= reader.GetNumPages(); i++ {
		page, _ := reader.GetPage(i)
		text, _ := page.GetPlainText(nil)
		fullText += text + "\n"
	}

	// Simple chunking (you can improve this later)
	chunk := fullText[:min(1000, len(fullText))] // First 1000 chars for demo

	// Embed the chunk
	embedding, err := h.embedder.Embed(chunk)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to embed text"})
		return
	}

	// TODO: Save chunk + embedding to database (we'll do this properly in Day 5)

	c.JSON(http.StatusOK, gin.H{
		"message":  "PDF uploaded and processed successfully",
		"filename": file.Filename,
		"chunk":    chunk[:100] + "...", // Show preview
	})
}

// Helper function for min (Go doesn't have built-in min for int until 1.21+)
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
