package database

import (
	"database/sql"
	"fmt"

	"desemrola/internal/config"

	_ "github.com/lib/pq"
)

// NewPostgresConnection formata a string de conexão, conecta e envia um ping de teste
// para garantir que o banco de dados oficial está operando perfeitamente.
func NewPostgresConnection(cfg *config.AppConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSslMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("falha ao instanciar o pacote postgres: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("não foi possível pingar a porta do servidor de banco: %w", err)
	}

	return db, nil
}
