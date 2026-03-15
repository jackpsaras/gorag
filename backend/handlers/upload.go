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

	// Open the uploaded file
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
	numPages := reader.NumPage()

	for i := 1; i <= numPages; i++ {
		page := reader.Page(i)
		text, _ := page.GetPlainText(nil)
		fullText += text + "\n"
	}

	// Simple chunking
	chunk := fullText
	if len(chunk) > 1000 {
		chunk = chunk[:1000]
	}

	// Create embedding (this fixes "declared and not used")
	embedding, err := h.embedder.Embed(chunk)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create embedding"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "PDF uploaded and processed successfully",
		"filename":  file.Filename,
		"pages":     numPages,
		"chunk":     chunk + "...",
		"embedding": embedding,
	})
}
