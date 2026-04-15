package handler

import (
	"desemrola/internal/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRoutes centraliza e agrupa todos os mapeamentos de endpoints 
// deixando o arquivo main.go blindado e os roteamentos fáceis de visualizar.
func SetupRoutes(r *gin.Engine, ticketHandler *TicketHandler, userHandler *UserHandler, jwtSecret string) {
	// Grupo público / usuários normais - Agora protegido por Autenticação!
	publicApi := r.Group("/api/tickets").Use(middleware.AuthMiddleware(jwtSecret))
	{
		publicApi.POST("/upload", ticketHandler.UploadCSV)
		publicApi.GET("/all", ticketHandler.GetTickets)
	}

	// Módulo Isolado de Autenticação / Base de Usuarios Livre
	userApi := r.Group("/api/users")
	{
		userApi.POST("/register", userHandler.Register)
		userApi.POST("/login", userHandler.Login)
	}

	// Grupo exclusivo para administradores
	adminApi := r.Group("/api/admin/tickets").Use(middleware.AuthMiddleware(jwtSecret))
	{
		adminApi.POST("/clean", ticketHandler.CleanDB)
	}

	adminUsersApi := r.Group("/api/admin/users").Use(middleware.AuthMiddleware(jwtSecret))
	{
		adminUsersApi.GET("/all", userHandler.GetAll)
	}
}
