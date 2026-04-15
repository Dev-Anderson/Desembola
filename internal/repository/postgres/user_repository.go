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

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var u domain.User
	// Precisamos do senha_hash apenas no momento do Login pelo banco,
	// porém o campo json:"-" do domain.User garante que não vaze para APIs públicas caso o Dev falhe.
	err := r.db.QueryRow(`SELECT id, nome, email, senha_hash FROM usuarios WHERE email = $1`, email).
		Scan(&u.ID, &u.Nome, &u.Email, &u.SenhaHash)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("usuário não encontrado")
		}
		return nil, err
	}

	return &u, nil
}
