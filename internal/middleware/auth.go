package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware valida a presença e validade de um Token JWT no cabeçalho "Authorization"
func AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token de autenticação não fornecido no cabeçalho."})
			return
		}

		// Valida e processa o Bearer Token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido. Use 'Bearer <token>'"})
			return
		}

		tokenString := parts[1]

		// Parse e validação fina da assinatura JWT
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			// Checar se o método de assinatura é HMAC (o que nós criamos)
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de assinatura inesperado: %v", t.Header["alg"])
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Acesso Negado: Token inválido ou expirado."})
			return
		}

		// Se tiver tudo ok, pode recuperar e gravar os claims caso queira na memória do gin context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("userID", claims["user_id"])
			c.Set("userEmail", claims["email"])
		}

		// Libera a roda para avançar pros Endpoints ou UseCases!
		c.Next()
	}
}
