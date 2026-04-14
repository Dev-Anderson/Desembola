package usecase

import (
	"fmt"

	"desemrola/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

type CreateUserUseCase struct {
	repo domain.UserRepository
}

func NewCreateUserUseCase(repo domain.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{repo: repo}
}

func (uc *CreateUserUseCase) Execute(nome, email, senhaPura string) (*domain.User, error) {
	// Validações Simples
	if nome == "" || email == "" || senhaPura == "" {
		return nil, fmt.Errorf("campos obrigatórios não preenchidos (nome, email, senha)")
	}

	// Bcrypt transformando a string vinda em um hash incrackeavel complexo
	hash, err := bcrypt.GenerateFromPassword([]byte(senhaPura), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("erro fatal ao tentar hashear a senha: %w", err)
	}

	usr := &domain.User{
		Nome:      nome,
		Email:     email,
		SenhaHash: string(hash), // O hash é salvo
	}

	if err := uc.repo.SaveUser(usr); err != nil {
		return nil, err
	}

	// Limpar do objeto em memória a senha pura, já que agora finalizamos.
	usr.Senha = ""
	
	return usr, nil
}
