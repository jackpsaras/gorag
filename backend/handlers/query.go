// backend/handlers/query.go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackpsaras/gorag/services"
)

type QueryHandler struct {
	rag *services.RAGService
}

func NewQueryHandler(rag *services.RAGService) *QueryHandler {
	return &QueryHandler{rag: rag}
}

func (h *QueryHandler) Query(c *gin.Context) {
	var req struct {
		Question string `json:"question"`
	}
	c.BindJSON(&req)

	answer, _ := h.rag.Query(req.Question)

	c.JSON(http.StatusOK, gin.H{"answer": answer})
}
