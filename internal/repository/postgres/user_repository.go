package postgres

import (
	"database/sql"
	"fmt"

	"desemrola/internal/domain"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) SaveUser(u *domain.User) error {
	err := r.db.QueryRow(`
		INSERT INTO usuarios (nome, email, senha_hash) 
		VALUES ($1, $2, $3) RETURNING id
	`, u.Nome, u.Email, u.SenhaHash).Scan(&u.ID)
	
	if err != nil {
		return fmt.Errorf("falha ao inserir o usuário no banco: %v", err)
	}
	return nil
}

func (r *userRepository) GetAllUsers() ([]*domain.User, error) {
	rows, err := r.db.Query(`SELECT id, nome, email FROM usuarios`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var u domain.User
		// Nota: Nunca consultamos e colocamos o senha_hash de volta na memória publicamente sem haver necessidade!
		if err := rows.Scan(&u.ID, &u.Nome, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	return users, nil
}
