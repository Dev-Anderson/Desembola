package usecase

import (
	"errors"
	"time"

	"desemrola/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginUseCase struct {
	repo      domain.UserRepository
	jwtSecret string
}

func NewLoginUseCase(repo domain.UserRepository, jwtSecret string) *LoginUseCase {
	return &LoginUseCase{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

// Execute valida as credenciais do usuário e gera um JWT novo em caso de sucesso.
func (uc *LoginUseCase) Execute(email, senha string) (string, error) {
	// 1. Busca os dados brutos no repo, incluindo a senha do banco (segura)
	user, err := uc.repo.FindByEmail(email)
	if err != nil {
		return "", errors.New("credenciais inválidas ou inexistentes")
	}

	// 2. Compara a criptografia bcrypt gravada no banco com a senha em texto limpo do Formulario
	err = bcrypt.CompareHashAndPassword([]byte(user.SenhaHash), []byte(senha))
	if err != nil {
		return "", errors.New("credenciais inválidas ou inexistentes") // Omitimos qual lado tá errado por segurança
	}

	// 3. Credenciais válidas! Vamos emitir o Passaporte Digital (JWT)
	// Payload a ser guardado publicamente (mas inalterável) dentro do Token.
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"nome":    user.Nome,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Dura 24 Horas
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(uc.jwtSecret))
	if err != nil {
		return "", errors.New("falha interna persistente ao tentar gerar token JWT de acesso")
	}

	return tokenString, nil
}
