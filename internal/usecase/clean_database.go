package usecase

import (
	"fmt"
	"log"

	"desemrola/internal/domain"
)

type CleanDatabaseUseCase struct {
	repo domain.TicketRepository
}

func NewCleanDatabaseUseCase(repo domain.TicketRepository) *CleanDatabaseUseCase {
	return &CleanDatabaseUseCase{repo: repo}
}

func (uc *CleanDatabaseUseCase) Execute(clean bool) error {
	if !clean {
		return nil
	}

	log.Println("Executando operação: Limpando todos os tickets e status antigos...")
	if err := uc.repo.ClearDatabase(); err != nil {
		return fmt.Errorf("erro ao limpar banco de dados: %w", err)
	}
	return nil
}
