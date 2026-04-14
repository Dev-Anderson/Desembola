package usecase

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"desemrola/internal/domain"

	"golang.org/x/text/encoding/charmap"
)

type ImportTicketsUseCase struct {
	repo domain.TicketRepository
}

func NewImportTicketsUseCase(repo domain.TicketRepository) *ImportTicketsUseCase {
	return &ImportTicketsUseCase{repo: repo}
}

func (uc *ImportTicketsUseCase) Execute(file io.Reader) error {
	decoder := charmap.Windows1252.NewDecoder().Reader(file)
	reader := csv.NewReader(decoder)
	reader.Comma = ';'
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("erro ao ler o CSV: %w", err)
	}

	layoutDate := "02/01/2006 15:04"
	registrosInseridos := 0

	for i, row := range records {
		if i == 0 || len(row) < 8 {
			continue
		}

		numeroStr := strings.TrimSpace(row[0])
		assunto := strings.TrimSpace(row[1])
		abertoEmStr := strings.TrimSpace(row[2])
		statusStr := strings.TrimSpace(row[3])
		numTaskStr := strings.TrimSpace(row[4])
		departamento := strings.TrimSpace(row[5])
		categoria := strings.TrimSpace(row[6])
		responsavel := strings.TrimSpace(row[7])

		numero, err := strconv.Atoi(numeroStr)
		if err != nil {
			log.Printf("Linha %d: Erro ao converter Numero '%s': %v", i+1, numeroStr, err)
			continue
		}

		var abertoEm *time.Time
		if abertoEmStr != "" {
			pt, err := time.Parse(layoutDate, abertoEmStr)
			if err != nil {
				log.Printf("Linha %d: Erro ao dar parse na data '%s': %v", i+1, abertoEmStr, err)
				continue
			}
			abertoEm = &pt
		}

		idStatus, err := uc.repo.GetOrCreateStatus(statusStr)
		if err != nil {
			log.Printf("Linha %d: %v", i+1, err)
			continue
		}

		numeroTask := parseNumeroTask(numTaskStr)

		ticket := &domain.Ticket{
			ID:           numero,
			Descricao:    assunto,
			AbertoEm:     abertoEm,
			Status:       statusStr,
			StatusID:     idStatus,
			NumeroTask:   numeroTask,
			Departamento: departamento,
			Categoria:    categoria,
			Responsavel:  responsavel,
		}

		if err := uc.repo.SaveTicket(ticket); err != nil {
			log.Printf("Linha %d: Erro ao salvar ticket %d: %v", i+1, numero, err)
			continue
		}

		registrosInseridos++
	}

	fmt.Printf("Migração finalizada com sucesso. %d registros lidos e processados para a tabela ticket.\n", registrosInseridos)
	return nil
}

func parseNumeroTask(s string) *int {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	s = strings.ReplaceAll(s, ",00", "")
	val, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	return &val
}
