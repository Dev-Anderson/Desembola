package domain

import "time"

type Ticket struct {
	ID           int        `json:"id"`
	Descricao    string     `json:"descricao"`
	AbertoEm     *time.Time `json:"aberto_em,omitempty"`
	Status       string     `json:"status"`
	StatusID     int        `json:"id_status"`
	NumeroTask   *int       `json:"numero_task,omitempty"`
	Departamento string     `json:"departamento"`
	Categoria    string     `json:"categoria"`
	Responsavel  string     `json:"responsavel"`
}

type TicketRepository interface {
	ClearDatabase() error
	GetOrCreateStatus(descricao string) (int, error)
	SaveTicket(t *Ticket) error
	GetAllTickets() ([]*Ticket, error)
}
