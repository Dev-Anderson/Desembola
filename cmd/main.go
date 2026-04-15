package main

import (
	"log"

	"desemrola/internal/config"
	"desemrola/internal/database"
	"desemrola/internal/handler"
	"desemrola/internal/repository/postgres"
	"desemrola/internal/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Carrega configuração do env
	cfg := config.LoadConfig(".env")

	// 2. Conecta ao banco de dados isolado na pasta de inicialização de infraestrutura
	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("Falha crítica ao subir aplicação sem banco de dados: %v", err)
	}
	defer db.Close()

	// [NOVO] - Engrenagem da Migration entra em campo aqui:
	if err := database.RunMigrations(cfg); err != nil {
		log.Fatalf("Ocorreu um erro letal ao checar versões do Banco de Dados: %v", err)
	}

	// 3. Inicializa as camadas do sistema (Injeção de dependência)
	ticketRepo := postgres.NewTicketRepository(db)
	userRepo := postgres.NewUserRepository(db)
	
	// Construção das Classes do Sistema de Tickets
	importUC := usecase.NewImportTicketsUseCase(ticketRepo)
	cleanUC := usecase.NewCleanDatabaseUseCase(ticketRepo)
	getUC := usecase.NewGetTicketsUseCase(ticketRepo)

	// Construção das Classes do Sistema de Login/Acesso
	createUserUC := usecase.NewCreateUserUseCase(userRepo)
	getUsersUC := usecase.NewGetUsersUseCase(userRepo)
	loginUC := usecase.NewLoginUseCase(userRepo, cfg.JWTSecret)

	// Instanciando Controllers
	ticketHandler := handler.NewTicketHandler(importUC, cleanUC, getUC)
	userHandler := handler.NewUserHandler(createUserUC, getUsersUC, loginUC)

	// 4. Configurando Rotas Web HTTP / GIN e as delegando!
	r := gin.Default()
	handler.SetupRoutes(r, ticketHandler, userHandler, cfg.JWTSecret)

	// 5. Iniciar Escuta na Porta Definida no .env
	port := cfg.APIPort
	if port == "" {
		port = "8080"
	}

	log.Printf("Sua API está subindo na porta %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Erro crítico ao tentar inicializar o servidor Web: %v", err)
	}
}
