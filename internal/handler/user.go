package handler

import (
	"log"
	"net/http"

	"desemrola/internal/usecase"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	createUC *usecase.CreateUserUseCase
	getUC    *usecase.GetUsersUseCase
	loginUC  *usecase.LoginUseCase
}

func NewUserHandler(createUC *usecase.CreateUserUseCase, getUC *usecase.GetUsersUseCase, loginUC *usecase.LoginUseCase) *UserHandler {
	return &UserHandler{
		createUC: createUC,
		getUC:    getUC,
		loginUC:  loginUC,
	}
}

type RegisterRequest struct {
	Nome  string `json:"nome"`
	Email string `json:"email"`
	Senha string `json:"senha"`
}

// Register trata do endpoint de cadastro que pode ser invocado pelo frontend
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato inválido. Envie os campos 'nome', 'email' e 'senha'."})
		return
	}

	user, err := h.createUC.Execute(req.Nome, req.Email, req.Senha)
	if err != nil {
		log.Printf("Erro de cadastro: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível registrar o usuário. (Talvez o e-mail já exista?)"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Bem vindo! Usuário cadastrado com sucesso.",
		"usuario": user,
	})
}

// GetAll é o ponto de consulta administrativo
func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.getUC.Execute()
	if err != nil {
		log.Printf("Erro de visualização: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao tentar buscar usuários do sistema."})
		return
	}

	c.JSON(http.StatusOK, users)
}

type LoginRequest struct {
	Email string `json:"email"`
	Senha string `json:"senha"`
}

// Login valida o e-mail/senha e devolve um JWT
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verifique os dados enviados (email, senha)."})
		return
	}

	token, err := h.loginUC.Execute(req.Email, req.Senha)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
