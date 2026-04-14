package handler

import "github.com/gin-gonic/gin"

// SetupRoutes centraliza e agrupa todos os mapeamentos de endpoints 
// deixando o arquivo main.go blindado e os roteamentos fáceis de visualizar.
func SetupRoutes(r *gin.Engine, ticketHandler *TicketHandler, userHandler *UserHandler) {
	// Grupo público / usuários normais
	publicApi := r.Group("/api/tickets")
	{
		publicApi.POST("/upload", ticketHandler.UploadCSV)
		publicApi.GET("/all", ticketHandler.GetTickets)
	}

	// Módulo Isolado de Autenticação / Base de Usuarios Livre
	userApi := r.Group("/api/users")
	{
		userApi.POST("/register", userHandler.Register)
	}

	// Grupo exclusivo para administradores
	adminApi := r.Group("/api/admin/tickets")
	// Pode-se facilmente usar adminApi.Use(AuthMiddleware()) para blindar tudo aqui dentro!
	{
		adminApi.POST("/clean", ticketHandler.CleanDB)
	}

	adminUsersApi := r.Group("/api/admin/users")
	{
		adminUsersApi.GET("/all", userHandler.GetAll)
	}
}
