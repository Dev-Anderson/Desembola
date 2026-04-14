package database

import (
	"fmt"
	"log"
	"strings"

	"desemrola/internal/config"

	"github.com/golang-migrate/migrate/v4"
	// Drivers necessários do golang-migrate
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations lê a pasta local de sql/migrations e efetua update de schemas no db.
func RunMigrations(cfg *config.AppConfig) error {
	// Padrão URL exigido pelo migrate: postgres://user:password@host:port/dbname?query
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSslMode)

	log.Println("Inciando verificação de Migrations do Banco de Dados...")
	
	// source URL aponta pro disco local onde os dados .sql estão
	m, err := migrate.New("file://sql/migrations", dbURL)
	if err != nil {
		return fmt.Errorf("erro ao configurar driver de migration: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("Tudo limpo! Nenhuma migration pendente para o PostgreSQL.")
			return nil
		}

		// Auto-Healing: Destrava o banco quando a primeira migration sujou (dirty)
		if strings.Contains(err.Error(), "Dirty database") {
			log.Println("Destravando o estado 'Dirty' do banco e forçando a versão 1!")
			
			// Força a versão 1 para arrumar o buraco (ignora o erro anterior)
			if forceErr := m.Force(1); forceErr != nil {
				return fmt.Errorf("falha ao tentar reparar o banco: %w", forceErr)
			}
			
			// Retoma o fluxo chamando o UP novamente para subir as demais versões de boas.
			if retryErr := m.Up(); retryErr != nil && retryErr != migrate.ErrNoChange {
				return fmt.Errorf("falha na tentativa final de subida de migrations: %w", retryErr)
			}
			
			log.Println("Banco destravado e Migrations avançadas com sucesso!")
			return nil
		}

		return fmt.Errorf("erro ao aplicar migrations (Up): %w", err)
	}

	log.Println("Migrations aplicadas ao banco de dados com sucesso!")
	return nil
}
