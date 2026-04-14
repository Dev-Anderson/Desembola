package usecase

import "desemrola/internal/domain"

type GetTicketsUseCase struct {
	repo domain.TicketRepository
}

func NewGetTicketsUseCase(repo domain.TicketRepository) *GetTicketsUseCase {
	return &GetTicketsUseCase{repo: repo}
}

func (uc *GetTicketsUseCase) Execute() ([]*domain.Ticket, error) {
	return uc.repo.GetAllTickets()
}
