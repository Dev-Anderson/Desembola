package postgres

import (
	"database/sql"
	"fmt"

	"desemrola/internal/domain"
)

type ticketRepository struct {
	db *sql.DB
}

func NewTicketRepository(db *sql.DB) domain.TicketRepository {
	return &ticketRepository{db: db}
}

func (r *ticketRepository) ClearDatabase() error {
	_, err := r.db.Exec("TRUNCATE TABLE ticket, status_task CASCADE")
	return err
}

func (r *ticketRepository) GetOrCreateStatus(descricao string) (int, error) {
	var id int
	err := r.db.QueryRow("SELECT id FROM status_task WHERE descricao = $1", descricao).Scan(&id)
	if err == sql.ErrNoRows {
		err = r.db.QueryRow("INSERT INTO status_task (descricao, ativo) VALUES ($1, true) RETURNING id", descricao).Scan(&id)
		if err != nil {
			return 0, fmt.Errorf("erro ao inserir status %s: %v", descricao, err)
		}
		return id, nil
	} else if err != nil {
		return 0, fmt.Errorf("erro ao buscar status %s: %v", descricao, err)
	}
	return id, nil
}

func (r *ticketRepository) SaveTicket(t *domain.Ticket) error {
	var abertoEm sql.NullTime
	if t.AbertoEm != nil {
		abertoEm = sql.NullTime{Time: *t.AbertoEm, Valid: true}
	}

	_, err := r.db.Exec(`
		INSERT INTO ticket 
		(id, descricao, aberto_em, id_status, id_epic, departamento, categoria, responsavel)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO UPDATE SET 
			descricao = EXCLUDED.descricao,
			aberto_em = EXCLUDED.aberto_em,
			id_status = EXCLUDED.id_status,
			id_epic = EXCLUDED.id_epic,
			departamento = EXCLUDED.departamento,
			categoria = EXCLUDED.categoria,
			responsavel = EXCLUDED.responsavel
	`, t.ID, t.Descricao, abertoEm, t.StatusID, t.NumeroTask, t.Departamento, t.Categoria, t.Responsavel)

	return err
}

func (r *ticketRepository) GetAllTickets() ([]*domain.Ticket, error) {
	rows, err := r.db.Query(`
		SELECT 
			t.id, t.descricao, t.aberto_em, t.id_status, st.descricao as status_descricao,
			t.id_epic, t.departamento, t.categoria, t.responsavel
		FROM ticket t
		LEFT JOIN status_task st ON t.id_status = st.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []*domain.Ticket
	for rows.Next() {
		var t domain.Ticket
		var abertoEm sql.NullTime
		var numTask sql.NullInt64
		var statusDesc sql.NullString
		
		err := rows.Scan(
			&t.ID, &t.Descricao, &abertoEm, &t.StatusID, &statusDesc,
			&numTask, &t.Departamento, &t.Categoria, &t.Responsavel,
		)
		if err != nil {
			return nil, err
		}
		if abertoEm.Valid {
			t.AbertoEm = &abertoEm.Time
		}
		if statusDesc.Valid {
			t.Status = statusDesc.String
		}
		if numTask.Valid {
			val := int(numTask.Int64)
			t.NumeroTask = &val
		}

		tickets = append(tickets, &t)
	}

	return tickets, nil
}
