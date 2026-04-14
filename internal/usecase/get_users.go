package usecase

import "desemrola/internal/domain"

type GetUsersUseCase struct {
	repo domain.UserRepository
}

func NewGetUsersUseCase(repo domain.UserRepository) *GetUsersUseCase {
	return &GetUsersUseCase{repo: repo}
}

func (uc *GetUsersUseCase) Execute() ([]*domain.User, error) {
	return uc.repo.GetAllUsers()
}
