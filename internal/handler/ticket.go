package handler

import (
	"log"
	"net/http"

	"desemrola/internal/usecase"

	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	importUC *usecase.ImportTicketsUseCase
	cleanUC  *usecase.CleanDatabaseUseCase
	getUC    *usecase.GetTicketsUseCase
}

func NewTicketHandler(importUC *usecase.ImportTicketsUseCase, cleanUC *usecase.CleanDatabaseUseCase, getUC *usecase.GetTicketsUseCase) *TicketHandler {
	return &TicketHandler{
		importUC: importUC,
		cleanUC:  cleanUC,
		getUC:    getUC,
	}
}

// UploadCSV recebe um arquivo via multipart/form-data e processa diretamente para o banco
func (h *TicketHandler) UploadCSV(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Arquivo CSV não enviado ou campo 'file' ausente no form-data."})
		return
	}

	f, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao tentar ler o arquivo enviado."})
		return
	}
	defer f.Close()

	if err := h.importUC.Execute(f); err != nil {
		log.Printf("Erro na importação: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "CSV importado e processado no banco de dados com sucesso!"})
}

type CleanRequest struct {
	Clean bool `json:"limpar"`
}

// CleanDB faz parse no body para {"limpar": true/false}
func (h *TicketHandler) CleanDB(c *gin.Context) {
	var req CleanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Corpo da requisição inválido. Envie um JSON como {\"limpar\": true}."})
		return
	}

	if err := h.cleanUC.Execute(req.Clean); err != nil {
		log.Printf("Erro ao limpar banco: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao limpar os dados no banco."})
		return
	}

	if req.Clean {
		c.JSON(http.StatusOK, gin.H{"message": "O banco de dados de tickets foi deletado com sucesso!"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Nenhuma ação realizada pois o valor de 'limpar' era false."})
	}
}

// GetTickets retorna todos os tickets num json de array limpo
func (h *TicketHandler) GetTickets(c *gin.Context) {
	tickets, err := h.getUC.Execute()
	if err != nil {
		log.Printf("Erro na busca de tickets: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao consultar sistema de tickets."})
		return
	}
	
	// Retornando a slice preenchida (se vazia, ele retorna nulo/vazio JSON corretamente)
	c.JSON(http.StatusOK, tickets)
}
