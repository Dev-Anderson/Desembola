package domain

type User struct {
	ID        int    `json:"id"`
	Nome      string `json:"nome"`
	Email     string `json:"email"`
	Senha     string `json:"senha,omitempty"` // Campo que a web/usario nos passa cru
	SenhaHash string `json:"-"`               // Invisivel na serialização JSON de saída
}

type UserRepository interface {
	SaveUser(u *User) error
	GetAllUsers() ([]*User, error)
	FindByEmail(email string) (*User, error)
}
